/**
* @Author: SAmwake
* @Date: 2021/11/12 17:05
 */
package v1

import (
	"github.com/gin-gonic/gin"
	"relocate/middleware"
	"relocate/service/area_service"
	"relocate/util/app"
)

type AreaTransferBody struct {
	ContractNo       string  `json:"contract_no" form:"contract_no" validate:"required"`
	TargetContractNo string  `json:"target_contract_no" form:"target_contract_no" validate:"required"`
	Area             float64 `json:"area" form:"area" validate:"required"`
	AreaUnitPrice    float64 `json:"area_unit_price" form:"area_unit_price" validate:"required"`
}

func AreaTransfer(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	var body AreaTransferBody
	if !appG.ParseJSONRequest(&body) {
		return
	}

	err := area_service.AreaTransfer(body.ContractNo, body.TargetContractNo, body.Area, body.AreaUnitPrice)
	if appG.HasError(err) {
		return
	}

	appG.SuccessResponse("创建合同成功")
}

type AreaContractUpdateStatusBody struct {
	Id     uint `json:"id" form:"id" validate:"required"`
	Status bool `json:"status" form:"statue" validate:"omitempty,required"`
}

func AreaContractUpdateStatus(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	var body AreaContractUpdateStatusBody
	if !appG.ParseJSONRequest(&body) {
		return
	}

	claim := middleware.GetClaims(context)
	err := area_service.AreaContractUpdateStatus(body.Id, body.Status, claim.Issuer)
	if appG.HasError(err) {
		return
	}

	appG.SuccessResponse("状态修改成功")
}
