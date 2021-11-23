package router

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	v1 "relocate/api/v1"
	_ "relocate/docs"
	"relocate/middleware"
	"relocate/util/sign"
)

//说明
//middleware.JWT()：管理员、用户共有
//middleware.JWT(sign.AdminClaimsType)：管理员独有
//middleware.JWT(sign.UserClaimsType)：用户独有

//初始化路由信息
func InitRouter() *gin.Engine {
	r := gin.New()
	//全局 Recovery 中间件从任何 panic 恢复，如果出现 panic，它会写一个 500 错误。
	r.Use(gin.Recovery())
	//全局 日志中间件
	r.Use(middleware.LoggerToFile())
	//全局 跨域中间件
	r.Use(middleware.Cors())
	//加载模板文件
	r.LoadHTMLGlob("router/templates/*")
	//加载静态文件
	r.Static("/web", "router/static")
	//swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//v1版本
	apiV1 := r.Group("/api/v1")
	initAdminRouter(apiV1)
	initConfigRouter(apiV1)
	initStagingRouter(apiV1)
	initContractRouter(apiV1)
	initUserRouter(apiV1)
	initTimeRouter(apiV1)
	initHuxingRouter(apiV1)
	initInformationRouter(apiV1)
	initDeclarationRouter(apiV1)
	initAreaDetailRouter(apiV1)
	initResultRouter(apiV1)
	initLoggingRouter(apiV1)
	initAccountingRouter(apiV1)
	initCheckRouter(apiV1)
	initAreaTransferRouter(apiV1)
	return r
}

// todo已完成 合同面积转移
func initAreaTransferRouter(apiV1 *gin.RouterGroup) {
	AreaTransfer := apiV1.Group("/area_transfer")
	{
		// 创建面积转移合同
		AreaTransfer.POST("", middleware.JWT(sign.AdminClaimsType), v1.AreaTransfer)
		// 根据面积合同修改状态
		AreaTransfer.POST("/updateStatus", middleware.JWT(sign.AdminClaimsType), v1.AreaContractUpdateStatus)
	}
}

func initAdminRouter(apiV1 *gin.RouterGroup) {
	admin := apiV1.Group("/admin")
	{
		admin.POST("/login", v1.AdminLogin)
		admin.POST("/update", middleware.JWT(sign.AdminClaimsType), v1.AdminUpdatePassword)
		admin.GET("/get", middleware.JWT(sign.AdminClaimsType), v1.GetAllAdmin)
	}
}

func initConfigRouter(apiV1 *gin.RouterGroup) {
	config := apiV1.Group("/config")
	{
		config.GET("/staging/setting", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.SettingStagingConfig)
		config.GET("/rounds/setting", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.SettingNowRoundsConfig)
		config.GET("/staging/get", v1.GetStagingConfig)
		config.GET("/rounds/get", v1.GetNowRoundsConfig)
		config.POST("/huxing/groupingOptional/setting", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.SettingHuxingGroupingOptionalConfig)
		config.GET("/huxing/groupingOptional/get", v1.GetHuxingGroupingOptionalConfig)
		config.GET("/declaration/active/setting", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.SettingIncludeTotalConfig)
		config.GET("/declaration/active/get", v1.GetIncludeTotalConfig)

		// 查询总可选套数
		config.GET("/huxing/optional/get", v1.GetAllOptionHuxing)
		//设置总可选套数
		config.GET("/huxing/optional/setting", middleware.JWT(sign.AdminClaimsType), v1.SetAllOptionHuxing)
	}
}

func initStagingRouter(apiV1 *gin.RouterGroup) {
	staging := apiV1.Group("/staging")
	{
		staging.POST("/import", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.ImportContractByStaging)
		staging.GET("/getCount", middleware.JWT(sign.AdminClaimsType), v1.GetStagingContractCount)
		// 分页获取分期列表
		staging.GET("/get", v1.GetAllStagingPage)
		// 新增分期
		staging.POST("/new", middleware.JWT(sign.AdminClaimsType), v1.AddStaging)
	}
}

