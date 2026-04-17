package server

import (
	adminApi "gitee.com/liujit/shop/server/api/admin"
	appApi "gitee.com/liujit/shop/server/api/app"
	configApi "gitee.com/liujit/shop/server/api/config"
	loginApi "gitee.com/liujit/shop/server/api/login"
	payApi "gitee.com/liujit/shop/server/api/pay"
	"gitee.com/liujit/shop/server/cmd/assets"
	"gitee.com/liujit/shop/server/internal/service"
	"gitee.com/liujit/shop/server/pkg/service/admin"
	adminBiz "gitee.com/liujit/shop/server/pkg/service/admin/biz"
	"gitee.com/liujit/shop/server/pkg/service/app"
	"gitee.com/liujit/shop/server/pkg/service/config"
	"gitee.com/liujit/shop/server/pkg/service/file"
	"gitee.com/liujit/shop/server/pkg/service/login"
	"gitee.com/liujit/shop/server/pkg/service/pay"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	authnEngine "go.newcapec.cn/ncttools/nmskit-auth/authn/engine"
	authzEngine "go.newcapec.cn/ncttools/nmskit-auth/authz/engine"
	nmskitAuthData "go.newcapec.cn/ncttools/nmskit-auth/data"
	bootstrapConf "go.newcapec.cn/ncttools/nmskit-bootstrap/conf"
	bootstrapServer "go.newcapec.cn/ncttools/nmskit-bootstrap/server"
	swaggerUI "go.newcapec.cn/ncttools/nmskit-swagger-ui"
	"go.newcapec.cn/ncttools/nmskit/log"
	"go.newcapec.cn/ncttools/nmskit/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	cfg *bootstrapConf.Server_HTTP,
	jwt *bootstrapConf.Authentication_Jwt,
	authenticator authnEngine.Authenticator,
	authorizer authzEngine.Engine,
	userToken *nmskitAuthData.UserToken,

	adminBaseUserCase *adminBiz.BaseUserCase,

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

	adminPayBill *admin.PayBillServiceImpl,

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
) *http.Server {
	if cfg == nil {
		log.Errorf("HTTP server adminBaseConfig is nil")
		return nil
	}
	srv := bootstrapServer.CreateHttpServer(cfg,
		newHttpMiddleware(cfg.GetMiddleware(), logger, adminBaseUserCase, authenticator, authorizer, userToken, jwt)...,
	)

	adminApi.RegisterAuthServiceHTTPServer(srv, adminAuth)
	adminApi.RegisterBaseApiServiceHTTPServer(srv, adminBaseApi)
	adminApi.RegisterBaseConfigServiceHTTPServer(srv, adminBaseConfig)
	adminApi.RegisterBaseDeptServiceHTTPServer(srv, adminBaseDept)
	adminApi.RegisterBaseDictServiceHTTPServer(srv, adminBaseDict)
	adminApi.RegisterBaseJobServiceHTTPServer(srv, adminBaseJob)
	adminApi.RegisterBaseLogServiceHTTPServer(srv, adminBaseLog)
	adminApi.RegisterBaseMenuServiceHTTPServer(srv, adminBaseMenu)
	adminApi.RegisterBaseRoleServiceHTTPServer(srv, adminBaseRole)
	adminApi.RegisterBaseUserServiceHTTPServer(srv, adminBaseUser)

	adminApi.RegisterDashboardServiceHTTPServer(srv, adminDashboard)

	adminApi.RegisterGoodsCategoryServiceHTTPServer(srv, adminGoodsCategory)
	adminApi.RegisterGoodsServiceHTTPServer(srv, adminGoods)
	adminApi.RegisterGoodsPropServiceHTTPServer(srv, adminGoodsProp)
	adminApi.RegisterGoodsSkuServiceHTTPServer(srv, adminGoodsSku)
	adminApi.RegisterGoodsSpecServiceHTTPServer(srv, adminGoodsSpec)

	adminApi.RegisterOrderServiceHTTPServer(srv, adminOrder)

	adminApi.RegisterPayBillServiceHTTPServer(srv, adminPayBill)

	adminApi.RegisterShopBannerServiceHTTPServer(srv, adminShopBanner)
	adminApi.RegisterShopHotServiceHTTPServer(srv, adminShopHot)
	adminApi.RegisterShopServiceServiceHTTPServer(srv, adminShopService)

	adminApi.RegisterUserStoreServiceHTTPServer(srv, adminUserStore)

	appApi.RegisterAuthServiceHTTPServer(srv, appAuth)
	appApi.RegisterBaseAreaServiceHTTPServer(srv, appBaseArea)
	appApi.RegisterBaseDictServiceHTTPServer(srv, appBaseDict)

	appApi.RegisterGoodsCategoryServiceHTTPServer(srv, appGoodsCategory)
	appApi.RegisterGoodsServiceHTTPServer(srv, appGoods)

	appApi.RegisterOrderServiceHTTPServer(srv, appOrder)

	appApi.RegisterShopBannerServiceHTTPServer(srv, appShopBanner)
	appApi.RegisterShopHotServiceHTTPServer(srv, appShopHot)
	appApi.RegisterShopServiceServiceHTTPServer(srv, appShopService)

	appApi.RegisterUserAddressServiceHTTPServer(srv, appUserAddress)
	appApi.RegisterUserCartServiceHTTPServer(srv, appUserCart)
	appApi.RegisterUserCollectServiceHTTPServer(srv, appUserCollect)
	appApi.RegisterUserStoreServiceHTTPServer(srv, appUserStore)

	configApi.RegisterConfigServiceHTTPServer(srv, config)
	// 修改http接口实现
	service.RegisterFileServiceHTTPServer(srv, file)
	loginApi.RegisterLoginServiceHTTPServer(srv, login)
	payApi.RegisterPayServiceHTTPServer(srv, pay)

	srv.Handle("/metrics", promhttp.Handler())
	swaggerUI.RegisterSwaggerUIServerWithOption(
		srv,
		swaggerUI.WithTitle("NMS Service"),
		swaggerUI.WithMemoryData(assets.OpenApiData, "yaml"),
	)
	return srv
}
