package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/kelseyhightower/envconfig"
	gologger "pkg.tanyudii.me/go-pkg/go-logger"
	"strings"
)

type config struct {
	Username string `envconfig:"REDIS_USERNAME"`
	Password string `envconfig:"REDIS_PASSWORD"`
	Host     string `envconfig:"REDIS_HOST" required:"127.0.0.1"`
	Port     string `envconfig:"REDIS_PORT" required:"6370"`
	Network  string `envconfig:"REDIS_NETWORK" default:"tcp"`

	client *redis.Client
}

func (c *config) connect() *redis.Client {
	if c.client != nil {
		return c.client
	}
	c.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.Host, c.Port),
		Username: c.Username,
		Password: c.Password,
	})
	return c.client
}

var mapCfg = make(map[string]*config)

func Connect(name ...string) *redis.Client {
	var prefix string
	if len(name) > 0 {
		prefix = strings.ToUpper(name[0])
	}
	if cfg, ok := mapCfg[prefix]; ok {
		return cfg.client
	}
	cfg := &config{}
	envconfig.MustProcess(prefix, cfg)
	cfg.connect()
	mapCfg[prefix] = cfg
	return cfg.client
}

func Close(cli *redis.Client) {
	if err := cli.Close(); err != nil {
		gologger.Fatalf("failed close connection redis: %v", err)
	}
}
