package v1

//TODO已完成: 新增资讯(管理员)
//TODO已完成: 分页获取资讯列表
//TODO已完成: 根据资讯ID修改资讯(管理员)
//TODO已完成: 根据资讯ID删除资讯(管理员)
//TODO已完成: 根据资讯ID更改发布状态(管理员)

import (
	"github.com/gin-gonic/gin"
	"relocate/api"
	"relocate/model"
	"relocate/service/information_server"
	"relocate/util/app"
	"strconv"
	"time"
)

type NewInformationBody struct {
	Title   string `json:"title" form:"title" validate:"required"`     //标题
	Content string `json:"content" form:"content" validate:"required"` //内容
	Picture string `json:"picture" form:"picture"`                     //图片
	Suffix  string `json:"suffix" form:"suffix"`                       //图片后缀
}

// @Tags 资讯
// @Summary 新增资讯
// @Description 后台管理员新增资讯
// @Produce  json
// @Security ApiKeyAuth
// @Param data body NewInformationBody true "资讯标题（必须），内容（必须），图片(base64字符串),图片后缀(例如.jpg)"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/information/new [post]
func SetInformation(ctx *gin.Context) {
	appG := app.Gin{Ctx: ctx}
	var body NewInformationBody
	if !appG.ParseJSONRequest(&body) {
		return
	}

	//name := time.Now().Unix()
	//filename := string(name) + ".png"
	//if err := ctx.SaveUploadedFile(body.Picture, "ImageServer/"+filename); err != nil {
	//	appG.HasError(err)
	//	return
	//}
	//photo := "ImageServer/" + filename
	if appG.HasError(information_server.AddInformation(body.Title, body.Content, body.Picture, body.Suffix)) {
		return
	}
	appG.SuccessResponse("新增成功")
}

type DeleteInformationBody struct {
	Id uint `json:"id" form:"id" validate:"required"`
}

// @Tags 资讯
// @Summary 删除资讯
// @Description 后台管理员删除资讯
// @Produce  json
// @Security ApiKeyAuth
// @Param data body DeleteInformationBody true "资讯id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/information/delete [post]
func DeleteInformation(ctx *gin.Context) {
	appG := app.Gin{Ctx: ctx}
	var body DeleteInformationBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	if appG.HasError(information_server.DeleteInformation(body.Id)) {
		return
	}
	appG.SuccessResponse("删除成功")
}

type UpdateInformationBody struct {
	Id      uint   `json:"id" form:"id" validate:"required"`
	Title   string `json:"title" form:"title" `     //标题
	Content string `json:"content" form:"content" ` //内容
	Picture string `json:"picture" form:"picture"`  //图片
	Suffix  string `json:"suffix" form:"suffix"`    //图片后缀
}

func UpdateInformation(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	var body UpdateInformationBody
	body.Id = 0
	if !appG.ParseJSONRequest(&body) {
		return
	}

	if err := information_server.UpdateInformation(body.Id, body.Content, body.Picture, body.Suffix, body.Title); err != nil {
		appG.BadResponse(err.Error())
		return
	}
	appG.SuccessResponse("修改成功")
}

type UpdateInformationStatusBody struct {
	Id     uint `json:"id" form:"id" validate:"required"`
	Status int  `json:"status" form:"status"`
}

// @Tags 资讯
// @Summary 修改资讯状态
// @Description 后台管理员修改资讯状态
// @Produce  json
// @Security ApiKeyAuth
// @Param data body UpdateInformationStatusBody true "资讯id(必须),status状态int（0待发布中,1已发布）"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/information/updateStatus [post]
func UpdateInformationStatus(ctx *gin.Context) {
	appG := app.Gin{Ctx: ctx}
	var body UpdateInformationStatusBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	if appG.HasError(information_server.UpdateInformationStatus(body.Id, model.InformationStatus(body.Status))) {
		return
	}
	appG.SuccessResponse("修改成功")
}

// @Tags 资讯
// @Summary 分页获取资讯列表
// @Description 分页获取资讯列表
// @Produce  json
// @Security ApiKeyAuth
// @Param page query int false "页码"
// @Param pageSize query int false "页面大小"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/information/get [get]
func GetInformationList(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body api.PaginationQueryBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	informationList, err := model.GetInformations(body.Page, body.PageSize)

	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse(informationList)
}

type UploadImageBody struct {
	Images string `json:"images" form:"images" validate:"required"` //base64图片
	Suffix string `json:"suffix" form:"suffix" validate:"required"` //图片后缀
}

// @Tags 资讯
// @Summary 上传图片，返回图片链接
// @Description 上传图片，返回图片链接
// @Produce  json
// @Param data body UploadImageBody true "上传图片的信息"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/information/upload [post]
func UploadImage(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body UploadImageBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	formatTime := strconv.FormatInt(time.Now().Unix(), 10)
	formatImg, err := api.UploadImg(body.Images, body.Suffix, "information", formatTime)
	if err != nil {
		return
	}
	appG.SuccessResponse(map[string]interface{}{"Image": formatImg})
}

// GetAllPublicInformation 分页获取已发布资讯列表
func GetAllPublicInformation(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	var body api.PaginationQueryBody
	if !appG.ParseQueryRequest(&body) {
		return
	}

	data, err := model.GetAllPublicInformation(body.Page, body.PageSize)
	if err != nil {
		appG.BadResponse(err.Error())
		return
	}

	appG.SuccessResponse(data)
}
