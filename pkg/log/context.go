package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ctxKeyType string

const (
	ctxLogFieldKey ctxKeyType = "ctx_log_fields"
	ctxLogKey      ctxKeyType = "ctx_log"
)

type EntryField struct {
	key   string
	value interface{}
}

func FromContext(ctx context.Context) *Logger {
	if ctx == nil {
		return NewLogger()
	}

	logger, ok := ctx.Value(ctxLogKey).(*Logger)
	if !ok {
		return NewLogger().WithContext(ctx)
	}

	return logger
}

func WithContext(ctx context.Context) *Logger {
	return NewLogger().WithContext(ctx)
}

func With(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, ctxLogKey, logger.WithContext(ctx))
}

func WithField(ctx context.Context, key string, value interface{}) context.Context {
	return WithFields(ctx, Field(key, value))
}

func WithFields(ctx context.Context, fields ...EntryField) context.Context {
	nextValues := append(contextFields(ctx), fields...)
	return context.WithValue(ctx, ctxLogFieldKey, nextValues)
}

func contextFields(ctx context.Context) []EntryField {
	ctxValues := ctx.Value(ctxLogFieldKey)
	if ctxValues == nil {
		return []EntryField{}
	} else if _, ok := ctxValues.([]EntryField); ok {
		return ctxValues.([]EntryField)
	}

	return []EntryField{}
}

func contextLogFields(ctx context.Context) logrus.Fields {
	fields := logrus.Fields{}
	if ctx == nil {
		return fields
	}

	ctxValues := ctx.Value(ctxLogFieldKey)
	if values, ok := ctxValues.([]EntryField); ok {
		for _, v := range values {
			fields[v.key] = v.value
		}
	}

	return fields
}

func Field(key string, value interface{}) EntryField {
	return EntryField{key, value}
}
