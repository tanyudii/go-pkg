package gin_auth

type Config struct {
	graphqlMode bool
}

type ConfigFunc func(c *Config)

func GraphQLMode(m bool) ConfigFunc {
	return func(c *Config) {
		c.graphqlMode = m
	}
}

func generate(args ...ConfigFunc) *Config {
	c := &Config{}
	for i := range args {
		args[i](c)
	}
	return c
}
