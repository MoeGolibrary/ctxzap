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

// ToContext put logger to context
func ToContext(ctx context.Context, logger *zap.Logger, fields []zap.Field) context.Context {
	return context.WithValue(ctx, zapCtxKey, &ctxLogger{logger: logger, fields: fields})
}

// AddFields attach more fields to the context
func AddFields(ctx context.Context, fields ...zap.Field) {
	if logger, ok := ctx.Value(zapCtxKey).(*ctxLogger); ok && logger != nil {
		logger.fields = append(logger.fields, fields...)
	}
}

// ReplaceField replace existing field or append new field
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

// Debug log
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	if logger, ok := ctx.Value(zapCtxKey).(*ctxLogger); ok && logger != nil {
		logger.logger.Debug(msg, append(logger.fields, fields...)...)
	}
}

// Info log
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	if logger, ok := ctx.Value(zapCtxKey).(*ctxLogger); ok && logger != nil {
		logger.logger.Info(msg, append(logger.fields, fields...)...)
	}
}

// Warn log
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	if logger, ok := ctx.Value(zapCtxKey).(*ctxLogger); ok && logger != nil {
		logger.logger.Warn(msg, append(logger.fields, fields...)...)
	}
}

// Error log
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	if logger, ok := ctx.Value(zapCtxKey).(*ctxLogger); ok && logger != nil {
		logger.logger.Error(msg, append(logger.fields, fields...)...)
	}
}
