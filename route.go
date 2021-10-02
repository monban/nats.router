package router

import (
	"context"

	"github.com/nats-io/nats.go"
)

type HandlerFunc func(context.Context, *nats.Msg)

type Route struct {
	Subject string
	Handler HandlerFunc
}

// Start a goroutine that will call Handler whenever a matching Subject arrives
func (r *Route) Start(ctx context.Context, nc *nats.Conn) {
	ch := make(chan *nats.Msg, 100)
	nc.ChanSubscribe(r.Subject, ch)
	go func() {
		for {
			select {
			case msg := <-ch:
				go r.Handler(ctx, msg)
			case <-ctx.Done():
				return
			}
		}
	}()
}
