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
func (r *Route) Start(ctx context.Context, nc *nats.Conn, buffsize uint) {
	ch := make(chan *nats.Msg, buffsize)
	nc.ChanSubscribe(r.Subject, ch)
	go ListenAndHandle(ctx, ch, r.Handler)
}

func ListenAndHandle(ctx context.Context, ch chan *nats.Msg, fn HandlerFunc) {
	func() {
		for {
			select {
			case msg := <-ch:
				go fn(ctx, msg)
			case <-ctx.Done():
				return
			}
		}
	}()
}
