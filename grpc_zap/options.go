package grpc_zap

import (
	"context"
	"go.uber.org/zap"
	"strings"
)

type options struct {
	metadataFields []string
	fieldsFunc     func(ctx context.Context) []zap.Field
	errFieldsFunc  func(ctx context.Context, err error) []zap.Field
	started        bool
	finished       bool
	send           bool
	recv           bool
}

type Option interface {
	apply(o *options)
}

type optFunc func(o *options)

func (f optFunc) apply(o *options) {
	f(o)
}

// WithMetadataFields attach all the specified keys to all the logs
// please note the values are extracted synchronously.
func WithMetadataFields(fields ...string) Option {
	return optFunc(func(o *options) {
		for _, k := range fields {
			o.metadataFields = append(o.metadataFields, strings.ToLower(k))
		}
	})
}

// WithFieldsFunc attach custom zap fields to all the logs
// please note the values are extracted synchronously.
func WithFieldsFunc(fn func(ctx context.Context) []zap.Field) Option {
	return optFunc(func(o *options) {
		o.fieldsFunc = fn
	})
}

// WithAccessLog enable/disable access log
// both started & finished log are enabled for stream server by default
// only finished log is enabled for unary server by default
func WithAccessLog(started, finished bool) Option {
	return optFunc(func(o *options) {
		o.started = started
		o.finished = finished
	})
}

// WithEventLog enable/disable event log
// only works with stream interceptor
// both send/recv are disabled by default
func WithEventLog(send, recv bool) Option {
	return optFunc(func(o *options) {
		o.send = send
		o.recv = recv
	})
}

// WithErrorFieldsFunc extract error detail fields
func WithErrorFieldsFunc(fn func(ctx context.Context, err error) []zap.Field) Option {
	return optFunc(func(o *options) {
		o.errFieldsFunc = fn
	})
}

func resolveOptions(o *options, opts []Option) *options {
	for _, opt := range opts {
		opt.apply(o)
	}
	return o
}
