package router

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
)

func RunServer(fn func(*nats.Conn)) {
	s := test.RunRandClientPortServer()
	nc, _ := nats.Connect(s.ClientURL())
	fn(nc)
	nc.Close()
	s.Shutdown()
}

func HandlerCounterFunc(i *int32, t *testing.T, name string) HandlerFunc {
	return func(ctx context.Context, msg *nats.Msg) {
		atomic.AddInt32(i, 1)
		msg.Respond([]byte(""))
	}
}

// Make sure our testing infrastructure is valid
func TestInfrastructure(t *testing.T) {
	subject := "HELLO"
	data := []byte("World!")
	wg := sync.WaitGroup{}
	wg.Add(1)
	RunServer(func(nc *nats.Conn) {
		nc.Subscribe(">", func(msg *nats.Msg) {
			l := fmt.Sprintf("Expected: {Subject: %v, Data: %v}, Received: {Subject: %v, Data: %v}", subject, string(data), msg.Subject, string(msg.Data))
			if subject == msg.Subject && bytes.Equal(data, msg.Data) {
				t.Log(l)
			} else {
				t.Error(l)
			}
			wg.Done()
		})
		nc.Publish(subject, data)
		wg.Wait()
	})
}
