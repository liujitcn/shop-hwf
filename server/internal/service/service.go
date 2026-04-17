package service

import (
	"gitee.com/liujit/shop/server/pkg/core"
	"gitee.com/liujit/shop/server/pkg/service/admin"
	"gitee.com/liujit/shop/server/pkg/service/app"
	"gitee.com/liujit/shop/server/pkg/service/config"
	"gitee.com/liujit/shop/server/pkg/service/file"
	"gitee.com/liujit/shop/server/pkg/service/login"
	"gitee.com/liujit/shop/server/pkg/service/pay"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	core.NewShopCore,
	admin.NewAuthServiceImpl,
	admin.NewBaseApiServiceImpl,
	admin.NewBaseConfigServiceImpl,
	admin.NewBaseDeptServiceImpl,
	admin.NewBaseDictServiceImpl,
	admin.NewBaseJobServiceImpl,
	admin.NewBaseLogServiceImpl,
	admin.NewBaseMenuServiceImpl,
	admin.NewBaseRoleServiceImpl,
	admin.NewBaseUserServiceImpl,

	admin.NewDashboardServiceImpl,

	admin.NewGoodsCategoryServiceImpl,
	admin.NewGoodsServiceImpl,
	admin.NewGoodsPropServiceImpl,
	admin.NewGoodsSkuServiceImpl,
	admin.NewGoodsSpecServiceImpl,

	admin.NewOrderServiceImpl,

	admin.NewPayBillServiceImpl,

	admin.NewShopBannerServiceImpl,
	admin.NewShopHotServiceImpl,
	admin.NewShopServiceServiceImpl,

	admin.NewUserStoreServiceImpl,

	app.NewAuthServiceImpl,
	app.NewBaseAreaServiceImpl,
	app.NewBaseDictServiceImpl,

	app.NewGoodsCategoryServiceImpl,
	app.NewGoodsServiceImpl,

	app.NewOrderServiceImpl,

	app.NewShopBannerServiceImpl,
	app.NewShopHotServiceImpl,
	app.NewShopServiceServiceImpl,

	app.NewUserAddressServiceImpl,
	app.NewUserCartServiceImpl,
	app.NewUserCollectServiceImpl,
	app.NewUserStoreServiceImpl,

	config.NewConfigServiceImpl,

	file.NewFileServiceImpl,

	login.NewLoginServiceImpl,

	pay.NewPayServiceImpl,
)
