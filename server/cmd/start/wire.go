//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package start

import (
	"gitee.com/liujit/shop/server/internal/biz"
	"gitee.com/liujit/shop/server/internal/config"
	"gitee.com/liujit/shop/server/internal/data"
	"gitee.com/liujit/shop/server/internal/server"
	"gitee.com/liujit/shop/server/internal/service"
	"go.newcapec.cn/ncttools/nmskit"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/conf"
	"go.newcapec.cn/ncttools/nmskit/registry"

	"go.newcapec.cn/ncttools/nmskit/log"

	"github.com/google/wire"
)

// initApp init NmsKit application.
func initApp(log.Logger, registry.Registrar, *conf.Bootstrap) (*nmskit.App, func(), error) {
	panic(wire.Build(server.ProviderSet, service.ProviderSet, data.ProviderSet, config.ProviderSet, biz.ProviderSet, newApp))
}
