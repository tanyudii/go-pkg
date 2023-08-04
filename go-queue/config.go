package go_queue

import "github.com/vmihailenco/taskq/v3"

type Config struct {
	name          string
	minNumWorker  int32
	maxNumWorker  int32
	workerLimit   int32
	maxNumFetcher int32
	redis         taskq.Redis
}

type ConfigFunc func(c *Config)

func Name(n string) ConfigFunc {
	return func(c *Config) {
		c.name = n
	}
}

func MinNumWorker(w int32) ConfigFunc {
	return func(c *Config) {
		c.minNumWorker = w
	}
}

func MaxNumWorker(w int32) ConfigFunc {
	return func(c *Config) {
		c.maxNumWorker = w
	}
}

func WorkerLimit(w int32) ConfigFunc {
	return func(c *Config) {
		c.workerLimit = w
	}
}

func MaxNumFetcher(f int32) ConfigFunc {
	return func(c *Config) {
		c.maxNumFetcher = f
	}
}

func Redis(r taskq.Redis) ConfigFunc {
	return func(c *Config) {
		c.redis = r
	}
}

func (c *Config) ToQueueOptions() *taskq.QueueOptions {
	return &taskq.QueueOptions{
		Name:          c.name,
		MinNumWorker:  c.minNumWorker,
		MaxNumWorker:  c.maxNumWorker,
		WorkerLimit:   c.workerLimit,
		MaxNumFetcher: c.maxNumFetcher,
		Redis:         c.redis,
	}
}

func generate(args ...ConfigFunc) *Config {
	c := &Config{}
	for i := range args {
		args[i](c)
	}
	return c
}
