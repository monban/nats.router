package router

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func TestSingleRoute(t *testing.T) {
	ctx, done := context.WithTimeout(context.Background(), 5*time.Second)
	subject := "FOO"
	var data []byte
	called := make(chan bool)
	fooHandler := func(ctx context.Context, msg *nats.Msg) {
		called <- true
	}
	r := &Router{Routes: []Route{
		{subject, fooHandler},
	}}
	RunServer(func(nc *nats.Conn) {
		go r.ListenAndHandle(ctx, nc)
		r.WaitUntilReady()
		nc.Publish(subject, data)
		select {
		case <-called:
			t.Log("handler called")
		case <-ctx.Done():
			t.Error(ctx.Err())
		}
		done()
	})
}

func TestMultipleRoutes(t *testing.T) {
	ctx, done := context.WithTimeout(context.Background(), 5*time.Second)
	var counter int32
	var data []byte
	var expected_count int32 = 4
	h := HandlerCounterFunc(&counter, t, "Handler")
	r := Router{Routes: []Route{
		{"FOO", h},
		{"BAR", h},
		{"BAZ", h},
	}}

	RunServer(func(nc *nats.Conn) {
		go r.ListenAndHandle(ctx, nc)
		r.WaitUntilReady()
		nc.Publish("BAZ", data)
		nc.Publish("BAR", data)
		nc.Publish("BAR", data)
		nc.Publish("FOO", data)
		for counter != expected_count {
			select {
			case <-ctx.Done():
				t.Error(ctx.Err())
				return
			default:
				runtime.Gosched()
			}
		}
		done()
	})
}

func TestHierarchicalRoutes(t *testing.T) {
	ctx, done := context.WithTimeout(context.Background(), 1*time.Second)
	var catchAllCounter int32
	var fooSubCounter int32
	var fooBarCounter int32
	var fooBarStarCounter int32
	var data []byte
	r := Router{[]Route{
		{">", HandlerCounterFunc(&catchAllCounter, t, ">")},
		{"FOO.>", HandlerCounterFunc(&fooSubCounter, t, ">")},
		{"FOO.BAR", HandlerCounterFunc(&fooBarCounter, t, ">")},
		{"FOO.BAR.*", HandlerCounterFunc(&fooBarStarCounter, t, ">")},
	}, false}

	RunServer(func(nc *nats.Conn) {
		go r.ListenAndHandle(ctx, nc)
		r.WaitUntilReady()
		nc.Publish("FOO", data)
		nc.Publish("FOO.BAR", data)
		nc.Publish("FOO.BAR.BAZ", data)
		nc.Publish("FOO.QUX", data)
		for !(catchAllCounter == 4 && fooSubCounter == 3 && fooBarCounter == 1 && fooBarStarCounter == 1) {
			select {
			case <-ctx.Done():
				t.Error(ctx.Err())
				t.Errorf("catchAllCounter: %v", catchAllCounter)
				t.Errorf("fooSubCounter: %v", fooSubCounter)
				t.Errorf("fooBarCounter: %v", fooBarCounter)
				t.Errorf("fooBarStarCounter: %v", fooBarStarCounter)
				return
			default:
				runtime.Gosched()
			}
		}
		done()
	})
}
