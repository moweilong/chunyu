package data

import (
	"log/slog"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/google/wire"
	"github.com/moweilong/chunyu/internal/config"
	"github.com/moweilong/chunyu/pkg/orm"
	"github.com/moweilong/chunyu/pkg/system"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(SetupDB)

// SetupDB 初始化数据存储
func SetupDB(c *config.Bootstrap, l *slog.Logger) (*gorm.DB, error) {
	cfg := c.Data.Database
	dsn := cfg.Dsn
	var dial gorm.Dialector
	switch cfg.DBType {
	case "postgres":
		dial = postgres.New(postgres.Config{
			DriverName: "pgx",
			DSN:        dsn,
		})
	case "mysql":
		dial = mysql.New(mysql.Config{
			DSN: dsn,
		})
	default:
		dial = sqlite.Open(filepath.Join(system.Getwd(), dsn))
	}

	if cfg.DBType == "sqlite" {
		cfg.MaxIdleConns = 1
		cfg.MaxOpenConns = 1
	}
	db, err := orm.New(true, dial, orm.Config{
		MaxIdleConns:    int(cfg.MaxIdleConns),
		MaxOpenConns:    int(cfg.MaxOpenConns),
		ConnMaxLifetime: cfg.ConnMaxLifetime.Duration(),
		SlowThreshold:   cfg.SlowThreshold.Duration(),
	}, orm.NewLogger(l, c.Debug, cfg.SlowThreshold.Duration()))
	slog.Info("数据库连接成功", "type", cfg.DBType, "dsn", cfg.Dsn)
	return db, err
}
