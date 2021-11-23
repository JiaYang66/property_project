package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"relocate/util/app"
	"relocate/util/sign"
	"strings"
)

func JWT(claimsType ...sign.ClaimsType) gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{Ctx: c}
		auth := c.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			appG.UnauthorizedResponse("请登录后再进行操作")
			c.Abort() // 阻断执行下一个HandlerFunc，仅会执行接下去的代码
			return
		}
		if len(strings.Fields(auth)) > 1 {
			auth = strings.Fields(auth)[1]
		}
		// 校验token
		claims, err := sign.ParseToken(auth)
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				appG.UnauthorizedResponse("token 已过期")
			default:
				appG.UnauthorizedResponse("token 验证失败")
			}
			c.Abort()
			return
		}
		if len(claimsType) > 0 {
			switch claimsType[0] {
			case sign.AdminClaimsType:
				if claims.Type != sign.AdminClaimsType {
					appG.UnauthorizedResponse("请使用管理员账号登录")
					c.Abort()
					return
				}
				break
			case sign.UserClaimsType:
				if claims.Type != sign.UserClaimsType {
					appG.UnauthorizedResponse("请使用用户账号登录")
					c.Abort()
					return
				}
				break
			default:
				appG.UnauthorizedResponse("unknown")
				c.Abort()
				return
			}
		}
		// 跨中间件设置值
		c.Set("claims", claims)
		// 立即执行下一个HandlerFunc
		c.Next()
	}
}

func GetClaims(c *gin.Context) *sign.Claims {
	// 跨中间件取值
	claims, ok := c.Get("claims")
	if ok {
		// 通过断言获取claims内的所有内容
		return claims.(*sign.Claims)
	}
	return nil
}
