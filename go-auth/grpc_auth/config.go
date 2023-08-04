package grpc_auth

type Config struct {
	InternalCallPassword string
}

type ConfigFunc func(c *Config)

func InternalCallPassword(pwd string) ConfigFunc {
	return func(c *Config) {
		c.InternalCallPassword = pwd
	}
}

func generate(args ...ConfigFunc) *Config {
	c := &Config{}
	for i := range args {
		args[i](c)
	}
	return c
}
