package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"relocate/model"
	"relocate/service/user_service"
	"relocate/util/app"
	"relocate/util/gredis"
	"relocate/util/random"
	"strings"
)

//TODO已完成: 用户通过手机验证码注册
//TODO已完成: 用户通过手机号码+密码登录
//TODO已完成: 获取用户名下合同号(用户合同表)
//TODO已完成: 用户进行身份核验
//TODO已完成: 用户通过验证码修改密码

type UserGetCodeBody struct {
	PhoneNumber string `json:"phone_number" form:"phone_number" validate:"required,checkMobile,len=11"`
}

// @Tags 用户
// @Summary 获取验证码(注册)
// @Produce json
// @Param phone_number query string true "手机号码"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/user/get [get]
func UserRegisterSendCode(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body UserGetCodeBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	_, err := model.GetUserInfo(body.PhoneNumber)
	if err == nil {
		appG.BadResponse("用户已存在")
		return
	}
	code := random.Code(6)
	if !random.SendCode(body.PhoneNumber, code) {
		appG.BadResponse("发送验证码错误")
		return
	}
	user := make(map[string]string)
	user["user-register-"+body.PhoneNumber] = code
	gredis.Set(user, 200)
	appG.SuccessResponse("发送验证码成功")
}

type UserRegisterBody struct {
	Username string `json:"username" validate:"required,len=11,checkMobile" `
	Password string `json:"password" validate:"required"`
	Code     string `json:"code" validate:"required,len=6"` //验证码为6位数字
}

// @Tags 用户
// @Summary 注册(通过手机号码+密码注册+验证码)
// @Produce json
// @Param data body UserRegisterBody true "注册信息"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/user/register [post]
func UserRegister(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body UserRegisterBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	_, err := model.GetUserInfo(body.Username)
	if err == nil {
		appG.BadResponse("用户已存在")
		return
	}
	code, _ := gredis.Get("user-register-" + body.Username)
	if code != body.Code { //用户没有获取验证码或者验证码过期
		appG.BadResponse("验证码错误")
		return
	}
	if appG.HasError(user_service.CreateUser(body.Username, body.Password)) {
		return
	}
	gredis.Delete("user-register-" + body.Username)
	appG.SuccessResponse("注册成功")
}

type UserLoginBody struct {
	Username string `json:"username" form:"username" validate:"required,checkMobile,len=11"`
	Password string `json:"password" form:"password" validate:"required"`
}

type UserContractBody struct {
	CardNumber string `json:"card_number" form:"card_number"`
}

type UserInfoBody struct {
	PhoneNumber   string `json:"phone_number" form:"phone_number" validate:"required,checkMobile,len=11"` //当前登录用户的手机号码
	Peoples       string `json:"peoples" form:"peoples" validate:"required"`                              //真实姓名
	IdNumberType  int    `json:"id_number_type" form:"id_number_type"`                                    //证件类型 0大陆身份证 1香港身份证 2护照
	CardNumber    string `json:"card_number" form:"card_number" validate:"required"`                      //身份证
	PositiveImage string `json:"positive_image" form:"positive_image" validate:"required"`                //身份证正面,要求上传的图片的base64格式
	NegativeImage string `json:"negative_image" form:"negative_image" validate:"required"`                //身份证反面,要求上传的图片的base64格式
	SuffixA       string `json:"suffix_a" form:"suffix_a" validate:"required"`                            //正面图片原先的后缀(例：.jpg .png等)
	SuffixB       string `json:"suffix_b" form:"suffix_b" validate:"required"`                            //背面图片原先的后缀(例：.jpg .png等)
}

// @Tags 用户
// @Summary 用户进行身份核验
// @Produce json
// @Security ApiKeyAuth
// @Param data body UserInfoBody true "用户身份信息"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/user/validator [post]
func UserValidator(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body UserInfoBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	// 中文换英文括号
	cardNumber := strings.Replace(body.CardNumber, "（", "(", -1)
	cardNumber = strings.Replace(cardNumber, "）", ")", -1)
	//更新
	if appG.HasError(user_service.UpdateUser(body.IdNumberType, body.PhoneNumber, body.Peoples, cardNumber, body.PositiveImage, body.NegativeImage, body.SuffixA, body.SuffixB)) {
		return
	}
	//自动匹配
	if model.AutoMatchUser(body.PhoneNumber, body.Peoples, cardNumber) {
		//修改状态
		if appG.HasError(user_service.UpdateStatus(body.PhoneNumber)) {
			return
		}
	}
	appG.SuccessResponse("提交成功")
}

type UserBody struct {
	Phone string `json:"phone" form:"phone" validate:"required,checkMobile,len=11"`
}

