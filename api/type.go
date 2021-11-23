package api

type PaginationQueryBody struct {
	Page     uint `json:"page" form:"page"`
	PageSize uint `json:"pageSize" form:"pageSize"`
}

type ContractFilterBody struct {
	StagingId  uint   `json:"stagingId" form:"stagingId" validate:"required"`
	FilterName string `json:"filterName" form:"filterName"`
	Page       uint   `json:"page" form:"page"`
	PageSize   uint   `json:"pageSize" form:"pageSize"`
}

type UserValidateFilter struct {
	FilterName   string `json:"filter_name" form:"filter_name"`
	FilterStatus int    `json:"filter_status" form:"filter_status"`
	Page         uint   `json:"page" form:"page"`
	PageSize     uint   `json:"pageSize" form:"pageSize"`
}

type BaseContractBody struct {
	StagingID                         uint   `json:"staging_id" form:"staging_id" validate:"required"`                                 //分期id(必须)
	ContractNo                        string `json:"contract_no" form:"contract_no" validate:"required"`                               //合同号(必须)
	SocialCategory                    string `json:"social_category" form:"social_category" `                                          //社别
	HouseNumber                       string `json:"house_number" form:"house_number" validate:"required"`                             //房屋栋号(必须)
	OldAddress                        string `json:"old_address" form:"old_address" validate:"required"`                               //被拆迁房屋地址(必须)
	DateOfSigning                     string `json:"date_of_signing" form:"date_of_signing" validate:"required"`                       //签署协议日期(必须)
	DateOfDelivery                    string `json:"date_of_delivery" form:"date_of_delivery" validate:"required"`                     //交楼日期(必须)
	Signatory                         string `json:"signatory" form:"signatory" validate:"required"`                                   //签约人(必须)
	Registration                      string `json:"registration" form:"registration"`                                                 //登记字号
	CanDeclare                        bool   `json:"can_declare" form:"can_declare"`                                                   //可否申报
	CollectiveLandPropertyCertificate string `json:"collective_land_property_certificate" form:"collective_land_property_certificate"` //集体土地房产证字
	Proprietor                        string `json:"proprietor" form:"proprietor"`                                                     //证载产权人
	ChangeMethod                      string `json:"change_method" form:"change_method"`                                               //变更方式
	PhoneNumber1                      string `json:"phone_number1" form:"phone_number1"`                                               //手机号码1
	PhoneNumber2                      string `json:"phone_number2" form:"phone_number2"`
	Desc                              string `json:"desc" form:"desc"`                       //备注
	HouseWriteOff                     bool   `json:"house_write_off" form:"house_write_off"` //是否完成房屋注销
}

type AddContractBody struct {
	StagingID                         uint    `json:"staging_id" form:"staging_id" validate:"required"`                                 //分期id(必须)
	ContractNo                        string  `json:"contract_no" form:"contract_no" validate:"required"`                               //合同号(必须)
	SocialCategory                    string  `json:"social_category" form:"social_category"  validate:"required"`                      //社别(必须)
	Peoples                           string  `json:"peoples" form:"peoples" validate:"required"`                                       //被拆迁人(至少一个，可能有多个)
	HouseNumber                       string  `json:"house_number" form:"house_number" validate:"required"`                             //房屋栋号(必须)
	CardNumber                        string  `json:"card_number" form:"card_number" validate:"required"`                               //身份证号码(必须，可能有多个)
	OldAddress                        string  `json:"old_address" form:"old_address" validate:"required"`                               //被拆迁的地址(必须)
	DateOfSigning                     string  `json:"date_of_signing" form:"date_of_signing" validate:"required"`                       //签署协议日期(必须)
	DateOfDelivery                    string  `json:"date_of_delivery" form:"date_of_delivery" validate:"required"`                     //交楼日期(必须)
	Signatory                         string  `json:"signatory" form:"signatory" validate:"required"`                                   //签约人(必须)
	Registration                      string  `json:"registration" form:"registration"`                                                 //登记字号
	InitialHQArea                     float64 `json:"initial_hq_area" form:"initial_hq_area" validate:"gt=0"`                           //初始回迁面积(必须)
	HouseWriteOff                     bool    `json:"house_write_off" form:"house_write_off"`                                           //是否完成房屋注销
	IsDelivery                        bool    `json:"is_delivery" form:"is_delivery" validate:"eq=1"`                                   //是否交齐楼（必须为1）
	CanDeclare                        bool    `json:"can_declare" form:"can_declare"`                                                   //是否可申报(不写默认不可，1为可)
	PhoneNumber1                      string  `json:"phone_number1" form:"phone_number1"`                                               //手机号码1
	PhoneNumber2                      string  `json:"phone_number2" form:"phone_number2"`                                               //手机号码2
	Proprietor                        string  `json:"proprietor" form:"proprietor"`                                                     //证载产权人
	ChangeMethod                      string  `json:"change_method" form:"change_method"`                                               //变更方式 	//受托人关系
	Desc                              string  `json:"desc" form:"desc"`                                                                 //备注
	CollectiveLandPropertyCertificate string  `json:"collective_land_property_certificate" form:"collective_land_property_certificate"` //集体土地房产证字
}

