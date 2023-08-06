package go_graphql

import "github.com/99designs/gqlgen/graphql"

const (
	DefaultGraphQLPort      = "4000"
	DefaultGraphQLPath      = "/graphql"
	DefaultPlaygroundPath   = "/playground"
	DefaultEnablePlayground = true
	DefaultEnableCORS       = true
)

type Config struct {
	graphQLPort        string
	graphQLPath        string
	playgroundPath     string
	enableCORS         bool
	enablePlayground   bool
	recoverFunc        graphql.RecoverFunc
	errorPresenterFunc graphql.ErrorPresenterFunc
}

type ConfigFunc func(c *Config)

func GraphQLPort(p string) ConfigFunc {
	if p == "" {
		p = DefaultGraphQLPort
	}
	return func(c *Config) {
		c.graphQLPort = p
	}
}

func GraphQLPath(p string) ConfigFunc {
	if p == "" {
		p = DefaultGraphQLPath
	}
	return func(c *Config) {
		c.graphQLPath = p
	}
}

func PlaygroundPath(p string) ConfigFunc {
	if p == "" {
		p = DefaultPlaygroundPath
	}
	return func(c *Config) {
		c.playgroundPath = p
	}
}

func EnablePlayground(e bool) ConfigFunc {
	return func(c *Config) {
		c.enablePlayground = e
	}
}

func EnableCORS(e bool) ConfigFunc {
	return func(c *Config) {
		c.enableCORS = e
	}
}

func RecoverFunc(fn graphql.RecoverFunc) ConfigFunc {
	return func(c *Config) {
		c.recoverFunc = fn
	}
}

func ErrorPresenterFunc(fn graphql.ErrorPresenterFunc) ConfigFunc {
	return func(c *Config) {
		c.errorPresenterFunc = fn
	}
}

func generate(args ...ConfigFunc) *Config {
	c := &Config{
		graphQLPort:        DefaultGraphQLPort,
		graphQLPath:        DefaultGraphQLPath,
		playgroundPath:     DefaultPlaygroundPath,
		enablePlayground:   DefaultEnablePlayground,
		enableCORS:         DefaultEnableCORS,
		recoverFunc:        Recover,
		errorPresenterFunc: ErrorPresenter,
	}
	for i := range args {
		args[i](c)
	}
	return c
}
