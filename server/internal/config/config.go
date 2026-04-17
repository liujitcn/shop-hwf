package config

import (
	"gitee.com/liujit/shop/server/pkg/config"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(config.ParseShopConfig, config.ParseWxMiniApp, config.ParseWxPay)
