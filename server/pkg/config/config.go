package config

import (
	"errors"
	_const "gitee.com/liujit/shop/server/lib/const"
	"go.newcapec.cn/nctcommon/nmslib"
	nmslibConfig "go.newcapec.cn/nctcommon/nmslib/config"
	"go.newcapec.cn/ncttools/nmskit/log"
	"strconv"
	"time"
)

func ParseShopConfig() *ShopConfig {
	var shopConfig *ShopConfig
	var shopSection ShopSection
	err := nmslibConfig.Scan(&shopSection)
	if err != nil {
		log.Fatalf("nmslibConfig.Read err: %v", err)
	} else {
		shopConfig = shopSection.GetShopConfig()
	}
	return shopConfig
}

func ParseWxMiniApp(cfg *ShopConfig) (*WxMiniApp, error) {
	wxMiniApp := cfg.GetWxMiniApp()
	if wxMiniApp == nil {
		return nil, errors.New("微信登录配置信息错误")
	}
	appid := wxMiniApp.GetAppid()
	secret := wxMiniApp.GetSecret()
	if appid == "" || secret == "" {
		return nil, errors.New("微信登录配置信息错误")
	}
	return wxMiniApp, nil
}

func ParseWxPay(cfg *ShopConfig) (*WxPay, error) {
	wxPay := cfg.GetWxPay()
	if wxPay == nil {
		return nil, errors.New("支付配置信息错误")
	}
	appid := wxPay.GetAppid()
	mchId := wxPay.GetMchId()
	mchCertSn := wxPay.GetMchCertSn()
	mchCertPath := wxPay.GetMchCertPath()
	mchAPIv3Key := wxPay.GetMchAPIv3Key()
	if appid == "" || mchId == "" || mchCertSn == "" || mchCertPath == "" || mchAPIv3Key == "" {
		return nil, errors.New("支付配置信息错误")
	}
	return wxPay, nil
}

func ParsePayTimeout() time.Duration {
	cache := nmslib.Runtime.GetCache()
	if cache == nil {
		// 默认30分钟
		return time.Duration(_const.PayTimeout) * time.Minute
	}

	v, err := cache.Get(_const.CacheKeyConfig + _const.CacheKeyPayTimeout)
	if err != nil {
		// 默认30分钟
		return time.Duration(_const.PayTimeout) * time.Minute
	}
	var payTimeout int
	payTimeout, err = strconv.Atoi(v)
	if err != nil {
		// 默认30分钟
		return time.Duration(_const.PayTimeout) * time.Minute
	}
	// 默认30分钟
	return time.Duration(payTimeout) * time.Minute
}
