package server

import (
	adminApi "gitee.com/liujit/shop/server/api/admin"
	appApi "gitee.com/liujit/shop/server/api/app"
	configApi "gitee.com/liujit/shop/server/api/config"
	fileApi "gitee.com/liujit/shop/server/api/file"
	loginApi "gitee.com/liujit/shop/server/api/login"
	payApi "gitee.com/liujit/shop/server/api/pay"
	"gitee.com/liujit/shop/server/pkg/service/admin"
	"gitee.com/liujit/shop/server/pkg/service/app"
	"gitee.com/liujit/shop/server/pkg/service/config"
	"gitee.com/liujit/shop/server/pkg/service/file"
	"gitee.com/liujit/shop/server/pkg/service/login"
	"gitee.com/liujit/shop/server/pkg/service/pay"
	authnEngine "go.newcapec.cn/ncttools/nmskit-auth/authn/engine"
	authzEngine "go.newcapec.cn/ncttools/nmskit-auth/authz/engine"
	nmskitAuthData "go.newcapec.cn/ncttools/nmskit-auth/data"
	bootstrapConf "go.newcapec.cn/ncttools/nmskit-bootstrap/conf"
	bootstrapServer "go.newcapec.cn/ncttools/nmskit-bootstrap/server"
	"go.newcapec.cn/ncttools/nmskit/log"
	"go.newcapec.cn/ncttools/nmskit/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(cfg *bootstrapConf.Server_GRPC,
	jwt *bootstrapConf.Authentication_Jwt,
	authenticator authnEngine.Authenticator,
	authorizer authzEngine.Engine,
	userToken *nmskitAuthData.UserToken,
	logger log.Logger,
	adminAuth *admin.AuthServiceImpl,
	adminBaseApi *admin.BaseApiServiceImpl,
	adminBaseConfig *admin.BaseConfigServiceImpl,
	adminBaseDept *admin.BaseDeptServiceImpl,
	adminBaseDict *admin.BaseDictServiceImpl,
	adminBaseJob *admin.BaseJobServiceImpl,
	adminBaseLog *admin.BaseLogServiceImpl,
	adminBaseMenu *admin.BaseMenuServiceImpl,
	adminBaseRole *admin.BaseRoleServiceImpl,
	adminBaseUser *admin.BaseUserServiceImpl,

	adminDashboard *admin.DashboardServiceImpl,

	adminGoodsCategory *admin.GoodsCategoryServiceImpl,
	adminGoods *admin.GoodsServiceImpl,
	adminGoodsProp *admin.GoodsPropServiceImpl,
	adminGoodsSku *admin.GoodsSkuServiceImpl,
	adminGoodsSpec *admin.GoodsSpecServiceImpl,

	adminOrder *admin.OrderServiceImpl,

	adminPayBill *admin.PayBillServiceImpl,

	adminShopBanner *admin.ShopBannerServiceImpl,
	adminShopHot *admin.ShopHotServiceImpl,
	adminShopService *admin.ShopServiceServiceImpl,

	adminUserStore *admin.UserStoreServiceImpl,

	appAuth *app.AuthServiceImpl,
	appBaseArea *app.BaseAreaServiceImpl,
	appBaseDict *app.BaseDictServiceImpl,

	appGoodsCategory *app.GoodsCategoryServiceImpl,
	appGoods *app.GoodsServiceImpl,

	appOrder *app.OrderServiceImpl,

	appShopBanner *app.ShopBannerServiceImpl,
	appShopHot *app.ShopHotServiceImpl,
	appShopService *app.ShopServiceServiceImpl,

	appUserAddress *app.UserAddressServiceImpl,
	appUserCart *app.UserCartServiceImpl,
	appUserCollect *app.UserCollectServiceImpl,
	appUserStore *app.UserStoreServiceImpl,

	config *config.ConfigServiceImpl,
	file *file.FileServiceImpl,
	login *login.LoginServiceImpl,
	pay *pay.PayServiceImpl,
) *grpc.Server {
	if cfg == nil {
		panic("grpc server not configured")
	}
	srv := bootstrapServer.CreateGrpcServer(cfg,
		newGrpcMiddleware(cfg.GetMiddleware(), logger, authenticator, authorizer, userToken, jwt)...,
	)
	adminApi.RegisterAuthServiceServer(srv, adminAuth)
	adminApi.RegisterBaseApiServiceServer(srv, adminBaseApi)
	adminApi.RegisterBaseConfigServiceServer(srv, adminBaseConfig)
	adminApi.RegisterBaseDeptServiceServer(srv, adminBaseDept)
	adminApi.RegisterBaseDictServiceServer(srv, adminBaseDict)
	adminApi.RegisterBaseJobServiceServer(srv, adminBaseJob)
	adminApi.RegisterBaseLogServiceServer(srv, adminBaseLog)
	adminApi.RegisterBaseMenuServiceServer(srv, adminBaseMenu)
	adminApi.RegisterBaseRoleServiceServer(srv, adminBaseRole)
	adminApi.RegisterBaseUserServiceServer(srv, adminBaseUser)

	adminApi.RegisterDashboardServiceServer(srv, adminDashboard)

	adminApi.RegisterGoodsCategoryServiceServer(srv, adminGoodsCategory)
	adminApi.RegisterGoodsServiceServer(srv, adminGoods)
	adminApi.RegisterGoodsPropServiceServer(srv, adminGoodsProp)
	adminApi.RegisterGoodsSkuServiceServer(srv, adminGoodsSku)
	adminApi.RegisterGoodsSpecServiceServer(srv, adminGoodsSpec)

	adminApi.RegisterOrderServiceServer(srv, adminOrder)

	adminApi.RegisterPayBillServiceServer(srv, adminPayBill)

	adminApi.RegisterShopBannerServiceServer(srv, adminShopBanner)
	adminApi.RegisterShopHotServiceServer(srv, adminShopHot)
	adminApi.RegisterShopServiceServiceServer(srv, adminShopService)

	adminApi.RegisterUserStoreServiceServer(srv, adminUserStore)

	appApi.RegisterAuthServiceServer(srv, appAuth)
	appApi.RegisterBaseAreaServiceServer(srv, appBaseArea)
	appApi.RegisterBaseDictServiceServer(srv, appBaseDict)

	appApi.RegisterGoodsCategoryServiceServer(srv, appGoodsCategory)
	appApi.RegisterGoodsServiceServer(srv, appGoods)

	appApi.RegisterOrderServiceServer(srv, appOrder)

	appApi.RegisterShopBannerServiceServer(srv, appShopBanner)
	appApi.RegisterShopHotServiceServer(srv, appShopHot)
	appApi.RegisterShopServiceServiceServer(srv, appShopService)

	appApi.RegisterUserAddressServiceServer(srv, appUserAddress)
	appApi.RegisterUserCartServiceServer(srv, appUserCart)
	appApi.RegisterUserCollectServiceServer(srv, appUserCollect)
	appApi.RegisterUserStoreServiceServer(srv, appUserStore)

	configApi.RegisterConfigServiceServer(srv, config)
	fileApi.RegisterFileServiceServer(srv, file)
	loginApi.RegisterLoginServiceServer(srv, login)
	payApi.RegisterPayServiceServer(srv, pay)

	return srv
}
