package clog

import (
	"context"
	"time"

	"go.uber.org/zap/zapcore"
)

// Defines common log fields.
const (
	KeyRequestID   string = "requestID"
	KeyUsername    string = "username"
	KeyWatcherName string = "watcher"
)

// Field is an alias for the field structure in the underlying log frame.
type Field = zapcore.Field

// Level is an alias for the level structure in the underlying log frame.
type Level = zapcore.Level

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func milliSecondsDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendFloat64(float64(d) / float64(time.Millisecond))
}

type key int

const (
	logContextKey key = iota
)

// WithContext returns a copy of context in which the log value is set.
func WithContext(ctx context.Context) context.Context {
	return std.WithContext(ctx)
}

func (l *zapLogger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logContextKey, l)
}

// FromContext returns the value of the log key on the ctx.
func FromContext(ctx context.Context) Logger {
	if ctx != nil {
		logger := ctx.Value(logContextKey)
		if logger != nil {
			return logger.(Logger)
		}
	}

	return WithName("Unknown-Context")
}
