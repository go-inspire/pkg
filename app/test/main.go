package main

import (
	"context"
	"inspire-pkg/app"
	"log"
	"time"
)

type server struct {
	name   string
	ctx    context.Context
	cancel context.CancelFunc
}

func NewServer(name string) *server {
	ctx, cancel := context.WithCancel(context.Background())
	return &server{
		name:   name,
		ctx:    ctx,
		cancel: cancel,
	}

}
func (s server) Start() error {
	log.Printf("%s started", s.name)

	for {
		select {
		case <-s.ctx.Done():
			log.Printf("%s done", s.name)
			return nil
		default:
			log.Printf("%s running", s.name)
			time.Sleep(time.Second)
		}
	}
}

func (s server) Stop() error {
	if s.cancel != nil {
		log.Printf("%s cancel\n", s.name)
		s.cancel()
	}
	return nil
}

func main() {
	app := app.New("test",
		app.Server(NewServer("server1"), NewServer("server2")),
	)

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}

}
