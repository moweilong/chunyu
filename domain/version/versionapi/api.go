package versionapi

import (
	"log/slog"

	"github.com/moweilong/chunyu/domain/version"
	"github.com/moweilong/chunyu/domain/version/store/versiondb"
	"github.com/moweilong/chunyu/pkg/orm"
	"gorm.io/gorm"
)

// 通过修改版本号，来控制是否执行表迁移
var (
	DBVersion = "0.0.1"
	DBRemark  = "debug"
)

// NewVersionCore ...
func NewVersionCore(db *gorm.DB) version.Core {
	vdb := versiondb.NewDB(db)
	core := version.NewCore(vdb)
	isOK := core.IsAutoMigrate(DBVersion)
	vdb.AutoMigrate(isOK)
	if isOK {
		slog.Info("更新数据库表结构")
		if err := core.RecordVersion(DBVersion, DBRemark); err != nil {
			slog.Error("RecordVersion", "err", err)
		}
	}
	orm.SetEnabledAutoMigrate(isOK)
	return core
}
