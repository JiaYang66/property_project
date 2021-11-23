package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"relocate/api"
	"relocate/middleware"
	"relocate/model"
	"relocate/service/declaration_service"
	"relocate/util/app"
	"relocate/util/sign"
	"strconv"
)

//TODO已完成: 根据分期期数（必须）、身份证号、合同号、手机号 模糊查询分页查询申报表列表（筛选功能）(管理员)
//TODO已完成: 根据合同号查看申报表详情
//TODO已完成: 根据合同号新增申报(操作人 管理员姓名、登录人姓名)
//TODO已完成: 根据合同号更改申报状态(管理员)
//TODO已完成: 申报表打印管理员姓名(管理员)
//TODO已完成: 根据合同号更改中签状态（是否录入结果）(管理员)

// @Tags 申报单
// @Summary 根据合同号新增申报(操作人 管理员姓名、登录人姓名)
// @Description 根据合同号新增申报(操作人 管理员姓名、登录人姓名)
// @Produce  json
// @Security ApiKeyAuth
// @Param data body api.AddDeclarationBody true "申报信息"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/declaration/addNew [post]
func AddDeclarationNew(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body api.AddDeclarationBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	//判断是管理员登录还是用户登录 -- true：管理员 false：用户
	claim := middleware.GetClaims(c)
	ok := false
	if claim.Type == sign.AdminClaimsType {
		ok = true
	}
	if appG.HasError(declaration_service.AddDeclarationNew(body, claim.Issuer, ok)) {
		return
	}
	appG.SuccessResponse("申报成功")
}

// @Tags 申报单
// @Summary 根据合同号新增申报(操作人 管理员姓名、登录人姓名)
// @Description 根据合同号新增申报(操作人 管理员姓名、登录人姓名)
// @Produce  json
// @Security ApiKeyAuth
// @Param data body api.AddDeclarationBody true "申报信息"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/declaration/add [post]
func AddDeclaration(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body api.AddDeclarationBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	//判断是管理员登录还是用户登录 -- true：管理员 false：用户
	claim := middleware.GetClaims(c)
	ok := false
	if claim.Type == sign.AdminClaimsType {
		ok = true
	}
	if appG.HasError(declaration_service.AddDeclarationNew(body, claim.Issuer, ok)) {
		return
	}
	appG.SuccessResponse("申报成功")
}

// @Tags 申报单
// @Summary 获取所有的申报表列表
// @Description 获取所有的申报表列表
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/declaration/getAll [get]
func GetAllDeclaration(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	data, err := model.GetAllDeclaration()
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(data)
}

type QueryDeclarationBody struct {
	ContractNo    string `json:"contract_no" form:"contract_no" validate:"required"`
	StagingId     string `json:"staging_id" form:"staging_id" validate:"required"`
	DeclarationID int    `json:"declaration_id" form:"declaration_id" validate:"required"`
}

type QueryDeclarationByContractNoBody struct {
	ContractNo string `json:"contract_no" form:"contract_no" validate:"required"`
}

// @Tags 申报单
// @Summary 申报表打印管理员姓名(管理员)
// @Description 申报表打印管理员姓名(管理员)
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/declaration/getAdminName [get]
func GetAdminName(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	claim := middleware.GetClaims(c)
	admin, err := model.GetAdminInfo(claim.Issuer)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(map[string]interface{}{
		"name": admin.AdminName,
	})
}

type UpdateStatusBody struct {
	DeclarationID uint `json:"declaration_id" form:"declaration_id"`
	Status        int  `json:"status" form:"status" validate:"oneof=0 1"`
}

// @Tags 申报单
// @Summary 根据申报表ID更改申报状态(管理员)
// @Description 根据申报表ID更改申报状态(管理员)
// @Produce  json
// @Security ApiKeyAuth
// @Param declaration_id query int true "申报ID（必须）"
// @Param status query int true "申报表申报状态（必须）-- 0：表示进行中；1：表示已确定"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/declaration/updateDeclarationStatus [post]
func UpdateDeclarationStatus(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body UpdateStatusBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	claim := middleware.GetClaims(c)
	if appG.HasError(declaration_service.ChangeDeclarationStatus(body.DeclarationID, body.Status, claim.Issuer)) {
		return
	}
	appG.SuccessResponse("修改成功")
}

type PrintingBody struct {
	DeclarationID uint `json:"declaration_id" form:"declaration_id" validate:"required"`
}