// @Tags 用户
// @Summary 用户通过手机查看自己的身份信息
// @Produce json
// @Security ApiKeyAuth
// @Param phone query string true "用户通过手机号查看身份信息"  前提是通过核验
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/user/findMyself [get]
func UserQuery(c *gin.Context) {
	appG := app.Gin{Ctx: c}
	var body UserBody
	if !appG.ParseQueryRequest(&body) {
		return
	}
	date, err := model.GetUserInfo(body.Phone)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.BadResponse("账号不存在")
			return
		}
		appG.HasError(err)
		return
	}
	appG.SuccessResponse(map[string]interface{}{"user": date})
}

type UpdateStatusByPhoneNumber struct {
	PhoneNumber string `json:"phone_number" form:"phone_number" validate:"required,checkMobile,len=11"`
	Pass        bool   `json:"pass" form:"phone_number"`
}

// UserForgetSendCode 获取验证码(忘记密码)
func UserForgetSendCode(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	var body UserGetCodeBody
	if !appG.ParseQueryRequest(&body) {
		return
	}

	if err := user_service.UserForgetSendCode(body.PhoneNumber); err != nil {
		appG.BadResponse(err)
		return
	}
	appG.SuccessResponse("发送验证码成功")
}

// UpdateUserPassword 修改密码  传入手机号，新密码，验证码
func UpdateUserPassword(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	var body UserRegisterBody
	if !appG.ParseJSONRequest(&body) {
		return
	}

	if err := user_service.UpdateUserPassword(body.Username, body.Password, body.Code); err != nil {
		appG.BadResponse(err)
		return
	}

	appG.SuccessResponse("修改密码成功")
}

type QueryContractByAccountIdBody struct {
	CardNumber string `json:"card_number" validate:"required" form:"card_number"` //验证码为6位数字
}

// QueryContractByAccountId 根据身份证获取合同表(用户合同表)
// 测试 440440111122223333
func QueryContractByAccountId(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	var body QueryContractByAccountIdBody
	if !appG.ParseQueryRequest(&body) {
		return
	}

	contracts, err := user_service.QueryContractByAccountId(body.CardNumber)
	if err != nil {
		appG.BadResponse(err.Error())
		return
	}
	result := make(map[string]interface{}, 0)
	result["contractList"] = contracts
	result["cardNumber"] = body.CardNumber
	appG.SuccessResponse(result)
}

// QueryResultByAccountId 用户证件号查询自己的摇珠结果
//id_number 440440111122223333
func QueryResultByAccountId(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	var body QueryContractByAccountIdBody
	if !appG.ParseQueryRequest(&body) {
		return
	}

	results, err := user_service.QueryResultByAccountId(body.CardNumber)
	if err != nil {
		appG.BadResponse(err.Error())
		return
	}

	appG.SuccessResponse(results)
}

// UserLoginByPhoneAndNumber 登录(通过手机号码+密码登录)
func UserLoginByPhoneAndNumber(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	var body UserLoginBody
	if !appG.ParseJSONRequest(&body) {
		return
	}

	token, err := user_service.GenerateToken(body.Username, body.Password)
	if appG.HasError(err) {
		return
	}
	user, err := model.FindUserByPhone(body.Username)
	if appG.HasError(err) {
		return
	}

	result := make(map[string]interface{})
	result["token"] = token
	result["phoneNumber"] = user.PhoneNumber
	result["name"] = user.Name
	result["idNumber"] = user.IdNumber
	result["userStatus"] = user.UserStatus
	appG.SuccessResponse(result)
}

type UserQueryByStatusOrFilterNameBody struct {
	FilterStatus string `json:"filter_status" form:"filter_status"` // 筛选状态，不选查全部
	FilterName   string `json:"filter_name" form:"filter_name"`     // 姓名，身份证或电话号码，不选查全部
	Page         uint   `json:"page"`
	PageSize     uint   `json:"page_size"`
}

// UserQueryByStatusOrFilterName 根据用户状态查看核验列表
func UserQueryByStatusOrFilterName(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	var body UserQueryByStatusOrFilterNameBody
	if !appG.ParseQueryRequest(&body) {
		return
	}

	data, err := user_service.UserQueryByStatusOrFilterName(body.FilterStatus, body.FilterName, body.Page, body.PageSize)

	if err != nil {
		appG.BadResponse(err.Error())
		return
	}

	appG.SuccessResponse(data)

}

type UpdateUserStatusBody struct {
	PhoneNumber string `json:"phone_number" form:"phone_number" validate:"required,checkMobile,len=11"`
	pass        bool   `json:"pass" form:"pass"`
}

// UpdateUserStatus 修改用户的校验状态
func UpdateUserStatus(context *gin.Context) {
	appG := app.Gin{Ctx: context}
	var body UpdateUserStatusBody
	if !appG.ParseJSONRequest(&body) {
		return
	}
	fmt.Println(body)
	err := user_service.UpdateUserStatus(body.PhoneNumber, body.pass)
	if appG.HasError(err) {
		return
	}
	appG.SuccessResponse("修改成功")

}
