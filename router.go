package router

import (
	"context"
	"runtime"

	"github.com/nats-io/nats.go"
)

type Router struct {
	Routes []Route
	ready  bool
}

func (r *Router) ListenAndHandle(ctx context.Context, nc *nats.Conn) {
	r.ready = false
	for _, route := range r.Routes {
		route.Start(ctx, nc)
	}
	r.ready = true
	<-ctx.Done()
}

// Blocks until the router has loaded all routes
func (r *Router) WaitUntilReady() {
	for !r.ready {
		runtime.Gosched()
	}
}