func initContractRouter(apiV1 *gin.RouterGroup) {
	contract := apiV1.Group("/contract")
	{
		contract.POST("/updateHouseWriteOff", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.UpdateHouseWriteOffList)
		contract.POST("/new", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.NewContract)
		contract.POST("/addCardNumber", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.AddCardNumber)
		// 补充导入指标安置面积、计算临迁费面积
		contract.POST("/supplement", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.SupplementContract)
		// 根据分期期数（必须）、姓名、身份证号、合同号、手机号 模糊查询分页查询合同单列表（筛选功能）
		contract.GET("/get", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.QueryContractByFuzzyNameAndStagingId)
		// 根据合同号获取申报中的记录数和结果数(ActiveState = true)
		contract.GET("/getDeclarationCount", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.QueryDeclarationCountByContractOn)
		// 根据分期id，设置该分期的所有合同号的申报状态
		contract.POST("/setCanDeclareByStagingId", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.SetCanDeclareByStagingId)
		// 根据合同号修改数据
		contract.POST("/update", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.UpdateContract)
		//传入一个合同号数组和一个要设置的申报状态（是否可以申报 can_declare），批量修改
		contract.POST("/updateCanDeclare", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.UpdateContractCanDeclareByContractOn)
	}
}

func initUserRouter(apiV1 *gin.RouterGroup) {
	user := apiV1.Group("/user")
	{
		// 根据 phone 获取
		user.GET("/findMyself", middleware.JWT(sign.UserClaimsType), v1.UserQuery)
		// 发送验证码
		user.GET("/get", v1.UserRegisterSendCode)
		// 根据phone 以及验证码验证
		user.POST("/register", v1.UserRegister)
		// 提供电话号码，姓名，证件类型，身份证，身份证正面，身份证反面，正面照片，背面照片
		user.POST("/validator", middleware.JWT(sign.UserClaimsType), v1.UserValidator)
		// 获取验证码(忘记密码)
		user.GET("/forget", v1.UserForgetSendCode)
		// 修改密码
		user.POST("/updatePassword", v1.UpdateUserPassword)
		// 根据身份证获取合同表(用户合同表)
		user.GET("/getContract", middleware.JWT(sign.UserClaimsType), v1.QueryContractByAccountId)
		// 用户证件号查询自己的摇珠结果 card_number
		user.POST("/getResult", middleware.JWT(sign.UserClaimsType), v1.QueryResultByAccountId)
		// 登录(通过手机号码+密码登录)
		user.POST("/login", v1.UserLoginByPhoneAndNumber)
		// 根据用户状态查看核验列表
		user.GET("/getAllValidate", v1.UserQueryByStatusOrFilterName)
		// 修改用户的校验状态
		user.POST("/updateStatus", v1.UpdateUserStatus)
	}
}
func initTimeRouter(apiV1 *gin.RouterGroup) {
	// todo已完成 时间的增删改查
	time := apiV1.Group("/time")
	{
		// 根据分期id 获取现场确认时间
		time.GET("/get", v1.QueryTimeByStagingId)
		// 添加时间
		time.POST("/new", middleware.JWT(sign.AdminClaimsType), v1.AddTime)
		// 修改现场确认时间
		time.POST("/update", middleware.JWT(sign.AdminClaimsType), v1.UpdateTime)
		// 根据传入id删除视
		time.POST("/delete", middleware.JWT(sign.AdminClaimsType), v1.DeleteTimeById)
	}
}

func initHuxingRouter(apiV1 *gin.RouterGroup) {
	huxing := apiV1.Group("/huxing")
	{
		huxing.GET("/get/optional", v1.GetOptionalHuxing)
		// 获取所有户型
		huxing.GET("/get", v1.GetAllHuxing)
		// 删除户型
		huxing.POST("/delete", middleware.JWT(sign.AdminClaimsType), v1.DeleteHuxing)
		// 新增户型
		huxing.POST("/new", middleware.JWT(sign.AdminClaimsType), v1.AddHuxing)
		// 修改户型
		huxing.POST("/update", middleware.JWT(sign.AdminClaimsType), v1.UpdateHuxing)
	}
}

func initInformationRouter(apiV1 *gin.RouterGroup) {
	information := apiV1.Group("/information")
	{
		information.POST("/upload", v1.UploadImage)
		information.GET("/get", middleware.JWT(sign.AdminClaimsType), v1.GetInformationList)
		information.POST("/new", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.SetInformation)
		information.POST("/updateStatus", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.UpdateInformationStatus)
		information.POST("/delete", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.DeleteInformation)

		// 分页获取已发布资讯列表
		information.GET("/getpublic", v1.GetAllPublicInformation)
		// 修改资讯
		information.POST("/update", middleware.JWT(sign.AdminClaimsType), v1.UpdateInformation)
	}
}

