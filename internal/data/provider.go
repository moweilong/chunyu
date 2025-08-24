package data

import (
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/glebarez/sqlite"
	"github.com/google/wire"
	"github.com/moweilong/chunyu/internal/config"
	"github.com/moweilong/chunyu/pkg/orm"
	"github.com/moweilong/chunyu/pkg/system"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(SetupDB)

// SetupDB 初始化数据存储
func SetupDB(c *config.Bootstrap, l *slog.Logger) (*gorm.DB, error) {
	cfg := c.Data.Database
	dial, isSQLite := getDialector(cfg.Dsn)
	if isSQLite {
		cfg.MaxIdleConns = 1
		cfg.MaxOpenConns = 1
	}
	db, err := orm.New(true, dial, orm.Config{
		MaxIdleConns:    int(cfg.MaxIdleConns),
		MaxOpenConns:    int(cfg.MaxOpenConns),
		ConnMaxLifetime: cfg.ConnMaxLifetime.Duration(),
		SlowThreshold:   cfg.SlowThreshold.Duration(),
	}, orm.NewLogger(l, c.Debug, cfg.SlowThreshold.Duration()))
	return db, err
}

// getDialector 返回 dial 和 是否 sqlite
func getDialector(dsn string) (gorm.Dialector, bool) {
	if strings.HasPrefix(dsn, "postgres") {
		return postgres.New(postgres.Config{
			DriverName: "pgx",
			DSN:        dsn,
		}), false
	}
	return sqlite.Open(filepath.Join(system.Getwd(), dsn)), true
}
