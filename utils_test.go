package router

import (
	"github.com/nats-io/nats-server/test"
	"github.com/nats-io/nats.go"
)

func RunServer(fn func(*nats.Conn)) {
	opts := test.DefaultTestOptions
	s := test.RunServer(&opts)
	nc, _ := nats.Connect("nats://127.0.0.1")
	fn(nc)
	nc.Close()
	s.Shutdown()
}