func initDeclarationRouter(apiV1 *gin.RouterGroup) {
	declaration := apiV1.Group("/declaration")
	{
		declaration.POST("/add", middleware.JWT(), v1.AddDeclaration)

		declaration.POST("/addNew", middleware.JWT(), v1.AddDeclarationNew)

		declaration.GET("/getAdminName", middleware.JWT(sign.AdminClaimsType), v1.GetAdminName)

		declaration.POST("/updateDeclarationStatus", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.UpdateDeclarationStatus)

		declaration.POST("/printing", middleware.JWT(sign.AdminClaimsType), v1.Printing)

		declaration.GET("/getAll", middleware.JWT(sign.AdminClaimsType), v1.GetAllDeclaration)
		// 删除申报单
		declaration.POST("/delete", middleware.JWT(sign.AdminClaimsType), v1.DeleteDeclarationById)
		// 根据申报表ID更改中签状态（是否录入结果）(管理员)
		declaration.POST("/enterResult", middleware.JWT(sign.AdminClaimsType), v1.DeclarationEnterResult)
		// 根据分期期数（必须）导出申报表列表(管理员)
		declaration.POST("/export", middleware.JWT(sign.AdminClaimsType), v1.ExportDeclarationByStagingId)
		// 根据分期期数（必须）、户型、状态（申报状态，摇珠状态，有效状态）、(身份证号、合同号、手机号)查询分页查询申报表列表（筛选功能）(管理员)
		declaration.GET("/get", middleware.JWT(sign.AdminClaimsType), v1.GetAllDeclarationFuzzy)
		// 根据合同号查看该合同单的申报
		declaration.GET("/getDeclaration", middleware.JWT(), v1.GetDeclarationByContractNo)
		// 根据合同号，分期数ID查看申报表详情
		declaration.GET("/getDetail", middleware.JWT(), v1.GetDeclarationByContractNoAndStringId)
		// updateActive 修改申报有效状态，根据申报ID
		declaration.POST("/updateActive", middleware.JWT(sign.AdminClaimsType), v1.UpdateDeclarationActiveStatus)
		// 根据申报表ID更改申报表数据
		// todo已测试
		declaration.POST("/updateDeclaration", middleware.JWT(), v1.UpdateDeclarationTrustee)
	}
}

func initAreaDetailRouter(apiV1 *gin.RouterGroup) {
	areaDetail := apiV1.Group("/areaDetail")
	{
		//分页获取合同号面积明细记录
		areaDetail.GET("/get", v1.GetAreaDetailPageByContextNo)
	}
}

func initResultRouter(apiV1 *gin.RouterGroup) {
	result := apiV1.Group("/result")
	{
		result.GET("/get", v1.GetResultList)
		result.POST("/export", middleware.JWT(sign.AdminClaimsType), v1.ExportResults)
		// 根据公示状态分页获取摇珠结果列表(移动端)
		result.GET("/getByStatus", middleware.JWT(sign.AdminClaimsType), v1.GetResultByStatus)
		// 根据申报号(批量)设置公示状态 { "declaration_id":[31,9], "publicity_status":false }
		result.POST("/updatePublicityStat", middleware.JWT(sign.AdminClaimsType), v1.UpdateResultPublicStat)
	}
}

func initLoggingRouter(apiV1 *gin.RouterGroup) {
	logging := apiV1.Group("/logging")
	{
		// 仅分页
		logging.GET("/getAll", middleware.JWT(sign.AdminClaimsType), v1.GetAllLogging)

		// 根据操作人(非必须)、操作 、分页获取日志数据列表(管理员)(模糊查询)
		logging.GET("/get", middleware.JWT(sign.AdminClaimsType), v1.GetLogging)
	}
}

func initAccountingRouter(apiV1 *gin.RouterGroup) {
	accounting := apiV1.Group("/accounting")
	{
		//根据合同号核算
		accounting.POST("", middleware.JWT(sign.AdminClaimsType), middleware.Visitors(true), v1.AddAccounting)
		// 根据合同号、被拆迁人模糊查询分页查询核算列表（筛选功能）(管理员)
		accounting.GET("", middleware.JWT(sign.AdminClaimsType), v1.GetAccountingList)

		accounting.GET("/export", middleware.JWT(sign.AdminClaimsType), v1.ExportAccounting)
	}
}

func initCheckRouter(apiV1 *gin.RouterGroup) {
	check := apiV1.Group("/check")
	{
		// 根据合同号 ContractNo 、被拆迁人 Peoples 模糊查询分页查询核算列表（筛选功能）（管理员）
		check.GET("", middleware.JWT(sign.AdminClaimsType), v1.CheckFuzzyQueryByContractNoOrPeoples)
		// 根据结果核算
		check.POST("", middleware.JWT(), v1.AccountingByResults)
		// 导出核算表(管理员)
		check.GET("/export", middleware.JWT(), v1.ExportCheck)

	}
}
