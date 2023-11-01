package grpc_zap

import (
	"context"
	"github.com/MoeGolibrary/ctxzap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

func invoke(ctx context.Context, logger *zap.Logger, o *options, method string, handler func(ctx context.Context) error) {
	fields := []zap.Field{
		zap.String("grpc.method", method),
	}
	if len(o.metadataFields) > 0 {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			for _, k := range o.metadataFields {
				fields = append(fields, zap.String("grpc."+k, strings.Join(md[k], ", ")))
			}
		}
	}
	if o.fieldsFunc != nil {
		fields = append(fields, o.fieldsFunc(ctx)...)
	}
	ctx = ctxzap.ToContext(ctx, logger, fields)
	if o.started {
		ctxzap.Info(ctx, "method invoke started", zap.String("grpc.event", "start"))
	}
	start := time.Now()
	err := handler(ctx)
	if o.finished {
		code := status.Code(err)
		finishedFields := []zap.Field{
			zap.String("grpc.event", "finish"),
			zap.Uint32("grpc.status", uint32(code)),
			zap.String("grpc.statusText", code.String()),
			zap.Duration("grpc.duration", time.Since(start)),
		}
		if err != nil {
			finishedFields = append(finishedFields, zap.Error(err))
			if o.errFieldsFunc != nil {
				finishedFields = append(finishedFields, o.errFieldsFunc(ctx, err)...)
			}
		}
		ctxzap.Info(ctx, "method invoke finished", finishedFields...)
	}
}

func UnaryServerInterceptor(logger *zap.Logger, opts ...Option) grpc.UnaryServerInterceptor {
	o := resolveOptions(&options{started: false, finished: true}, opts)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		invoke(ctx, logger, o, info.FullMethod, func(ctx context.Context) error {
			resp, err = handler(ctx, req)
			return err
		})
		return
	}
}

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
	*options
	sent     int
	received int
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

func (w *wrappedServerStream) SendMsg(m any) error {
	err := w.ServerStream.SendMsg(m)
	if w.recv {
		w.sent += 1
		fields := []zap.Field{
			zap.String("grpc.event", "send"),
			zap.Int("grpc.sent", w.sent),
			zap.Int("grpc.received", w.received),
		}
		if err != nil {
			fields = append(fields, zap.Error(err))
			if w.errFieldsFunc != nil {
				fields = append(fields, w.errFieldsFunc(w.ctx, err)...)
			}
		}
		ctxzap.Info(w.ctx, "message sent", fields...)
	}
	return err
}

func (w *wrappedServerStream) RecvMsg(m any) error {
	err := w.ServerStream.RecvMsg(m)
	if w.recv {
		w.received += 1
		fields := []zap.Field{
			zap.String("grpc.event", "receive"),
			zap.Int("grpc.sent", w.sent),
			zap.Int("grpc.received", w.received),
		}
		if err != nil {
			fields = append(fields, zap.Error(err))
			if w.errFieldsFunc != nil {
				fields = append(fields, w.errFieldsFunc(w.ctx, err)...)
			}
		}
		ctxzap.Info(w.ctx, "message received", fields...)
	}
	return err
}

func StreamServerInterceptor(logger *zap.Logger, opts ...Option) grpc.StreamServerInterceptor {
	o := resolveOptions(&options{started: true, finished: true}, opts)
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		invoke(ss.Context(), logger, o, info.FullMethod, func(ctx context.Context) error {
			w := &wrappedServerStream{
				ServerStream: ss,
				ctx:          ctx,
				options:      o,
			}
			err = handler(srv, w)
			return err
		})
		return
	}
}
