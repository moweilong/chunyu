# log

## GormLogger

GormLogger 是一个 GORM Logger 实现，用于将 GORM 日志集成到 slog 日志系统中。

### 功能

- 支持将 GORM 日志输出到 slog 日志系统中
- 支持自定义日志级别和慢查询阈值
- 支持忽略 RecordNotFound 错误

### 安装

```bash
go get github.com/moweilong/chunyu/pkg/log
```

### 使用

```go
import (
    "log/slog"
    "time"

    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// 创建slog日志器
logger := log.SetupSlog()

// 创建GORM日志器（启用debug模式，慢查询阈值为200ms）
gormLogger := log.NewGormLogger(logger, true, 200*time.Millisecond)

// 在初始化GORM DB时设置日志器
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: gormLogger,
})

// 也可以通过 LogMode 方法动态调整日志级别
db.Logger.LogMode(logger.Info) // 启用详细日志
// 或
db.Logger.LogMode(logger.Error) // 仅记录错误日志
```

### 配置

- `logger`: slog 日志器实例
- `debug`: 是否启用调试模式（true 表示输出详细日志，false 表示仅输出警告和错误）
- `slowThreshold`: 慢查询阈值，单位为毫秒

### 注意事项

- 确保 slog 日志器已正确配置，否则 GORM 日志将无法输出
- 建议在生产环境中设置 `debug` 为 false，避免输出过多日志
- 建议在开发环境中设置 `debug` 为 true，方便调试
- 警告和错误级别的日志始终会被记录，无论 `debug` 模式如何设置