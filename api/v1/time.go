package v1

import (
	"github.com/gin-gonic/gin"
	"relocate/service/time_service"
	"relocate/util/app"
)

//TODO已完成: 获取现场确认时间段数据列表
//TODO已完成: 根据ID删除现场确认时间段数据(管理员)
//TODO已完成: 根据ID修改现场确认时间段数据(管理员)
//TODO已完成: 新增现场确认时间段数据(管理员)

type NewTimeBody struct {
	Name        string `json:"timeName" form:"timeName" validate:"required"`       //时间段
	OptionalNum uint   `json:"optionalNum" form:"optionalNum" validate:"required"` //可选数
	StagingId   uint   `json:"staging_id" form:"staging_id" validate:"required"`   //分期id
	//SelectedNum uint   `json:"selectedNum" form:"selectedNum" validate:"ltefield=OptionalNum"` //已选数
}

type GetTimeByStagingId struct {
	StagingId uint `json:"staging_id" form:"staging_id" validate:"required"`
}

type UpdateTimeBody struct {
	Id          uint   `json:"id" form:"id" validate:"required"`
	Name        string `json:"name" form:"name" `                //时间段
	OptionalNum uint   `json:"optional_num" form:"optional_num"` //可选数
	//	SelectedNum uint   `json:"selected_num" form:"selected_num" ` //已选数
}

type DeleteTimeBody struct {
	Id uint `json:"id" form:"id" validate:"required"`
}

func QueryTimeByStagingId(context *gin.Context) {
	appG := app.Gin{
		context,
	}
	var body GetTimeByStagingId
	if !appG.ParseQueryRequest(&body) {
		return
	}

	times, err := time_service.GetTimeByStagingId(body.StagingId)

	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(times)

}

func AddTime(context *gin.Context) {
	appG := app.Gin{
		Ctx: context,
	}
	var body NewTimeBody
	if !appG.ParseJSONRequest(&body) {
		return
	}

	err := time_service.SaveUser(body.Name, body.StagingId, body.OptionalNum)

	if appG.HasError(err) {
		appG.BadResponse(err)
		return
	}

	appG.SuccessResponse("添加成功")

}

func UpdateTime(context *gin.Context) {
	appG := app.Gin{
		Ctx: context,
	}
	var body UpdateTimeBody
	body.Id = 0
	body.OptionalNum = 0
	if !appG.ParseJSONRequest(&body) {
		return
	}

	err := time_service.UpdateTime(body.Id, body.Name, body.OptionalNum)

	if !appG.HasError(err) {
		appG.SuccessResponse("修改成功")
		return
	}

	appG.BadResponse("修改失败")

}

func DeleteTimeById(context *gin.Context) {
	appG := app.Gin{
		Ctx: context,
	}
	var body DeleteTimeBody
	body.Id = 0
	if !appG.ParseJSONRequest(&body) {
		return
	}

	err := time_service.DeleteById(body.Id)

	if appG.HasError(err) {
		appG.BadResponse(err)
		return
	}

	appG.SuccessResponse("删除成功")
}
