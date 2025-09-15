package log

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"gorm.io/gorm/logger"
)

// GormLogger 实现 gorm.io/gorm/logger.Interface 接口
// 用于将 GORM 日志集成到 slog 日志系统
type GormLogger struct {
	Logger                    *slog.Logger
	SlowThreshold             time.Duration
	debug                     bool
	IgnoreRecordNotFoundError bool
}

// NewGormLogger 创建一个新的 GORM Logger 实例
func NewGormLogger(log *slog.Logger, debug bool, slowThreshold time.Duration) *GormLogger {
	if slowThreshold <= 0 {
		slowThreshold = 200 * time.Millisecond // 默认慢查询阈值
	}

	return &GormLogger{
		Logger:                    log,
		SlowThreshold:             slowThreshold,
		debug:                     debug,
		IgnoreRecordNotFoundError: true,
	}
}

// LogMode 设置日志级别，实现 gorm.io/gorm/logger.Interface 接口
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	// 根据 level 设置 debug 模式
	newLogger.debug = level >= logger.Info
	return &newLogger
}

// Info 记录 info 级别的日志，实现 gorm.io/gorm/logger.Interface 接口
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.debug {
		l.Logger.InfoContext(ctx, "gorm info", "message", fmt.Sprintf(msg, data...))
	}
}

// Warn 记录 warn 级别的日志，实现 gorm.io/gorm/logger.Interface 接口
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	// 警告级别的日志总是记录
	l.Logger.WarnContext(ctx, "gorm warn", "message", fmt.Sprintf(msg, data...))
}

// Error 记录 error 级别的日志，实现 gorm.io/gorm/logger.Interface 接口
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	// 错误级别的日志总是记录
	l.Logger.ErrorContext(ctx, "gorm error", "message", fmt.Sprintf(msg, data...))
}

// Trace 记录 SQL 执行的详细信息，实现 gorm.io/gorm/logger.Interface 接口
// 包括 SQL 语句、执行时间、影响行数和错误信息
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	// 处理错误
	if err != nil {
		// 如果是记录未找到的错误且配置为忽略此类错误，则不记录
		if l.IgnoreRecordNotFoundError && strings.Contains(err.Error(), "record not found") {
			return
		}
		// 记录错误日志
		l.Logger.ErrorContext(ctx, "gorm sql error",
			"error", err,
			"sql", sql,
			"elapsed", elapsed,
			"rows", rows,
		)
		return
	}

	// 记录慢查询
	if l.SlowThreshold > 0 && elapsed > l.SlowThreshold {
		l.Logger.WarnContext(ctx, "gorm slow sql",
			"sql", sql,
			"elapsed", elapsed,
			"rows", rows,
		)
		return
	}

	// 记录普通 SQL 日志（仅在 debug 模式下）
	if l.debug {
		l.Logger.DebugContext(ctx, "gorm sql",
			"sql", sql,
			"elapsed", elapsed,
			"rows", rows,
		)
	}
}
