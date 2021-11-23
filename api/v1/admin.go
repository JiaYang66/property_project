package v1

import (
	"github.com/gin-gonic/gin"
	"relocate/middleware"
	"relocate/model"
	"relocate/service/admin_service"
	"relocate/util/app"
)

//前端将管理员输入的密码MD5加密后传输到后端
//后端将MD5密文拼接上数据库保存的盐值继续MD5后与数据库保存的密码进行匹配
//TODO手动数据库插入: 管理员注册(不开放注册，上线关闭注册功能)
//TODO已完成: 管理员登录

type AdminLoginBody struct {
	AdminName string `json:"adminName" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

// @Tags 管理员
// @Summary 管理员登录账号admin密码123456
// @Produce json
// @Param data body AdminLoginBody true "登录信息"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/admin/login [post]
func AdminLogin(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body AdminLoginBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	token, err := admin_service.GenerateToken(body.AdminName, body.Password)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(token)
}

type AdminUpdateBody struct {
	AdminName   string `json:"adminName" validate:"required"`
	Password    string `json:"password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

// @Tags 管理员
// @Summary 修改密码(账号，旧密码，新密码)
// @Produce json
// @Security ApiKeyAuth
// @Param data body AdminUpdateBody true "手机号码，旧密码，新密码"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/admin/update [post]
func AdminUpdatePassword(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body AdminUpdateBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	if appG.HasError(admin_service.UpdateAdminPassword(body.AdminName, body.Password, body.NewPassword)) {
		return
	}
	appG.SuccessResponse("修改密码成功，请重新登录")
}

// @Tags 管理员
// @Summary 获取所有的管理员账号和姓名
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/admin/get [get]
func GetAllAdmin(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	adminList, err := model.GetAllAdminInfo()
	if appG.HasError(err) {
		return
	}
	claim := middleware.GetClaims(c)
	user, err := model.GetAdminInfo(claim.Issuer)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(map[string]interface{}{"adminList": adminList, "currentAdmin": user})
}