// @Tags 申报单
// @Summary 打印申报表(管理员)
// @Produce  json
// @Security ApiKeyAuth
// @Param declaration_id query int true "申报表ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/declaration/printing [post]
func Printing(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body PrintingBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	claim := middleware.GetClaims(c)
	admin, err := model.GetAdminInfo(claim.Issuer)
	if appG.HasError(err) {
		return
	}
	file, filename, err := declaration_service.GenerateExcelNew(body.DeclarationID, admin.AdminSignname)
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

type DeleteResultBody struct {
	DeclarationID uint `json:"declaration_id" form:"declaration_id" validate:"required"`
}

type UpdateDeclarationActiveBody struct {
	DeclarationID uint `json:"declaration_id" form:"declaration_id" validate:"required"`
	State         bool `json:"state" form:"state" validate:"omitempty,required"`
}

func DeleteDeclarationById(context *gin.Context) {
	// 删除申报
	appG := app.Gin{Ctx: context}
	var body DeleteResultBody
	if !appG.ParseQueryRequest(&body) {
		return
	}

	err := declaration_service.DeleteDeclarationById(body.DeclarationID)

	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse("删除成功")
}

func DeclarationEnterResult(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	var body api.DeclarationEnterResultBody
	if !appG.ParseJSONRequest(&body) {
		return
	}

	claim := middleware.GetClaims(context)
	err := declaration_service.DeclarationEnterResult(body, claim.Issuer)
	if appG.HasError(err) {
		return
	}

	appG.SuccessResponse("修改成功")
}

type ExportDeclarationByStagingIdBody struct {
	StagingId uint `json:"stagingId" form:"stagingId" validate:"required"`
}

func ExportDeclarationByStagingId(context *gin.Context) {
	appG := app.Gin{context}
	var body ExportDeclarationByStagingIdBody
	if !appG.ParseQueryRequest(&body) {
		return
	}

	data, err := model.GetAllDeclarationByStagingId(body.StagingId)

	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(data)
}

func GetAllDeclarationFuzzy(context *gin.Context) {
	appG := app.Gin{context}
	var body api.GetAllDeclarationFuzzyBody
	{
		body.HuxingId = -1
		body.TimeId = -1
		body.DeclarationStatus = -1
		body.WinningStatus = -1
		body.ActiveState = -1
	}
	if !appG.ParseQueryRequest(&body) {
		return
	}

	data, err := declaration_service.GetAllDeclarationFuzzy(body)

	if appG.HasError(err) {
		return
	}

	appG.SuccessResponse(data)
}

func GetDeclarationByContractNo(context *gin.Context) {
	appG := app.Gin{context}
	var body QueryDeclarationByContractNoBody
	if !appG.ParseQueryRequest(&body) {
		return
	}

	data, err := model.GetDeclarationByContractNo(body.ContractNo)
	if appG.HasError(err) {
		return
	}

	appG.SuccessResponse(data)
}

type GetDeclarationByContractNoAndStringIdBody struct {
	ContractNo string `json:"contract_no" form:"contract_no" validate:"required"`
	StagingId  uint   `json:"staging_id" form:"staging_id" validate:"required"`
}

func GetDeclarationByContractNoAndStringId(context *gin.Context) {
	appG := app.Gin{context}
	var body GetDeclarationByContractNoAndStringIdBody
	if !appG.ParseQueryRequest(&body) {
		return
	}

	allList, currentList, err := model.GetDeclarationByContractNoAndStagingId(body.ContractNo, body.StagingId)
	if appG.HasError(err) {
		return
	}

	data := make(map[string]interface{})
	data["declaration_list"] = allList
	data["current_declaration"] = currentList

	appG.SuccessResponse(data)
}

type UpdateDeclarationActiveStatusBody struct {
	DeclarationId uint `json:"declaration_id" form:"declaration_id" validate:"required"`
	State         bool `json:"state" form:"state" validate:"omitempty,required"`
}

func UpdateDeclarationActiveStatus(context *gin.Context) {
	appG := app.Gin{context}
	var body UpdateDeclarationActiveStatusBody
	if !appG.ParseQueryRequest(&body) {
		return
	}

	claim := middleware.GetClaims(context)
	err := declaration_service.UpdateDeclarationActiveStatus(body.DeclarationId, body.State, claim.Issuer)
	if appG.HasError(err) {
		return
	}

	appG.SuccessResponse("修改成功")
}

func UpdateDeclarationTrustee(context *gin.Context) {
	appG := app.Gin{context}
	var body api.UpdateDeclaration
	if !appG.ParseJSONRequest(&body) {
		return
	}

	claim := middleware.GetClaims(context)
	ok := false
	if claim.Type == sign.AdminClaimsType {
		ok = true
	}

	err := declaration_service.UpdateDeclarationTrustee(body, claim.Issuer, ok)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse("修改成功")
}
