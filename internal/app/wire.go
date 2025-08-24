//go:build wireinject
// +build wireinject

package app

import (
	"log/slog"
	"net/http"

	"github.com/google/wire"
	"github.com/moweilong/chunyu/internal/config"
	"github.com/moweilong/chunyu/internal/data"
	"github.com/moweilong/chunyu/internal/web/api"
)

func WireApp(bc *config.Bootstrap, log *slog.Logger) (http.Handler, func(), error) {
	panic(wire.Build(data.ProviderSet, api.ProviderVersionSet, api.ProviderSet))
}
