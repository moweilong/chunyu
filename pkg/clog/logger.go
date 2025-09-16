package clog

import (
	"context"
	"log"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"chunyu/pkg/clog/klog"
)

var (
	std = New(NewOptions())
	mu  sync.Mutex
)

// Init initializes logger with specified options.
func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()
	std = New(opts)
}

// New create logger by opts which can custmoized by command arguments.
func New(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	encodeLevel := zapcore.CapitalLevelEncoder
	// when output to local path, with color is forbidden
	if opts.Format == consoleFormat && opts.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     timeEncoder,
		EncodeDuration: milliSecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	loggerConfig := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Development:       opts.Development,
		DisableCaller:     opts.DisableCaller,
		DisableStacktrace: opts.DisableStacktrace,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         opts.Format,
		EncoderConfig:    encoderConfig,
		OutputPaths:      opts.OutputPaths,
		ErrorOutputPaths: opts.ErrorOutputPaths,
	}

	var err error
	l, err := loggerConfig.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	logger := &zapLogger{
		zapLogger: l.Named(opts.Name),
		infoLogger: infoLogger{
			log:   l,
			level: zap.InfoLevel,
		},
	}
	klog.InitLogger(l)
	zap.RedirectStdLog(l)

	return logger
}

// SugaredLogger returns global sugared logger.
func SugaredLogger() *zap.SugaredLogger {
	return std.zapLogger.Sugar()
}

// StdErrLogger returns logger of standard library which writes to supplied zap
// logger at error level.
func StdErrLogger() *log.Logger {
	if std == nil {
		return nil
	}
	if l, err := zap.NewStdLogAt(std.zapLogger, zapcore.ErrorLevel); err == nil {
		return l
	}

	return nil
}

// StdInfoLogger returns logger of standard library which writes to supplied zap
// logger at info level.
func StdInfoLogger() *log.Logger {
	if std == nil {
		return nil
	}
	if l, err := zap.NewStdLogAt(std.zapLogger, zapcore.InfoLevel); err == nil {
		return l
	}

	return nil
}

func (l *zapLogger) Write(p []byte) (n int, err error) {
	l.zapLogger.Info(string(p))

	return len(p), nil
}

// WithValues creates a child logger and adds adds Zap fields to it.
func WithValues(keysAndValues ...interface{}) Logger { return std.WithValues(keysAndValues...) }

func (l *zapLogger) WithValues(keysAndValues ...interface{}) Logger {
	newLogger := l.zapLogger.With(handleFields(l.zapLogger, keysAndValues)...)

	return NewLogger(newLogger)
}

// WithName adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
func WithName(s string) Logger { return std.WithName(s) }

func (l *zapLogger) WithName(name string) Logger {
	newLogger := l.zapLogger.Named(name)

	return NewLogger(newLogger)
}

// Flush calls the underlying Core's Sync method, flushing any buffered
// log entries. Applications should take care to call Sync before exiting.
func Flush() { std.Flush() }

func (l *zapLogger) Flush() {
	_ = l.zapLogger.Sync()
}

// ZapLogger used for other log wrapper such as klog.
func ZapLogger() *zap.Logger {
	return std.zapLogger
}

// CheckIntLevel used for other log wrapper such as klog which return if logging a
// message at the specified level is enabled.
func CheckIntLevel(level int32) bool {
	var lvl zapcore.Level
	if level < 5 {
		lvl = zapcore.InfoLevel
	} else {
		lvl = zapcore.DebugLevel
	}
	checkEntry := std.zapLogger.Check(lvl, "")

	return checkEntry != nil
}

// L method output with specified context value.
func L(ctx context.Context) *zapLogger {
	return std.L(ctx)
}

func (l *zapLogger) L(ctx context.Context) *zapLogger {
	lg := l.clone()

	if requestID := ctx.Value(KeyRequestID); requestID != nil {
		lg.zapLogger = lg.zapLogger.With(zap.Any(KeyRequestID, requestID))
	}
	if username := ctx.Value(KeyUsername); username != nil {
		lg.zapLogger = lg.zapLogger.With(zap.Any(KeyUsername, username))
	}
	if watcherName := ctx.Value(KeyWatcherName); watcherName != nil {
		lg.zapLogger = lg.zapLogger.With(zap.Any(KeyWatcherName, watcherName))
	}

	return lg
}

//nolint:predeclared
func (l *zapLogger) clone() *zapLogger {
	copy := *l

	return &copy
}

// NewLogger creates a new logr.Logger using the given Zap Logger to log.
func NewLogger(l *zap.Logger) Logger {
	return &zapLogger{
		zapLogger: l,
		infoLogger: infoLogger{
			log:   l,
			level: zap.InfoLevel,
		},
	}
}