type DeclarationFilterBody struct {
	StagingId         uint   `json:"stagingId" form:"stagingId" validate:"required"`
	HuxingId          uint   `json:"huxing_id" form:"huxing_id"`
	TimeId            uint   `json:"time_id" form:"time_id"`
	DeclarationStatus int    `json:"declaration_status" form:"declaration_status"`
	WinningStatus     int    `json:"winning_status" form:"winning_status"`
	ActiveState       int    `json:"active_state" form:"active_state"`
	FilterName        string `json:"filterName" form:"filterName"`
	Page              uint   `json:"page" form:"page"`
	PageSize          uint   `json:"pageSize" form:"pageSize"`
}

type AddDeclarationBody struct {
	ContractNo          string `json:"contract_no" form:"contract_no" validate:"required"` //合同号
	HuxingID            uint   `json:"huxing_id" form:"huxing_id"`                         //户型id
	TimeID              uint   `json:"time_id" form:"time_id"`                             //现场确认时段id
	Trustee             string `json:"trustee" form:"trustee"`                             //受托人
	TrusteePhoneNumber  string `json:"trustee_phone_number" form:"trustee_phone_number"`   //受托人手机号码
	TrusteeCardNumber   string `json:"trustee_card_number" form:"trustee_card_number"`     //受托人身份证号码
	TrusteeRelationship string `json:"trustee_relationship" form:"trustee_relationship"`   //受托人关系
}
type HuxingGroupingOptionalConfigJson struct {
	ID  int `json:"id" form:"id"`
	Num int `json:"num" form:"num"`
}
type AddDeclarationNewBody struct {
	ContractNo          string                             `json:"contract_no" form:"contract_no" validate:"required"` //合同号
	Config              []HuxingGroupingOptionalConfigJson `json:"huxing_config" form:"huxing_config"`
	TimeID              uint                               `json:"time_id" form:"time_id"`                           //现场确认时段id
	Trustee             string                             `json:"trustee" form:"trustee"`                           //受托人
	TrusteePhoneNumber  string                             `json:"trustee_phone_number" form:"trustee_phone_number"` //受托人手机号码
	TrusteeCardNumber   string                             `json:"trustee_card_number" form:"trustee_card_number"`   //受托人身份证号码
	TrusteeRelationship string                             `json:"trustee_relationship" form:"trustee_relationship"` //受托人关系
}

type AreaDetailFilterBody struct {
	ContractNo string `json:"contract_no" form:"contract_no"`
	Page       uint   `json:"page" form:"page"`
	PageSize   uint   `json:"pageSize" form:"pageSize"`
}

type ResultBody struct {
	DeclarationID uint   `json:"declaration_id" form:"declaration_id" validate:"required"` //申报表ID
	Status        int    `json:"status" form:"status" validate:"oneof=0 1"`                //状态 0：不中签；1：中签
	BuildingNo    string `json:"building_no" form:"building_no"`                           //楼栋名称
	RoomNo        string `json:"room_no" form:"room_no"`                                   //房间号
}

type UpdateDeclaration struct {
	DeclarationID       uint   `json:"declaration_id" form:"declaration_id"`              //申报表ID
	HuxingID            uint   `json:"huxing_id" form:"huxing_id"`                        //户型id
	Trustee             string `json:"trustee" form:"trustee"`                            //受托人
	TrusteePhoneNumber  string `json:"trustee_phone_number" form:"trustee_phone_number" ` //受托人手机号码
	TrusteeCardNumber   string `json:"trustee_card_number" form:"trustee_card_number"`    //受托人身份证号码
	TrusteeRelationship string `json:"trustee_relationship" form:"trustee_relationship"`  //受托人关系
}

type ResultFilterBody struct {
	StagingId       uint   `json:"stagingId" form:"stagingId" validate:"required"`
	HuxingId        uint   `json:"huxing_id" form:"huxing_id"`
	PublicityStatus int    `json:"publicity_status" form:"publicity_status"`
	FilterName      string `json:"filterName" form:"filterName"`
	Page            uint   `json:"page" form:"page"`
	PageSize        uint   `json:"pageSize" form:"pageSize"`
}

