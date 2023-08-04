package go_grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	ErrServerNotInitialized = errors.New("server not initialized")
)

type RESTHandler func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

type Service interface {
	Init()
	Shutdown(ctx context.Context) error
	GetServer() *grpc.Server
	RunGracefully(t int)
	RunServers(ctx context.Context) <-chan error
	ListenAndServeGRPC(ctx context.Context) error
	ListenAndServeREST(ctx context.Context) error
	RegisterUnaryServerInterceptor(i ...grpc.UnaryServerInterceptor)
	RegisterRESTHandler(handlers ...RESTHandler)
}

type service struct {
	cfg          *Config
	server       *grpc.Server
	interceptors Interceptors
	restHandlers []RESTHandler
}

type Interceptors struct {
	serverUnary []grpc.UnaryServerInterceptor
}

func NewService(args ...ConfigFunc) Service {
	return &service{
		cfg: generate(args...),
	}
}

func (s *service) Init() {
	s.initInterceptors()
	s.initConfigRestServeMuxOpts()
	s.initGRPCServer()
	s.initReflection()
}

func (s *service) Shutdown(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (s *service) RunGracefully(t int) {
	mainCtx, cancelMainCtx := context.WithCancel(context.Background())
	go func() {
		if err := <-s.RunServers(mainCtx); err != nil {
			fmt.Printf("gogo run servers err: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Printf("gogo is shutting down: for %ds %v\n", t, time.Now())
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t)*time.Second)
	defer cancel()
	cancelMainCtx()
	if err := s.Shutdown(ctx); err != nil {
		fmt.Printf("gogo shutdown err: %v\n", err)
	}
	fmt.Printf("gogo shutdown gracefully: %v\n", time.Now())
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
		fmt.Printf("gogo Initializing gRPC connection in port %s\n", s.cfg.gRPCPort)
		exitFunc(s.ListenAndServeGRPC(ctx))
	})

	go wg.Wrap(func() {
		fmt.Printf("gogo Initializing HTTP connection in port %s", s.cfg.restPort)
		exitFunc(s.ListenAndServeREST(ctx))
	})

	return exitCh
}

func (s *service) ListenAndServeGRPC(_ context.Context) error {
	if s.server == nil {
		return ErrServerNotInitialized
	}
	fmt.Printf("gogo listen and serve grpc: %v\n", s.cfg.gRPCPort)

	defer s.server.GracefulStop()
	lis, err := net.Listen("tcp", ":"+s.cfg.gRPCPort)
	if err != nil {
		return err
	}

	return s.server.Serve(lis)
}

func (s *service) ListenAndServeREST(ctx context.Context) error {
	handler, err := s.initRESTHandler(ctx)
	if err != nil {
		return err
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	srv := &http.Server{
		Addr:    ":" + s.cfg.restPort,
		Handler: r,
	}

	if s.cfg.enableCORS {
		r.Use(GinCORS())
	}
	if s.cfg.onlyJSON {
		r.Use(GinJSON())
	}

	r.Group("*{any}").Any("", gin.WrapH(handler))

	go func() {
		<-ctx.Done()
		if err = srv.Shutdown(context.Background()); err != nil {
			fmt.Printf("gogo listen and serve rest: failed to shutdown %v\n", err)
		}
	}()

	fmt.Printf("gogo listen and serve rest: %v\n", s.cfg.restPort)
	if err = srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("gogo listen and serve rest: failed to listen and serve %v", err)
		return err
	}

	return nil
}

func (s *service) GetServer() *grpc.Server {
	return s.server
}

func (s *service) RegisterUnaryServerInterceptor(interceptors ...grpc.UnaryServerInterceptor) {
	s.interceptors.serverUnary = append(s.interceptors.serverUnary, interceptors...)
}

func (s *service) RegisterRESTHandler(handlers ...RESTHandler) {
	s.restHandlers = append(s.restHandlers, handlers...)
}

func (s *service) initInterceptors() {
	s.RegisterUnaryServerInterceptor(
		RequestIDUnaryServerInterceptor(),
		RecoveryUnaryServerInterceptor(),
		AcceptLangUnaryServerInterceptor(),
	)
}

func (s *service) initConfigRestServeMuxOpts() {
	s.cfg.restServeMuxOpts = append(
		s.cfg.restServeMuxOpts,
		runtime.WithRoutingErrorHandler(MuxHandleRoutingError),
		runtime.WithErrorHandler(MuxErrorHandler),
	)
}

func (s *service) initGRPCServer() {
	s.server = grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(s.interceptors.serverUnary...)))
}

func (s *service) initReflection() {
	reflection.Register(s.GetServer())
}

func (s *service) initRESTHandler(ctx context.Context) (http.Handler, error) {
	mux := runtime.NewServeMux(s.cfg.restServeMuxOpts...)

	conn, err := s.dialSelf()
	if err != nil {
		return nil, err
	}

	if err = s.initHealthCheck(mux, conn); err != nil {
		return nil, err
	}

	endpoint := ":" + s.cfg.gRPCPort
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	for i := range s.restHandlers {
		h := s.restHandlers[i]
		if err := h(ctx, mux, endpoint, opts); err != nil {
			return nil, err
		}
	}

	return mux, nil
}

func (s *service) initHealthCheck(mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return mux.HandlePath(http.MethodGet, "/_health", func(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
		w.Header().Set("Content-Type", "text/plain")
		if state := conn.GetState(); state != connectivity.Ready {
			http.Error(w, fmt.Sprintf("gRPC server is %s", state), http.StatusBadGateway)
			return
		}
	})
}

func (s *service) dialSelf() (*grpc.ClientConn, error) {
	return dial("tcp", fmt.Sprintf("127.0.0.1:%s", s.cfg.gRPCPort))
}
