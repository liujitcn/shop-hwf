package data

import (
	"context"
	"github.com/google/wire"
	authnEngine "go.newcapec.cn/ncttools/nmskit-auth/authn/engine"
	"go.newcapec.cn/ncttools/nmskit-auth/authn/engine/jwt"
	authzEngine "go.newcapec.cn/ncttools/nmskit-auth/authz/engine"
	authzEngineCasbin "go.newcapec.cn/ncttools/nmskit-auth/authz/engine/casbin"
	nmskitAuthData "go.newcapec.cn/ncttools/nmskit-auth/data"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/cache"
	bootstrapCache "go.newcapec.cn/ncttools/nmskit-bootstrap/cache"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/conf"
	bootstrapConfigData "go.newcapec.cn/ncttools/nmskit-bootstrap/config/data"
	bootstrapOss "go.newcapec.cn/ncttools/nmskit-bootstrap/oss"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/pprof"
	bootstrapQueue "go.newcapec.cn/ncttools/nmskit-bootstrap/queue"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/sqldb"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(bootstrapConfigData.ParseServerHttp, bootstrapConfigData.ParseServerGrpc, bootstrapConfigData.ParseServerVue,
	bootstrapConfigData.ParseJwt,
	bootstrapConfigData.ParseDataRedis,
	bootstrapConfigData.ParseDataQueue,
	bootstrapCache.NewCache,
	bootstrapQueue.NewQueue,
	bootstrapConfigData.ParseDataDatabase,
	sqldb.NewGormSqlDb,
	bootstrapConfigData.ParseOss,
	bootstrapOss.NewOSS,
	bootstrapConfigData.ParsePprof,
	pprof.NewPprof,
	NewAuthenticator,
	NewAuthorizer,
	NewUserToken,
)

// NewAuthenticator 创建认证器
func NewAuthenticator(cfg *conf.Authentication_Jwt) authnEngine.Authenticator {
	authenticator, _ := jwt.NewAuthenticator(
		jwt.WithKey([]byte(cfg.GetSecret())),
		jwt.WithSigningMethod(cfg.GetMethod()),
	)
	return authenticator
}

// NewAuthorizer 创建权鉴器
func NewAuthorizer() (authzEngine.Engine, error) {
	return authzEngineCasbin.New(context.Background())
}

func NewUserToken(cfg *conf.Authentication_Jwt, cache cache.Cache, authenticator authnEngine.Authenticator) *nmskitAuthData.UserToken {
	const (
		userAccessTokenKeyPrefix  = "uat_"
		userRefreshTokenKeyPrefix = "urt_"
	)
	return nmskitAuthData.NewUserToken(cache, authenticator, userAccessTokenKeyPrefix, userRefreshTokenKeyPrefix, cfg.GetAccessTokenExpires().AsDuration(), cfg.GetRefreshTokenExpires().AsDuration())
}
