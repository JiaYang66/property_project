package middleware

import (
	"github.com/gin-gonic/gin"
	"relocate/model"
	"relocate/util/app"
)

func Visitors(isAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{Ctx: c}
		claim := GetClaims(c)
		if claim != nil {
			admin, err := model.GetAdminInfo(claim.Issuer)
			if err != nil {
				c.Abort()
				return
			}
			if isAdmin {
				//只有超级管理员可以操作的
				if admin.Visitors {
					//当前管理员是访客
					appG.VisitorsResponse("您当前处于访客模式，无权限访问")
					c.Abort()
					return
				}
			}
		}
		// 立即执行下一个HandlerFunc
		c.Next()
	}
}
