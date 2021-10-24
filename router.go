package router

import (
	"context"

	"github.com/nats-io/nats.go"
)

type Router struct {
	routes []route
	nc     *nats.Conn
	ctx    context.Context
	buffsz uint
}

func New(ctx context.Context, nc *nats.Conn, buffersize uint) Router {
	return Router{nc: nc, buffsz: buffersize}
}

// Add a new route to the routing table and start listening on it
func (r *Router) Route(subject string, handler HandlerFunc) {
	route := route{
		subject: subject,
		handler: handler,
	}
	r.routes = append(r.routes, route)
	route.start(r.ctx, r.nc, r.buffsz)
}
