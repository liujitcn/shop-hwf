package util

import (
	"context"
	_const "gitee.com/liujit/shop/server/lib/const"
	"go.newcapec.cn/ncttools/nmskit-auth/data"
	authMiddleware "go.newcapec.cn/ncttools/nmskit-auth/middleware"
)

func IsMember(ctx context.Context) bool {
	authInfo, err := authMiddleware.FromContext(ctx)
	if err != nil {
		return false
	}
	if len(authInfo.RoleCode) == 0 || authInfo.RoleCode == _const.BaseRoleCode_Guest {
		return false
	}
	if authInfo.RoleCode == _const.BaseRoleCode_User {
		return true
	}
	return false
}

func IsMemberByAuthInfo(authInfo *data.UserTokenPayload) bool {
	if len(authInfo.RoleCode) == 0 || authInfo.RoleCode == _const.BaseRoleCode_Guest {
		return false
	}
	if authInfo.RoleCode == _const.BaseRoleCode_User {
		return true
	}
	return false
}
