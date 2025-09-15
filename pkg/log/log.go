package log

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
)

// SetupSlog 初始化 slog
func SetupSlog(opts *Options) func() {
	// 应用 zap 日志等级
	SetLevel(opts.Level)
	// 采样率默认配置
	if opts.TickSec <= 0 {
		opts.TickSec = 1
		opts.First = 5
		opts.Thereafter = 5
	}
	// 日志轮替
	r := rotatelog(opts)
	// slog 初始化
	log := slog.New(
		newSlog(
			NewJSONLogger(opts.Debug, r, opts).Core(),
			zapslog.WithCaller(opts.Debug),
		),
	)

	if opts.ID != "" {
		log = log.With("serviceID", opts.ID)
	}
	if opts.Version != "" {
		log = log.With("serviceVersion", opts.Version)
	}
	slog.SetDefault(log)

	file, err := os.OpenFile(filepath.Join(opts.Dir, "crash.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err == nil {
		_ = SetCrashOutput(file)
	}
	return func() {
		if file != nil {
			file.Close()
		}
	}
}

// Level 日志级别
var Level = zap.NewAtomicLevelAt(zap.InfoLevel)

// SetLevel 设置日志级别 debug/warn/error
func SetLevel(l string) {
	switch strings.ToLower(l) {
	case "debug":
		Level.SetLevel(zap.DebugLevel)
	case "warn":
		Level.SetLevel(zap.WarnLevel)
	case "error":
		Level.SetLevel(zap.ErrorLevel)
	default:
		Level.SetLevel(zap.InfoLevel)
	}
}

func rotatelog(opt *Options) *rotatelogs.RotateLogs {
	if opt.MaxAge <= 0 {
		opt.MaxAge = 7 * 24 * time.Hour
	}
	if opt.RotationTime <= 0 {
		opt.RotationTime = 12 * time.Hour
	}
	if opt.RotationSize <= 0 {
		opt.RotationSize = 10 * 1024 * 1024
	}
	r, _ := rotatelogs.New(
		filepath.Join(opt.Dir, "%Y%m%d_%H_%M_%S.log"),
		rotatelogs.WithMaxAge(opt.MaxAge),
		rotatelogs.WithRotationTime(opt.RotationTime),
		rotatelogs.WithRotationSize(opt.RotationSize),
	)
	return r
}

// NewJSONLogger 创建JSON日志
func NewJSONLogger(debug bool, w io.Writer, opt *Options) *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	config.NameKey = ""
	mulitWriteSyncer := []zapcore.WriteSyncer{
		zapcore.AddSync(w),
	}
	if debug {
		mulitWriteSyncer = append(mulitWriteSyncer, zapcore.AddSync(os.Stdout))
	}
	// level = zap.ErrorLevel
	core := zapcore.NewSamplerWithOptions(zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.NewMultiWriteSyncer(mulitWriteSyncer...),
		Level,
	), time.Duration(opt.TickSec)*time.Second, opt.First, opt.Thereafter)
	return zap.New(core, zap.AddCaller())
}

// SetCrashOutput 设置一个文件作为程序崩溃信息的输出目标，当程序发生未被 recover 的 panic 或致命错误时，
// 系统会自动将崩溃信息（包括堆栈跟踪）写入到指定文件中。
// 函数依赖于 Go 1.23 版本中新增的 runtime/debug 包的 SetCrashOutput 功能，这是 Go 语言在错误处理机制上的一个重要改进。
func SetCrashOutput(f *os.File) error {
	return debug.SetCrashOutput(f, debug.CrashOptions{})
}
