package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"relocate/middleware"
	"relocate/model"
	"relocate/service/staging_service"
	"relocate/util/app"
)

//TODO已完成: 设置当前分期(全局配置)(管理员)
//TODO已完成: 查询当前分期(全局配置)(管理员)
//TODO已完成: 设置总可选套数(全局配置)(管理员)
//TODO已完成: 查询总可选套数(全局配置)(管理员)

type QueryStagingIdBody struct {
	StagingId uint `json:"stagingId" form:"stagingId" validate:"required"`
}

// @Tags 全局配置
// @Summary 设置当前分期
// @Description 后台管理员设置当前分期
// @Produce  json
// @Security ApiKeyAuth
// @Param stagingId query int true "分期id值"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/config/staging/setting [get]
func SettingStagingConfig(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body QueryStagingIdBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	data, err := json.Marshal(body)
	if appG.HasError(err) {
		return
	}
	stagingID, err := model.GetNowStagingConfig()
	staging, _ := staging_service.GetStagingInfoById(stagingID)
	nowStagingName := staging.StagingName
	staging, err = staging_service.GetStagingInfoById(body.StagingId)
	if appG.HasError(err) {
		return
	}
	if appG.HasError(model.SettingNowStagingConfig(staging.ID)) {
		return
	}
	claim := middleware.GetClaims(c)
	logging := model.Logging{
		Username:    claim.Issuer,
		StagingName: nowStagingName,
		Operation:   "设置当前分期为" + staging.StagingName,
		Details:     string(data),
	}
	logging.Create()
	appG.SuccessResponse("设置成功")
}

// @Tags 全局配置
// @Summary 获取当前分期
// @Description 获取当前分期
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/config/staging/get [get]
func GetStagingConfig(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	stagingId, err := model.GetNowStagingConfig()
	if appG.HasError(err) {
		return
	}
	if stagingId == 0 {
		appG.SuccessResponse("配置未初始化")
		return
	}
	staging, err := staging_service.GetStagingInfoById(stagingId)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(*staging)
}

type QueryNowRoundsBody struct {
	Rounds uint `json:"rounds" form:"rounds" validate:"required"`
}

// @Tags 全局配置
// @Summary 设置当前轮次
// @Description 后台管理员设置当前轮次
// @Produce  json
// @Security ApiKeyAuth
// @Param rounds query int true "分期id值"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/config/rounds/setting [get]
func SettingNowRoundsConfig(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body QueryNowRoundsBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	if appG.HasError(model.SettingNowRoundsConfig(body.Rounds)) {
		return
	}
	appG.SuccessResponse("设置成功")
}

// @Tags 全局配置
// @Summary 获取当前轮次
// @Description 获取当前轮次
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/config/rounds/get [get]
func GetNowRoundsConfig(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	rounds, err := model.GetNowRoundsConfig()
	if appG.HasError(err) {
		return
	}
	if rounds == 0 {
		appG.SuccessResponse("配置未初始化")
		return
	}
	appG.SuccessResponse(map[string]interface{}{"rounds": rounds})
}

type QueryNumberBody struct {
	Num uint `json:"num" form:"num" validate:"min=1"`
}

type GroupingOptionalBody struct {
	Config []model.HuxingGroupingOptionalConfigJson `json:"config"`
}

// @Tags 全局配置
// @Summary 设置分组总可选套数
// @Description 后台管理员设置分组总可选套数
// @Produce  json
// @Security ApiKeyAuth
// @Param data body GroupingOptionalBody true "配置信息"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/config/huxing/groupingOptional/setting [post]
func SettingHuxingGroupingOptionalConfig(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body GroupingOptionalBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	if appG.HasError(model.SettingHuxingGroupingOptionalConfig(body.Config)) {
		return
	}
	appG.SuccessResponse("设置成功")
}

// @Tags 全局配置
// @Summary 查询分组总可选套数
// @Description 查询分组总可选套数
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/config/huxing/groupingOptional/get [get]
func GetHuxingGroupingOptionalConfig(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	config, err := model.GetHuxingGroupingOptionalConfig()
	if appG.HasError(err) {
		return
	}
	if config == nil {
		appG.SuccessResponse("配置未初始化")
		return
	}
	appG.SuccessResponse(config)
}

type QueryStateBody struct {
	State uint `json:"state" form:"state" validate:"oneof=0 1"`
}

// @Tags 全局配置
// @Summary 设置摇珠不中是否计入总数
// @Description 后台管理员设置摇珠不中是否计入总数
// @Produce  json
// @Security ApiKeyAuth
// @Param state query int true "状态 --0：申报不计入总数；1：申报计入总数"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/config/declaration/active/setting [get]
func SettingIncludeTotalConfig(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body QueryStateBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	if appG.HasError(model.SettingIncludeTotalConfig(body.State)) {
		return
	}
	appG.SuccessResponse("设置成功")
}

// @Tags 全局配置
// @Summary 查询摇珠不中是否计入总数状态(0：表示不计入总数，1：表示计入总数)
// @Description 查询摇珠不中是否计入总数状态(0：表示不计入总数，1：表示计入总数)
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/config/declaration/active/get [get]
func GetIncludeTotalConfig(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	state, err := model.GetIncludeTotalConfig()
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(map[string]uint{"IncludeTotalConfig": state})
}

type QueryHuxingTotalBody struct {
	Num uint `json:"num" form:"num" validate:"required"`
}

func SetAllOptionHuxing(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	var body QueryHuxingTotalBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	if appG.HasError(model.SetHuxingOptionalTotal(body.Num)) {
		return
	}
	appG.SuccessResponse("设置成功")
}

func GetAllOptionHuxing(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	config, err := model.GetHuxingOptionalTotal()
	if appG.HasError(err) {
		return
	}
	if *config == 0 {
		appG.SuccessResponse("配置未初始化")
		return
	}
	data := make(map[string]interface{})
	data["huxingOptionalConfig"] = *config
	appG.SuccessResponse(data)
}
