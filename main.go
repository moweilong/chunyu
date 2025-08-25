package main

import (
	"expvar"
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/moweilong/chunyu/internal/app"
	"github.com/moweilong/chunyu/internal/config"
	"github.com/moweilong/chunyu/pkg/system"
)

var (
	buildVersion = "0.0.1" // 构建版本号
	gitBranch    = "dev"   // git 分支
	gitHash      = "debug" // git 提交点哈希值
	release      string    // 发布模式 true/false
	buildTime    string    // 构建时间戳
)

// 自定义配置目录
// TODO 配置文件路径参数无效
var configDir = flag.String("conf", "./configs", "config directory, eg: -conf /configs/")

func getBuildRelease() bool {
	v, _ := strconv.ParseBool(release)
	return v
}

func main() {
	flag.Parse()

	// 初始化配置
	var bc config.Bootstrap
	fileDir, _ := system.Abs(*configDir)
	_ = os.MkdirAll(fileDir, 0o755)
	filePath := filepath.Join(fileDir, "config.toml")
	configIsNotExistWrite(filePath)
	if err := config.SetupConfig(&bc, filePath); err != nil {
		panic(err)
	}
	bc.Debug = !getBuildRelease()
	bc.BuildVersion = buildVersion
	bc.ConfigDir = fileDir
	bc.ConfigPath = filePath

	{
		expvar.NewString("version").Set(buildVersion)
		expvar.NewString("git_branch").Set(gitBranch)
		expvar.NewString("git_hash").Set(gitHash)
		expvar.NewString("build_time").Set(buildTime)
		expvar.Publish("timestamp", expvar.Func(func() any {
			return time.Now().Format(time.DateTime)
		}))
	}

	app.Run(&bc)
}

// configIsNotExistWrite 配置文件不存在时，回写配置
func configIsNotExistWrite(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := config.WriteConfig(config.DefaultConfig(), path); err != nil {
			system.ErrPrintf("WriteConfig", "err", err)
		}
	}
}
