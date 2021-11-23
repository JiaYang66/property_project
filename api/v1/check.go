package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"relocate/api"
	"relocate/model"
	"relocate/service/check_service"
	"relocate/util/app"
	"strconv"
)

type AddCheckBody struct {
	ResultID    uint    `json:"result_id" validate:"required"`
	RealityArea float64 `json:"reality_area"`
}

type FuzzyQueryByContractNoOrPeoplesBody struct {
	ContractNo string `json:"contract_no" form:"contract_no"`
	Peoples    string `json:"peoples" form:"peoples"`
	Page       uint   `json:"page" form:"page"`
	PageSize   uint   `json:"pageSize" form:"pageSize"`
}

func CheckFuzzyQueryByContractNoOrPeoples(context *gin.Context) {
	appG := app.Gin{
		Ctx: context,
	}
	var body FuzzyQueryByContractNoOrPeoplesBody
	if !appG.ParseQueryRequest(&body) {
		return
	}

	data, err := model.QueryCheckByContractNoOrPeoples(body.ContractNo, body.Peoples, body.Page, body.PageSize)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(data)
}

type AccountingByResultsBody struct {
	ResultId uint `json:"result_id" validate:"required"`
}

func AccountingByResults(context *gin.Context) {
	appG := app.Gin{
		Ctx: context,
	}
	var body AccountingByResultsBody
	if !appG.ParseJSONRequest(&body) {
		return
	}

	err := check_service.AccountingByResults(body.ResultId)
	if appG.HasError(err) {
		return
	}

	appG.SuccessResponse("修改成功")
}

func ExportCheck(context *gin.Context) {
	appG := app.Gin{
		Ctx: context,
	}
	var body api.CheckFilterBody
	if !appG.ParseJSONRequest(&body) {
		return
	}

	file, filename, err := check_service.ExportCheck(body)
	if appG.HasError(err) {
		return
	}
	context.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url.QueryEscape(filename)))
	context.Writer.Header().Add("Content-Type", "application/octet-stream")
	buffer, err := file.WriteToBuffer()
	if appG.HasError(err) {
		return
	}
	context.Writer.Header().Add("Content-Length", strconv.FormatInt(int64(buffer.Len()), 10))
	if appG.HasError(file.Write(context.Writer)) {
		return
	}

}
