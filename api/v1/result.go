package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"relocate/api"
	"relocate/model"
	"relocate/service/result_service"
	"relocate/util/app"
	"strconv"
)

//TODO已完成: 根据申报表ID录入结果(录入结果人员 管理员姓名)(管理员)
//TODO已完成: 根据分期期数（必须）、户型、身份证号、合同号、手机号 模糊查询分页查询摇珠结果列表（筛选功能）(管理员)
//TODO已完成: 根据合同号(批量)设置公示状态(管理员)
//TODO已完成: 根据公示状态分页获取摇珠结果列表

// @Tags 公示
// @Summary 根据分期期数（必须）、户型、公示状态、身份证号、合同号、手机号 模糊查询分页查询摇珠结果列表（筛选功能）(管理员)
// @Description 根据分期期数（必须）、户型、身份证号、合同号、手机号 模糊查询分页查询摇珠结果列表（筛选功能）(管理员)
// @Produce  json
// @Security ApiKeyAuth
// @Param stagingId query int true "分期期数（必须）"
// @Param filterName query string false "搜索关键字(姓名、身份证号、合同号、手机号)"
// @Param huxing_id query int false "户型id"
// @Param 公示状态 query int false "是否公示  0表示公示状态为false的数据，1表示公示状态为true的数据，不选则显示所有"
// @Param page query int false "页码"
// @Param pageSize query int false "页面大小"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/result/get [get]
func GetResultList(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body api.ResultFilterBody
	body.PublicityStatus = -1
	if !appG.ParseQueryRequest(&body) {
		return
	}
	data, err := model.GetLikeResults(body)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(data)
}

type UpdateResultStatusBody struct {
	DeclarationID   []int `json:"declaration_id" form:"declaration_id"`     //要批量修改的申报号数组
	PublicityStatus bool  `json:"publicity_status" form:"publicity_status"` //要修改的公示状态(是否公示)
}

type ResultQueryBody struct {
	Page     uint `json:"page" form:"page"`
	PageSize uint `json:"pageSize" form:"pageSize"`
}

// @Tags 公示
// @Summary 导出摇珠表列表(管理员)
// @Description 根据分期期数（必须）导出摇珠表列表(管理员)
// @Produce  json
// @Security ApiKeyAuth
// @Param stagingId query int true "分期期数（必须）"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/result/export [post]
func ExportResults(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body api.ResultFilterBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	body.PageSize = 9999
	body.PublicityStatus = -1
	data, err := model.GetLikeResults(body)
	if appG.HasError(err) {
		return
	}
	resultDataList, ok := data.Data.(*[]model.ResultData)
	if !ok {
		appG.BadResponse("导出过程发生异常")
		return
	}
	file, filename, err := result_service.ExportExcel(*resultDataList)
	if appG.HasError(err) {
		return
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url.QueryEscape(filename)))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	buffer, err := file.WriteToBuffer()
	if appG.HasError(err) {
		return
	}
	c.Writer.Header().Add("Content-Length", strconv.FormatInt(int64(buffer.Len()), 10))
	if appG.HasError(file.Write(c.Writer)) {
		return
	}
}

func GetResultByStatus(context *gin.Context) {
	// 根据公示状态分页获取摇珠结果列表(移动端)
	// status 为 1才显示
	appG := app.Gin{Ctx: context}
	var body api.PaginationQueryBody
	if !appG.ParseQueryRequest(&body) {
		return
	}

	data, err := model.GetResultByStatus(body.Page, body.PageSize)
	if err != nil {
		appG.BadResponse(err.Error())
		return
	}

	appG.SuccessResponse(data)
}

type ResultPublicStatBody struct {
	DeclarationId   []int `json:"declaration_id"`
	PublicityStatus bool  `json:"publicity_status"`
}

func UpdateResultPublicStat(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	var body ResultPublicStatBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	fmt.Println(body.PublicityStatus)
	if err := result_service.UpdateResultPublicStat(body.DeclarationId, body.PublicityStatus); err != nil {
		appG.BadResponse(err.Error())
		return
	}

	appG.SuccessResponse("修改成功")
}
