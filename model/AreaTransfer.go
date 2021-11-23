/**
* @Author: SAmwake
* @Date: 2021/11/12 16:28
 */
package model

type StatusType int

const (
	UnderWag StatusType = iota
	Finish
)

func (it StatusType) String() string {
	switch it {
	case UnderWag:
		return "进行中"
	case Finish:
		return "已完成"
	default:
		return "unknown"
	}
}

// 面积转移表
// 当用户与用户之间需要转移面积的时候需要用到的表
type AreaTransfer struct {
	Model
	PhoneNumber       string     `json:"phone_number" gorm:"not null;comment:'手机号码'"`
	Name              string     `json:"name" gorm:"not null;comment:'真实姓名'"`
	ContractNo        string     `json:"contract_no" gorm:"not null;comment:'合同号'"`
	Status            StatusType `json:"status" gorm:"not null;comment:'转移状态'"`
	TargetContractNo  string     `json:"target_contract_no" gorm:"not null;comment:'目标合同号'"`
	TargetPhoneNumber string     `json:"target_phone_number" gorm:"not null;comment:'目标手机号'"`
	TargetName        string     `json:"target_name" gorm:"not null;comment:'目标真实姓名'"`
	Area              float64    `json:"area" gorm:"not null;comment:'转移面积'"`
	DateOfSigning     string     `json:"date_of_signing" gorm:"not null;comment:'签署协议日期'"`
	AreaUnitPrice     float64    `json:"area_unit_price" gorm:"not null;comment:'面积单价'"`
}

func (c AreaTransfer) TableName() string {
	return "area_transfer"
}

func initAreaTransfer() {
	if !db.HasTable(&AreaTransfer{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&AreaTransfer{}).Error; err != nil {
			panic(err)
		}
	}
}

func IsAreaContractNotProcessed(contractNo string) bool {
	var areaContract AreaTransfer

	db.Model(AreaTransfer{}).Where("contract_no = ? and status = 0", contractNo).First(&areaContract)

	return areaContract.ID != 0
}

func (area *AreaTransfer) Create() error {
	return db.Create(&area).Error
}

func (area *AreaTransfer) Update() error {
	return db.Model(&AreaTransfer{}).Where("id = ?", area.ID).Updates(&area).Error
}

func (area *AreaTransfer) Get() error {
	return db.Where("id = ?", area.ID).First(&area).Error
}

func IsHasDeclarationNotProcessedByContractNo(contractNo string) (bool, error) {
	var declaration Declaration
	err := db.Model(Declaration{}).Where("contract_no = ? and declaration_status = 0", contractNo).First(&declaration).Error
	return declaration.ID != 0, err
}
