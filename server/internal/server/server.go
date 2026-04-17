package server

import (
	customLogging "gitee.com/liujit/shop/server/internal/middleware/logging"
	"gitee.com/liujit/shop/server/pkg/service/admin/biz"
	"github.com/google/wire"
	authnEngine "go.newcapec.cn/ncttools/nmskit-auth/authn/engine"
	authzEngine "go.newcapec.cn/ncttools/nmskit-auth/authz/engine"
	nmskitAuthData "go.newcapec.cn/ncttools/nmskit-auth/data"
	authMiddleware "go.newcapec.cn/ncttools/nmskit-auth/middleware"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/conf"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/vue"
	"go.newcapec.cn/ncttools/nmskit/log"
	"go.newcapec.cn/ncttools/nmskit/middleware"
	"go.newcapec.cn/ncttools/nmskit/middleware/logging"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer, vue.NewVueServer)

// NewMiddleware 创建中间件
func newHttpMiddleware(
	cfg *conf.Server_Middleware,
	logger log.Logger,
	userCase *biz.BaseUserCase,
	authenticator authnEngine.Authenticator,
	authorizer authzEngine.Engine,
	userToken *nmskitAuthData.UserToken,
	jwt *conf.Authentication_Jwt,
) []middleware.Middleware {
	var ms []middleware.Middleware
	var enableLogging bool
	if cfg != nil {
		enableLogging = cfg.GetEnableLogging()
	}
	if enableLogging {
		ms = append(ms, customLogging.Server(logger, userCase, authenticator))
	}
	ms = append(ms, authMiddleware.NewAuthMiddleware(authenticator, authorizer, userToken, jwt))

	return ms
}

// NewMiddleware 创建中间件
func newGrpcMiddleware(
	cfg *conf.Server_Middleware,
	logger log.Logger,
	authenticator authnEngine.Authenticator,
	authorizer authzEngine.Engine,
	userToken *nmskitAuthData.UserToken,
	jwt *conf.Authentication_Jwt,
) []middleware.Middleware {
	var ms []middleware.Middleware
	var enableLogging bool
	if cfg != nil {
		enableLogging = cfg.GetEnableLogging()
	}
	if enableLogging {
		ms = append(ms, logging.Server(logger))
	}
	ms = append(ms, authMiddleware.NewAuthMiddleware(authenticator, authorizer, userToken, jwt))
	return ms
}
