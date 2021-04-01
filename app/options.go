package app

import (
	"context"
	"inspire-pkg/app/server"
	"os"
)

// Option is an application option.
type Option func(o *options)

// options is an application options.
type options struct {
	id       string
	name     string
	version  string
	metadata map[string]string

	ctx  context.Context
	sigs []os.Signal

	servers []server.Server
}

// ID with service id.
func ID(id string) Option {
	return func(o *options) { o.id = id }
}

// Name with service name.
func Name(name string) Option {
	return func(o *options) { o.name = name }
}

// Version with service version.
func Version(version string) Option {
	return func(o *options) { o.version = version }
}

// Metadata with service metadata.
func Metadata(md map[string]string) Option {
	return func(o *options) { o.metadata = md }
}

// Context with service context.
func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// Signal with exit signals.
func Signal(sigs ...os.Signal) Option {
	return func(o *options) { o.sigs = sigs }
}

// Server with transport servers.
func Server(srv ...server.Server) Option {
	return func(o *options) { o.servers = srv }
}
