package app

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/moweilong/chunyu/domain/version/versionapi"
	"github.com/moweilong/chunyu/internal/config"
	"github.com/moweilong/chunyu/pkg/logger"
	"github.com/moweilong/chunyu/pkg/server"
	"github.com/moweilong/chunyu/pkg/system"
)

func Run(bc *config.Bootstrap) {
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// 以可执行文件所在目录为工作目录，防止以服务方式运行时，工作目录切换到其它位置
	bin, _ := os.Executable()
	if err := os.Chdir(filepath.Dir(bin)); err != nil {
		slog.Error("change work dir fail", "err", err)
	}

	log, clean := SetupLog(bc)
	defer clean()

	// 如果需要执行表迁移，递增此版本号和表更新说明
	versionapi.DBVersion = "0.0.7"
	versionapi.DBRemark = "user表结构优化"

	handler, cleanUp, err := WireApp(bc, log)
	if err != nil {
		slog.Error("程序构建失败", "err", err)
		panic(err)
	}
	defer cleanUp()

	svc := server.New(handler,
		server.Port(strconv.Itoa(bc.Server.HTTP.Port)),
		server.ReadTimeout(bc.Server.HTTP.Timeout.Duration()),
		server.WriteTimeout(bc.Server.HTTP.Timeout.Duration()),
	)
	go svc.Start()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("服务启动成功 port:", bc.Server.HTTP.Port)

	select {
	case s := <-interrupt:
		slog.Info(`<-interrupt`, "signal", s.String())
	case err := <-svc.Notify():
		system.ErrPrintf("err: %s\n", err.Error())
		slog.Error(`<-server.Notify()`, "err", err)
	}
	if err := svc.Shutdown(); err != nil {
		slog.Error(`server.Shutdown()`, "err", err)
	}
}

// SetupLog 初始化日志
func SetupLog(bc *config.Bootstrap) (*slog.Logger, func()) {
	logDir := filepath.Join(system.Getwd(), bc.Log.Dir)
	return logger.SetupSlog(logger.Config{
		Dir:          logDir,                            // 日志地址
		Debug:        bc.Debug,                          // 服务级别Debug/Release
		MaxAge:       bc.Log.MaxAge.Duration(),          // 日志存储时间
		RotationTime: bc.Log.RotationTime.Duration(),    // 循环时间
		RotationSize: bc.Log.RotationSize * 1024 * 1024, // 循环大小
		Level:        bc.Log.Level,                      // 日志级别
	})
}
