package router

import (
	"context"

	"github.com/nats-io/nats.go"
)

type HandlerFunc func(context.Context, *nats.Msg)

type route struct {
	subject string
	handler HandlerFunc
	ctx     context.Context
}

// Start a goroutine that will call Handler whenever a matching Subject arrives
func (r *route) start(ctx context.Context, nc *nats.Conn, buffsize uint) {
	ch := make(chan *nats.Msg, buffsize)
	nc.ChanSubscribe(r.subject, ch)
	go listenAndHandle(ctx, ch, r.handler)
}

func listenAndHandle(ctx context.Context, ch chan *nats.Msg, fn HandlerFunc) {
	func() {
		for {
			select {
			case msg := <-ch:
				go fn(ctx, msg)
				// TODO method of removing routes
				// TODO below code causes SIGSEGV
				//case <-ctx.Done():
				//return
			}
		}
	}()
}
