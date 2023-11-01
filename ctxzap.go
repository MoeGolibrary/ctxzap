package ctxzap

import (
	"context"

	"go.uber.org/zap"
)

type zapCtx struct{}

type ctxLogger struct {
	logger *zap.Logger
	fields []zap.Field
}

var (
	zapCtxKey = zapCtx{}
)

// FromContext extract logger from context
// Please use this with caution, this method will create a new logger,
// leveled logging functions like Info, Warn is preferred.
func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(zapCtxKey).(*ctxLogger); ok && logger != nil {
		return logger.logger.With(logger.fields...)
	}
	return nil
}

func ToContext(ctx context.Context, logger *zap.Logger, fields []zap.Field) context.Context {
	return context.WithValue(ctx, zapCtxKey, &ctxLogger{logger: logger, fields: fields})
}

func AddFields(ctx context.Context, fields ...zap.Field) {
	if logger, ok := ctx.Value(zapCtxKey).(*ctxLogger); ok && logger != nil {
		logger.fields = append(logger.fields, fields...)
	}
}

func ReplaceField(ctx context.Context, field zap.Field) {
	if logger, ok := ctx.Value(zapCtxKey).(*ctxLogger); ok && logger != nil {
		for i, f := range logger.fields {
			if f.Key == field.Key {
				logger.fields[i] = field
				return
			}
		}
		logger.fields = append(logger.fields, field)
	}
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	if logger, ok := ctx.Value(zapCtxKey).(*ctxLogger); ok && logger != nil {
		logger.logger.Debug(msg, append(logger.fields, fields...)...)
	}
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	if logger, ok := ctx.Value(zapCtxKey).(*ctxLogger); ok && logger != nil {
		logger.logger.Info(msg, append(logger.fields, fields...)...)
	}
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	if logger, ok := ctx.Value(zapCtxKey).(*ctxLogger); ok && logger != nil {
		logger.logger.Warn(msg, append(logger.fields, fields...)...)
	}
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	if logger, ok := ctx.Value(zapCtxKey).(*ctxLogger); ok && logger != nil {
		logger.logger.Error(msg, append(logger.fields, fields...)...)
	}
}
