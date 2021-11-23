package v1

import (
	"github.com/gin-gonic/gin"
	"relocate/model"
	"relocate/util/app"
)

//TODO已完成: 分页获取合同号面积明细记录

type GetAreaDetailPageByContextNoBody struct {
	ContractNo string `json:"contract_no" form:"contract_on" validate:"required"`
	model.PaginationQ
}

func GetAreaDetailPageByContextNo(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	var body GetAreaDetailPageByContextNoBody
	if !appG.ParseJSONRequest(&body) {
		return
	}

	data, err := model.GetAreaDetailPageByContextNo(body.ContractNo, body.Page, body.PageSize)

	if appG.HasError(err) {

	}

	appG.SuccessResponse(data)

}
