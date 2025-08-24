package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/moweilong/chunyu/domain/uniqueid"
	"github.com/moweilong/chunyu/domain/uniqueid/store/uniqueiddb"
	"github.com/moweilong/chunyu/domain/version/versionapi"
	"github.com/moweilong/chunyu/internal/config"
	"github.com/moweilong/chunyu/pkg/orm"
	"github.com/moweilong/chunyu/pkg/web"
	"gorm.io/gorm"
)

var (
	ProviderVersionSet = wire.NewSet(versionapi.NewVersionCore)
	ProviderSet        = wire.NewSet(
		wire.Struct(new(Usecase), "*"),
		NewHTTPHandler,
		versionapi.New,
	)
)

type Usecase struct {
	Conf    *config.Bootstrap
	DB      *gorm.DB
	Version versionapi.API
}

// NewHTTPHandler 生成Gin框架路由内容
func NewHTTPHandler(uc *Usecase) http.Handler {
	cfg := uc.Conf
	// 检查是否设置了 JWT 密钥，如果未设置，则生成一个长度为 32 的随机字符串作为密钥
	if cfg.Server.HTTP.JwtSecret == "" {
		uc.Conf.Server.HTTP.JwtSecret = orm.GenerateRandomString(32)
	}
	// 如果不处于调试模式，将 Gin 设置为发布模式
	if !uc.Conf.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	g := gin.New()
	// 处理未找到路由的情况，返回 JSON 格式的 404 错误信息
	g.NoRoute(func(c *gin.Context) {
		c.JSON(404, "来到了无人的荒漠")
	})
	// 如果启用了 Pprof，设置 Pprof 监控
	if cfg.Server.HTTP.PProf.Enabled {
		web.SetupPProf(g, &cfg.Server.HTTP.PProf.AccessIps)
	}

	setupRouter(g, uc)
	return g
}

// NewUniqueID 生成唯一 id
func NewUniqueID(db *gorm.DB) uniqueid.Core {
	store := uniqueiddb.NewDB(db).AutoMigrate(orm.GetEnabledAutoMigrate())
	return uniqueid.NewCore(store, 6)
}
