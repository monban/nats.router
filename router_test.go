package router

import (
	"bytes"
	"fmt"
	"sync"
	"testing"

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
