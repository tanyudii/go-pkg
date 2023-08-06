package go_queue

import (
	"context"
	"fmt"
	"github.com/vmihailenco/taskq/v3"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Worker interface {
	GetTasks() []*taskq.TaskOptions
}

type Service interface {
	Shutdown(ctx context.Context) error
	RunGracefully(t int)
	GetQueue() taskq.Queue
	RegisterWorker(workers ...Worker)
	AddMessage(ctx context.Context, name string, args ...interface{}) error
	AddMessageRaw(msg *taskq.Message) error
	AddMessageWithSchedule(ctx context.Context, taskName string, schedule *time.Time, args ...interface{}) error
}

type service struct {
	cfg   *Config
	queue taskq.Queue
}

func NewService(f taskq.Factory, args ...ConfigFunc) Service {
	cfg := generate(args...)
	return &service{
		cfg:   cfg,
		queue: f.RegisterQueue(cfg.ToQueueOptions()),
	}
}

func (s *service) Shutdown(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (s *service) RunGracefully(t int) {
	mainCtx, cancelMainCtx := context.WithCancel(context.Background())
	go func() {
		if err := s.queue.Consumer().Start(mainCtx); err != nil {
			fmt.Printf("go queue consumer err: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Printf("go queue is shutting down: for %ds %v\n", t, time.Now())
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t)*time.Second)
	defer cancel()
	cancelMainCtx()
	if err := s.Shutdown(ctx); err != nil {
		fmt.Printf("go queue shutdown err: %v\n", err)
	}
	fmt.Printf("go queue shutdown completed: %v\n", time.Now())
}

func (s *service) GetQueue() taskq.Queue {
	return s.queue
}

func (s *service) RegisterWorker(workers ...Worker) {
	for _, worker := range workers {
		for _, task := range worker.GetTasks() {
			taskq.RegisterTask(task)
		}
	}
}

func (s *service) AddMessage(ctx context.Context, taskName string, args ...interface{}) error {
	msg := taskq.NewMessage(ctx, args...)
	msg.TaskName = taskName
	return s.AddMessageRaw(msg)
}

func (s *service) AddMessageWithSchedule(ctx context.Context, taskName string, schedule *time.Time, args ...interface{}) error {
	now := time.Now()
	msg := taskq.NewMessage(ctx, args...)
	msg.TaskName = taskName
	if schedule != nil && schedule.After(now) {
		msg.SetDelay(schedule.Sub(now))
	}
	return s.queue.Add(msg)
}

func (s *service) AddMessageRaw(msg *taskq.Message) error {
	return s.queue.Add(msg)
}
