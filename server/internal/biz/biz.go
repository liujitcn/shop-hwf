package biz

import (
	"gitee.com/liujit/shop/server/lib/data"
	adminBiz "gitee.com/liujit/shop/server/pkg/service/admin/biz"
	"gitee.com/liujit/shop/server/pkg/service/admin/task"
	appBiz "gitee.com/liujit/shop/server/pkg/service/app/biz"
	configBiz "gitee.com/liujit/shop/server/pkg/service/config/biz"
	fileBiz "gitee.com/liujit/shop/server/pkg/service/file/biz"
	payBiz "gitee.com/liujit/shop/server/pkg/service/pay/biz"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	data.NewData,
	data.NewTransaction,
	data.NewBaseApiRepo,
	data.NewBaseAreaRepo,
	data.NewBaseConfigRepo,
	data.NewBaseDeptRepo,
	data.NewBaseDictRepo,
	data.NewBaseDictItemRepo,
	data.NewBaseJobRepo,
	data.NewBaseJobLogRepo,
	data.NewBaseLogRepo,
	data.NewBaseMenuRepo,
	data.NewBaseRoleRepo,
	data.NewBaseUserRepo,

	data.NewCasbinRuleRepo,

	data.NewGoodsCategoryRepo,
	data.NewGoodsRepo,
	data.NewGoodsPropRepo,
	data.NewGoodsSpecRepo,
	data.NewGoodsSkuRepo,

	data.NewOrderRepo,
	data.NewOrderAddressRepo,
	data.NewOrderCancelRepo,
	data.NewOrderGoodsRepo,
	data.NewOrderLogisticsRepo,
	data.NewOrderPaymentRepo,
	data.NewOrderRefundRepo,

	data.NewPayBillRepo,

	data.NewShopBannerRepo,
	data.NewShopHotRepo,
	data.NewShopHotGoodsRepo,
	data.NewShopHotItemRepo,
	data.NewShopServiceRepo,

	data.NewUserAddressRepo,
	data.NewUserCartRepo,
	data.NewUserCollectRepo,
	data.NewUserStoreRepo,

	adminBiz.NewBaseApiCase,
	adminBiz.NewBaseAreaCase,
	adminBiz.NewBaseConfigCase,
	adminBiz.NewBaseDeptCase,
	adminBiz.NewBaseDictCase,
	adminBiz.NewBaseDictItemCase,
	adminBiz.NewBaseJobCase,
	adminBiz.NewBaseJobLogCase,
	adminBiz.NewBaseLogCase,
	adminBiz.NewBaseMenuCase,
	adminBiz.NewBaseRoleCase,
	adminBiz.NewBaseUserCase,

	adminBiz.NewCasbinRuleCase,

	adminBiz.NewDashboardCase,

	adminBiz.NewGoodsCategoryCase,
	adminBiz.NewGoodsCase,
	adminBiz.NewGoodsPropCase,
	adminBiz.NewGoodsSkuCase,
	adminBiz.NewGoodsSpecCase,

	adminBiz.NewOrderCase,
	adminBiz.NewOrderAddressCase,
	adminBiz.NewOrderCancelCase,
	adminBiz.NewOrderGoodsCase,
	adminBiz.NewOrderLogisticsCase,
	adminBiz.NewOrderPaymentCase,
	adminBiz.NewOrderRefundCase,

	adminBiz.NewPayBillCase,

	adminBiz.NewShopBannerCase,
	adminBiz.NewShopHotCase,
	adminBiz.NewShopHotItemCase,
	adminBiz.NewShopServiceCase,

	adminBiz.NewUserStoreCase,

	appBiz.NewBaseAreaCase,
	appBiz.NewBaseUserCase,
	appBiz.NewBaseRoleCase,
	appBiz.NewBaseDictCase,
	appBiz.NewBaseDictItemCase,

	appBiz.NewGoodsCategoryCase,
	appBiz.NewGoodsCase,
	appBiz.NewGoodsPropCase,
	appBiz.NewGoodsSkuCase,
	appBiz.NewGoodsSpecCase,

	appBiz.NewOrderCase,
	appBiz.NewOrderAddressCase,
	appBiz.NewOrderCancelCase,
	appBiz.NewOrderGoodsCase,
	appBiz.NewOrderLogisticsCase,
	appBiz.NewOrderPaymentCase,
	appBiz.NewOrderRefundCase,

	appBiz.NewShopBannerCase,
	appBiz.NewShopHotCase,
	appBiz.NewShopHotItemCase,
	appBiz.NewShopServiceCase,

	appBiz.NewUserAddressCase,
	appBiz.NewUserCartCase,
	appBiz.NewUserCollectCase,
	appBiz.NewUserStoreCase,

	configBiz.NewConfigCase,
	fileBiz.NewFileCase,

	payBiz.NewOrderSchedulerCase,
	payBiz.NewPayCase,
	payBiz.NewPayBillCase,
	payBiz.NewWxPayCase,

	task.NewTradeBill,
	task.NewTaskList,
)
