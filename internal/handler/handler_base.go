package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qwqserver/internal/auth"
	"qwqserver/internal/common"
)

type HandleBase interface {
	AuthMiddleware(c *gin.Context)
}

type HandleBaseImpl struct {
	AuthCode  auth.CodeType
	AuthAbort bool
}

// AuthMiddleware 认证中间件 true abort false continue
func (h *HandleBaseImpl) AuthMiddleware(c *gin.Context, authAbort ...bool) (bool, *common.HTTPResult) {
	res := &common.HTTPResult{}
	res.Msg = "认证成功"
	isValid := true
	if len(authAbort) >= 1 {
		isValid = authAbort[0]
	}

	//

	authCode, ok := c.Get(auth.IdentityStatusKey)
	if !ok {
		res.Code = http.StatusUnauthorized
		res.Msg = "认证中间件错误"
	} else if authCode.(auth.CodeType) != auth.IdentityOK {
		res.Code = http.StatusUnauthorized
		res.Msg = "认证失败"
	}
	if authCode == auth.IdentityErrNoToken {
		//msg = "未提供认证令牌"
		res.Code = http.StatusUnauthorized
		res.Msg = "未提供认证令牌"
	} else if authCode == auth.IdentityErrTokenFormat {
		//msg = "令牌格式错误"
		res.Code = http.StatusUnauthorized
		res.Msg = "令牌格式错误"
	} else if authCode == auth.IdentityErrInvalidToken {
		//msg = "无效令牌"
		res.Code = http.StatusUnauthorized
		res.Msg = "无效令牌"
	} else if authCode == auth.IdentityErrTokenExpired {
		//msg = "令牌已过期"
		res.Code = http.StatusUnauthorized
		res.Msg = "令牌已过期"
	}
	// 此条件用于特殊需求，比如注册、登录...如果允许放行，并且未存在令牌，则放行
	if !isValid {
		c.Next()
		return true, res
	}
	if authCode == auth.IdentityOK {
		//msg = "请不要重复登录"
		res.Code = http.StatusUnauthorized
		res.Msg = "请不要重复登录"
		//c.AbortWithStatusJSON(401, res)
		c.Abort()
		return false, res
	}
	//c.AbortWithStatusJSON(401, res)
	c.Abort()
	return false, res
}
