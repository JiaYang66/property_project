package v1

//TODO已完成: 获取户型数据列表
//TODO已完成: 根据ID删除户型数据(管理员)
//TODO已完成: 根据ID修改户型数据(管理员)
//TODO已完成: 新增户型数据(管理员)
import (
	"github.com/gin-gonic/gin"
	"relocate/model"
	"relocate/service/huxing_service"
	"relocate/util/app"
)

type GetHuxingListBody struct {
	StagingID uint `json:"staging_id" form:"staging_id" validate:"required"`
}

// @Tags 户型
// @Summary 获取户型（可选>0以上的）
// @Description 获取户型（可选>0以上的）
// @Produce  json
// @Param staging_id query int true "分期期数"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/huxing/get/optional [get]
func GetOptionalHuxing(ctx *gin.Context) {
	appG := app.Gin{Ctx: ctx}
	var body GetHuxingListBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	huxingList, err := model.FindAllOptionalHuxing(body.StagingID)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(huxingList)
}

func GetAllHuxing(ctx *gin.Context) {
	appG := app.Gin{Ctx: ctx}
	var body GetHuxingListBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	huxingList, err := model.FindAllHuxing(body.StagingID)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(huxingList)
}

type DeleteHuxingBody struct {
	Id uint `json:"id" form:"id" validate:"required"`
}

func DeleteHuxing(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	var body DeleteHuxingBody
	if !appG.ParseJSONRequest(&body) {
		return
	}

	err := model.DelteHuxingById(body.Id)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse("删除成功")

}

type AddHuxingBody struct {
	StagingID  uint   `json:"staging_id" `
	BuildingNo string `json:"building_no" validate:"required"`
	HuxingNo   string `json:"huxing_no" validate:"required"`
	Area       string `json:"area" validate:"required"`
	AreaShow   string `json:"area_show" validate:"required"`
	Quantity   uint   `json:"quantity" validate:"required"`
	Maximum    uint   `json:"maximum" `
	Rounds     uint   `json:"rounds" `
}

func AddHuxing(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	var body AddHuxingBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	err := huxing_service.AddHuxing(body.StagingID,
		body.BuildingNo,
		body.HuxingNo,
		body.Area,
		body.AreaShow,
		body.Quantity,
		body.Maximum,
		body.Rounds)

	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse("添加成功")
}

type UpdateHuxingBody struct {
	Id         uint   `json:"id" validate:"required"`
	BuildingNo string `json:"building_no" `
	HuxingNo   string `json:"huxing_no" `
	Area       string `json:"area" validate:"required"`
	AreaShow   string `json:"area_show" validate:"required"`
	Quantity   uint   `json:"quantity" `
	Maximum    uint   `json:"maximum" `
	Rounds     uint   `json:"rounds" `
}

func UpdateHuxing(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	var body UpdateHuxingBody
	if !appG.ParseJSONRequest(&body) {
		return
	}

	err := huxing_service.UpdateHuxing(
		body.Id,
		body.BuildingNo,
		body.HuxingNo,
		body.Area,
		body.AreaShow,
		body.Quantity,
		body.Maximum,
		body.Rounds)

	if appG.HasError(err) {
		return
	}

	appG.SuccessResponse("修改成功")

}
