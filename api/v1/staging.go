package v1

//TODO已完成: 新增分期(管理员)
//TODO已完成: 分页获取分期列表
//TODO已完成: 根据分期期数导入合同单列表(管理员)

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"relocate/api"
	"relocate/middleware"
	"relocate/model"
	"relocate/service/staging_service"
	"relocate/util/app"
)

type NewStagingBody struct {
	StagingName string `json:"stagingName" validate:"required"`
	NowStaging  bool   `json:"nowStaging"`
}

// @Tags 分期
// @Summary 导入合同单原始数据
// @Description 后台管理员根据分期期数导入合同单原始数据
// @Accept 	multipart/form-data
// @Produce  json
// @Security ApiKeyAuth
// @Param stagingId formData int true "分期期数"
// @Param excel formData file true "拆迁人Excel原始数据"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/staging/import [post]
func ImportContractByStaging(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body QueryStagingIdBody
	if !appG.ParseFormRequest(&body) {
		return
	}
	fileHeader, err := api.Upload(c, "excel", 10, []string{".xlsx", ".xls"})
	if appG.HasError(err) {
		return
	}
	file, err := fileHeader.Open()
	defer file.Close()
	if appG.HasError(err) {
		return
	}
	claim := middleware.GetClaims(c)
	count, err := staging_service.ImportContract(body.StagingId, file, claim.Issuer)
	if appG.HasError(err) {
		return
	}

	appG.SuccessResponse(fmt.Sprintf("导入成功,条数为:%d", count))
}

// @Tags 分期
// @Summary 分页查询分期信息(包括每个分期合同单的数量)
// @Description 分页根据分期id查询分期信息(包括每个分期合同单的数量)
// @Produce  json
// @Security ApiKeyAuth
// @Param page query int false "页码"
// @Param pageSize query int false "页面大小"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/staging/getCount [get]
func GetStagingContractCount(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body api.PaginationQueryBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	data, err := model.GetStagingContractCount(body.PageSize, body.Page)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.BadResponse("查询分期失败")
			return
		}
		if appG.HasError(err) {
			return
		}
	}
	stagingId, err := model.GetNowStagingConfig()
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(map[string]interface{}{"stagingList": data, "nowStagingID": stagingId})

}

func GetAllStagingPage(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	var body api.PaginationQueryBody
	if !appG.ParseQueryRequest(&body) {
		return
	}

	page, err := model.GetStagingPage(body.Page, body.PageSize)

	if err != nil {
		appG.BadResponse(err.Error())
		return
	}
	appG.SuccessResponse(page)
}

// AddStaging 新增分期
func AddStaging(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	var body NewStagingBody

	if !appG.ParseJSONRequest(&body) {
		return
	}

	if err := staging_service.AddStaging(body.StagingName, body.NowStaging); err != nil {
		appG.BadResponse(err.Error())
	}

	appG.SuccessResponse("添加成功")

}