type LoggingFilterBody struct {
	Operator  string `json:"operator" form:"operator"`
	Operation string `json:"operation" form:"operation"`
	Page      uint   `json:"page" form:"page"`
	PageSize  uint   `json:"pageSize" form:"pageSize"`
}

type NewHuxingBody struct {
	StagingID  uint   `json:"staging_id" form:"staging_id"`                       //分期ID
	BuildingNo string `json:"building_no" form:"building_no" validate:"required"` //栋号
	HuxingNo   string `json:"huxing_no" form:"huxing" validate:"required"`        //户型
	AreaShow   string `json:"area_show" form:"area_show" validate:"required"`     //面积描述
	Area       string `json:"area" form:"area" validate:"required"`               //面积
	Quantity   uint   `json:"quantity" form:"quantity" validate:"required"`       //数量
	Maximum    uint   `json:"maximum" form:"maximum"`                             //最大可选
	Rounds     uint   `json:"rounds" form:"rounds"`                               //轮次
}

type UpdateHuxingBody struct {
	Id         uint   `json:"id" form:"id" validate:"required"`               //id
	BuildingNo string `json:"building_no" form:"building_no" `                //栋号
	HuxingNo   string `json:"huxing_no" form:"huxing_no"`                     //户型
	AreaShow   string `json:"area_show" form:"area_show" validate:"required"` //面积描述
	Area       string `json:"area" form:"area" validate:"required"`           //面积
	Quantity   uint   `json:"quantity" form:"quantity"`                       //数量
	Maximum    uint   `json:"maximum" form:"maximum"`                         //最大可选
	Rounds     uint   `json:"rounds" form:"rounds"`                           //轮次
}

type AccountingFilterBody struct {
	FilterName string `json:"filterName" form:"filterName"`
	Page       uint   `json:"page" form:"page"`
	PageSize   uint   `json:"pageSize" form:"pageSize"`
}

type ContractOnAndCanDeclareBody struct {
	ContractNo string `json:"contract_no" form:"contract_no"`
	StagingId  uint   `json:"staging_id" form:"staging_id"`
}

type ContractUpdateBody struct {
	CanDeclare                        bool   `json:"can_declare" form:"can_declare"`
	ChangeMethod                      string `json:"change_method" form:"change_method" `
	CollectiveLandPropertyCertificate string `json:"collective_land_property_certificate" form:"collective_land_property_certificate"`
	ContractNo                        string `json:"contract_no" form:"contract_no" validate:"required"`
	DateOfDelivery                    string `json:"date_of_delivery" form:"date_of_delivery"  validate:"required"`
	DateOfSigning                     string `json:"date_of_signing" form:"date_of_signing" validate:"required" `
	Desc                              string `json:"desc" form:"desc"`
	HouseWriteOff                     bool   `json:"house_write_off" form:"house_write_off"`
	HouseNumber                       string `json:"house_number" form:"house_number"  validate:"required"`
	OldAddress                        string `json:"old_address" form:"old_address"  validate:"required"`
	PhoneNumber1                      string `json:"phone_number1" form:"phone_number1"`
	PhoneNumber2                      string `json:"phone_number2" form:"phone_number2" `
	Proprietor                        string `json:"proprietor" form:"proprietor"`
	Registration                      string `json:"registration" form:"registration"`
	Signatory                         string `json:"signatory" form:"signatory" validate:"required"`
	SocialCategory                    string `json:"social_category" form:"social_category"`
	StagingID                         uint   `json:"staging_id" form:"staging_id" validate:"required"`
}

type DeclarationEnterResultBody struct {
	DeclarationId uint   `json:"declaration_id" form:"declaration_id" validate:"required"`
	BuildingOn    string `json:"building_on" form:"building_on"`
	RoomOn        string `json:"room_on" form:"room_on"`
	Status        uint   `json:"status"`
}

type GetAllDeclarationFuzzyBody struct {
	StagingId         uint   `json:"stagingId" form:"stagingId" validate:"required"`
	HuxingId          int    `json:"huxing_id" form:"huxing_id" `
	TimeId            int    `json:"time_id" form:"time_id"`
	DeclarationStatus int    `json:"declaration_status" form:"declaration_status" `
	WinningStatus     int    `json:"winning_status" form:"winning_status"`
	ActiveState       int    `json:"active_state" form:"active_state"`
	FilterName        string `json:"filterName" form:"filterName"`
	PaginationQueryBody
}

type CheckFilterBody struct {
	FilterName string `json:"filterName" form:"filterName"`
	Page       uint   `json:"page" form:"page"`
	PageSize   uint   `json:"pageSize" form:"pageSize"`
}
