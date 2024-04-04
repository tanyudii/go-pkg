package go_graphql

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	gologger "pkg.tanyudii.me/go-pkg/go-logger"
	"sync"
	"syscall"
	"time"
)

var (
	ErrServerNotInitialized = errors.New("[ERROR]: Server not initialized")
)

type Service interface {
	Shutdown(ctx context.Context) error
	GetServer() *handler.Server
	RunGracefully(t int)
	RunServers(ctx context.Context) <-chan error
	RegisterExecutableSchema(schema graphql.ExecutableSchema)
	RegisterMiddleware(m gin.HandlerFunc)
}

type service struct {
	cfg        *Config
	server     *handler.Server
	schema     graphql.ExecutableSchema
	middleware []gin.HandlerFunc
}

func NewService(args ...ConfigFunc) Service {
	return &service{
		cfg: generate(args...),
	}
}

func (s *service) Init() {
	s.initHandlerServer()
}

func (s *service) Shutdown(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (s *service) GetServer() *handler.Server {
	return s.server
}

func (s *service) RunGracefully(t int) {
	mainCtx, cancelMainCtx := context.WithCancel(context.Background())
	go func() {
		if err := <-s.RunServers(mainCtx); err != nil {
			gologger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	gologger.Infof("go graphql is shutting down: for %ds %v", t, time.Now())
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t)*time.Second)
	defer cancel()
	cancelMainCtx()
	if err := s.Shutdown(ctx); err != nil {
		gologger.Fatalf("go graphql shutdown err: %v", err)
	}
	gologger.Infof("go grpc shutdown gracefully: %v", time.Now())
}

func (s *service) RunServers(ctx context.Context) <-chan error {
	var once sync.Once
	exitCh := make(chan error)
	wg := &Wg{}
	exitFunc := func(err error) {
		once.Do(func() {
			exitCh <- err
		})
	}

	go wg.Wrap(func() {
		gologger.Infof("go graphql initializing graphQL connection in port %s", s.cfg.graphQLPort)
		exitFunc(s.ListenAndServeGraphQL(ctx))
	})

	return exitCh
}

func (s *service) ListenAndServeGraphQL(ctx context.Context) (err error) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	srv := &http.Server{
		Addr:    ":" + s.cfg.graphQLPort,
		Handler: r,
	}

	r.Use(GinRequestID())
	r.Use(GinAcceptLanguage())
	if s.cfg.enableCORS {
		r.Use(GinCORS())
	}

	for _, m := range s.middleware {
		r.Use(m)
	}

	s.initHealthCheck(r)

	graphqlHandler, err := s.graphQLHandler()
	if err != nil {
		return err
	}

	r.POST(s.cfg.graphQLPath, graphqlHandler)

	if s.cfg.enablePlayground && s.cfg.playgroundPath != "" {
		r.GET(s.cfg.playgroundPath, s.playgroundHandler())
	}

	go func() {
		<-ctx.Done()
		if err = srv.Shutdown(context.Background()); err != nil {
			gologger.Errorf("go graphql listen and serve graphQL: failed to shutdown %v", err)
		}
	}()

	gologger.Infof("go graphql listen and serve graphQL: %v", s.cfg.graphQLPort)
	if err = srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		gologger.Errorf("go graphql listen and serve graphQL: failed to listen and serve %v", err)
		return err
	}

	return nil
}

func (s *service) RegisterExecutableSchema(schema graphql.ExecutableSchema) {
	s.schema = schema
}

func (s *service) RegisterMiddleware(m gin.HandlerFunc) {
	s.middleware = append(s.middleware, m)
}

func (s *service) initHandlerServer() {
	s.server = handler.NewDefaultServer(s.schema)
	s.server.SetRecoverFunc(Recover)
}

func (s *service) graphQLHandler() (gin.HandlerFunc, error) {
	if s.server == nil {
		return nil, ErrServerNotInitialized
	}
	s.server.SetRecoverFunc(Recover)
	s.server.SetErrorPresenter(s.cfg.errorPresenterFunc)
	return func(c *gin.Context) {
		s.server.ServeHTTP(c.Writer, c.Request)
	}, nil
}

func (s *service) playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", s.cfg.graphQLPath)
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *service) initHealthCheck(r *gin.Engine) {
	r.GET("/_health", func(c *gin.Context) {
		c.Header("Content-Type", "text/plain")
	})
}
