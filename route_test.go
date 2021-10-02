package router

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func TestRoute(t *testing.T) {
	ctx, done := context.WithTimeout(context.Background(), 5*time.Second)
	subject := "FOO"
	data := []byte("bar")
	called := make(chan bool, 10)
	route := Route{Subject: subject, Handler: func(ctx context.Context, msg *nats.Msg) {
		t.Logf("Received message with subject %v data %v", msg.Subject, string(msg.Data))
		called <- true
	}}
	RunServer(func(nc *nats.Conn) {
		route.Start(ctx, nc)
		nc.Publish(subject, data)
		select {
		case <-called:
			t.Log("handler was called")
		case <-ctx.Done():
			t.Error(ctx.Err())
		}
		done()
	})
}
