package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"relocate/api"
	"relocate/model"
	"relocate/service/accounting_service"
	"relocate/util/app"
	"strconv"
)

type AddAccountingBody struct {
	ContractNoList []string `json:"contract_no_list" validate:"required"`
}

// @Tags 核算
// @Summary 根据合同号核算
// @Produce  json
// @Security ApiKeyAuth
// @Param data body AddAccountingBody true "核算信息"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/accounting [post]
func AddAccounting(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body AddAccountingBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	if appG.HasError(accounting_service.AddAccounting(body.ContractNoList)) {
		return
	}
	appG.SuccessResponse("核算成功")
}

// @Tags 核算
// @Summary 根据合同号、被拆迁人模糊查询分页查询核算列表（筛选功能）(管理员)
// @Produce  json
// @Security ApiKeyAuth
// @Param filterName query string false "搜索关键字(合同号、被拆迁人)"
// @Param page query int false "页码"
// @Param pageSize query int false "页面大小"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/accounting [get]
func GetAccountingList(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body api.AccountingFilterBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	data, err := model.GetLikeAccounting(body)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(data)
}

// @Tags 核算
// @Summary 导出核算表(管理员)
// @Produce  json
// @Security ApiKeyAuth
// @Param filterName query string false "搜索关键字(合同号、被拆迁人)"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/accounting/export [get]
func ExportAccounting(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body api.AccountingFilterBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	body.Page = 1
	body.PageSize = 9999
	data, err := model.GetLikeAccounting(body)
	if appG.HasError(err) {
		return
	}
	accountingDataList, ok := data.Data.(*[]model.Accounting)
	if !ok {
		appG.BadResponse("导出过程发生异常")
		return
	}
	file, filename, err := accounting_service.ExportExcel(*accountingDataList)
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
