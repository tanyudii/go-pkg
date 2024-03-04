package go_grpc

import "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

const (
	DefaultGRPCPort           = "5758"
	DefaultRESTPort           = "8080"
	DefaultEnableCORS         = true
	DefaultOnlyJSON           = true
	DefaultRegisterReflection = true
)

type Config struct {
	gRPCPort           string
	restPort           string
	enableCORS         bool
	onlyJSON           bool
	registerReflection bool
	tls                bool
	restServeMuxOpts   []runtime.ServeMuxOption
}

type ConfigFunc func(c *Config)

func GRPCPort(p string) ConfigFunc {
	if p == "" {
		p = DefaultGRPCPort
	}
	return func(c *Config) {
		c.gRPCPort = p
	}
}

func RESTPort(p string) ConfigFunc {
	if p == "" {
		p = DefaultRESTPort
	}
	return func(c *Config) {
		c.restPort = p
	}
}

func EnableCORS(cors bool) ConfigFunc {
	return func(c *Config) {
		c.enableCORS = cors
	}
}

func OnlyJSON(j bool) ConfigFunc {
	return func(c *Config) {
		c.onlyJSON = j
	}
}

func RegisterReflection(r bool) ConfigFunc {
	return func(c *Config) {
		c.registerReflection = r
	}
}

func WithTLS(tls bool) ConfigFunc {
	return func(c *Config) {
		c.tls = tls
	}
}

func AddRestServerMuxOpt(opt ...runtime.ServeMuxOption) ConfigFunc {
	return func(c *Config) {
		c.restServeMuxOpts = append(c.restServeMuxOpts, opt...)
	}
}

func generate(args ...ConfigFunc) *Config {
	c := &Config{
		gRPCPort:           DefaultGRPCPort,
		restPort:           DefaultRESTPort,
		enableCORS:         DefaultEnableCORS,
		onlyJSON:           DefaultOnlyJSON,
		registerReflection: DefaultRegisterReflection,
	}
	for i := range args {
		args[i](c)
	}
	return c
}
