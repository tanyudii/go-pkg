package go_grpc

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
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
	RpcDurationsHistogram   = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "grpc_rpc_durations_histogram",
		Help:    "GRPC RPC latency distributions.",
		Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000},
	}, []string{"httpCode", "grpcCode", "grpcMethod", "statusCode"})
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
	ListenAndServePrometheus(ctx context.Context) error
	RegisterUnaryServerInterceptor(i ...grpc.UnaryServerInterceptor)
	RegisterRESTHandler(handlers ...RESTHandler)
}

type service struct {
	cfg                  *Config
	server               *grpc.Server
	interceptors         Interceptors
	restHandlers         []RESTHandler
	prometheusCollectors []prometheus.Collector
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
	s.initDefaultPrometheusCollectors()
	s.registerHealthServer()
}

func (s *service) Shutdown(ctx context.Context) error {
	<-ctx.Done()
	return nil
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
	gologger.Infof("go grpc is shutting down: for %ds %v", t, time.Now())
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t)*time.Second)
	defer cancel()
	cancelMainCtx()
	if err := s.Shutdown(ctx); err != nil {
		gologger.Fatalf("go grpc shutdown err: %v", err)
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
		gologger.Infof("go grpc initializing gRPC connection in port %s", s.cfg.gRPCPort)
		exitFunc(s.ListenAndServeGRPC(ctx))
	})

	go wg.Wrap(func() {
		gologger.Infof("go grpc initializing HTTP connection in port %s", s.cfg.restPort)
		exitFunc(s.ListenAndServeREST(ctx))
	})

	go wg.Wrap(func() {
		gologger.Infof("go grpc initializing Prometheus connection in port %s", s.cfg.prometheusPort)
		exitFunc(s.ListenAndServePrometheus(ctx))
	})

	return exitCh
}

func (s *service) ListenAndServeGRPC(_ context.Context) error {
	if s.server == nil {
		return ErrServerNotInitialized
	}
	gologger.Infof("go grpc listen and serve grpc: %v", s.cfg.gRPCPort)

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

	srv := &http.Server{
		Addr:    ":" + s.cfg.restPort,
		Handler: MuxCORS(handler),
	}

	go func() {
		<-ctx.Done()
		if err = srv.Shutdown(context.Background()); err != nil {
			gologger.Errorf("go grpc listen and serve rest: failed to shutdown %v", err)
		}
	}()

	gologger.Infof("go grpc listen and serve rest: %v", s.cfg.restPort)
	if err = srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		gologger.Errorf("go grpc listen and serve rest: failed to listen and serve %v", err)
		return err
	}

	return nil
}

func (s *service) ListenAndServePrometheus(ctx context.Context) (err error) {
	for _, c := range s.prometheusCollectors {
		if err = prometheus.Register(c); err != nil {
			gologger.Errorf("go grpc listen and serve prometheus: failed to register %v", err)
			return err
		}
	}

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.cfg.prometheusPort),
		Handler: MuxCORS(mux),
	}

	mux.Handle("/metrics", promhttp.Handler())

	go func() {
		<-ctx.Done()
		if err = srv.Shutdown(context.Background()); err != nil {
			gologger.Errorf("go grpc listen and serve prometheus: failed to shutdown %v", err)
		}
	}()

	gologger.Infof("go grpc listen and serve prometheus: %v", s.cfg.prometheusPort)
	if err = srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		gologger.Errorf("go grpc listen and serve prometheus: failed to listen and serve %v", err)
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
		runtime.WithIncomingHeaderMatcher(MuxIncomingHeaderMatcher),
		runtime.WithForwardResponseOption(MuxHandleRoutingRedirect),
		runtime.WithHealthEndpointAt(grpc_health_v1.NewHealthClient(ClientConn(":"+s.cfg.gRPCPort)), "/_health"),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: s.cfg.discardUnknown,
			},
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:     true,
				EmitUnpopulated:   true,
				EmitDefaultValues: true,
			},
		}),
	)
}

func (s *service) initGRPCServer() {
	s.server = grpc.NewServer(grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(s.interceptors.serverUnary...)))
}

func (s *service) initReflection() {
	reflection.Register(s.GetServer())
}

func (s *service) initDefaultPrometheusCollectors() {
	s.prometheusCollectors = append(s.prometheusCollectors, RpcDurationsHistogram)
}

func (s *service) initRESTHandler(ctx context.Context) (http.Handler, error) {
	mux := runtime.NewServeMux(s.cfg.restServeMuxOpts...)

	creds := insecure.NewCredentials()
	if s.cfg.tls {
		creds = credentials.NewTLS(&tls.Config{})
	}

	endpoint := ":" + s.cfg.gRPCPort
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	for i := range s.restHandlers {
		h := s.restHandlers[i]
		if err := h(ctx, mux, endpoint, opts); err != nil {
			return nil, err
		}
	}

	return mux, nil
}

func (s *service) registerHealthServer() {
	grpc_health_v1.RegisterHealthServer(s.GetServer(), newHealthServer())
}
