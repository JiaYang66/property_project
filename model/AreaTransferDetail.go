/**
* @Author: SAmwake
* @Date: 2021/11/12 16:51
 */
package model

// 面积转移细节表
// 用户之间通过合同转移面积成功之后，创建面积转移明细表
type AreaTransferDetail struct {
	Model
	PhoneNumber       string  `json:"phone_number" gorm:"not null;comment:'手机号码'"`
	Name              string  `json:"name" gorm:"not null;comment:'真实姓名'"`
	ContractNo        string  `json:"contract_no" gorm:"not null;comment:'合同号'"`
	IdNumber          string  `json:"id_number" gorm:"not null;comment:'身份证号码'"`
	Area              float64 `json:"area" gorm:"not null;comment:'转移面积'"`
	PreTransferArea   float64 `json:"pre_transfer_area" gorm:"not null;comment:'转移前剩余面积'"`
	AfterTransferArea float64 `json:"after_transfer_area" gorm:"not null;comment:'转移后剩余面积'"`

	TargetIdNumber          string  `json:"target_id_number" gorm:"not null;comment:'目标身份证号码'"`
	TargetContractNo        string  `json:"target_contract_no" gorm:"not null;comment:'目标合同号'"`
	TargetPhoneNumber       string  `json:"target_phone_number" gorm:"not null;comment:'目标手机号'"`
	TargetName              string  `json:"target_name" gorm:"not null;comment:'目标真实姓名'"`
	TargetPreTransferArea   float64 `json:"target_pre_transfer_area" gorm:"not null;comment:'目标转移前剩余面积'"`
	TargetAfterTransferArea float64 `json:"target_after_transfer_area" gorm:"not null;comment:'目标转移后剩余面积'"`

	Operator       string  `json:"operator" gorm:"not null;comment:'操作员'"`
	AreaUnitPrice  float64 `json:"area_unit_price" gorm:"not null;comment:'面积单价'"`
	AreaPriceTotal float64 `json:"area_price_total" gorm:"not null;comment:'面积总价'"`
	DateOfSigning  string  `json:"date_of_signing" gorm:"not null;comment:'签署协议日期'"`
}

func (c AreaTransferDetail) TableName() string {
	return "area_transfer_detail"
}

func initAreaTransferDetail() {
	if !db.HasTable(&AreaTransferDetail{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&AreaTransferDetail{}).Error; err != nil {
			panic(err)
		}
	}
}

func (area *AreaTransferDetail) Create() error {
	return db.Create(&area).Error
}
