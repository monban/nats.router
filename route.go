package router

import (
	"context"

	"github.com/nats-io/nats.go"
)

type Route struct {
	Subject string
	Handler func(context.Context, *nats.Msg)
}

func (r *Route) Start(ctx context.Context, nc *nats.Conn) {
	ch := make(chan *nats.Msg)
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
