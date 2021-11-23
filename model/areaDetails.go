package model

// 面积明细表
type AreaDetails struct {
	Model
	ContractNo         string `json:"contract_no" gorm:"not null;comment:'合同号 - 外键'"`
	OperationalDetails string `json:"operational_details" gorm:"not null;comment:'操作表述'"`
	OperationalArea    string `json:"operational_area" gorm:"not null;comment:'操作面积 +300㎡ -300㎡'"`
	RemainingArea      string `json:"remaining_area" gorm:"not null;comment:'剩余面积㎡'"`
}

func (ad AreaDetails) TableName() string {
	return "area_details"
}

func initAreaDetails() {
	if !db.HasTable(&AreaDetails{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&AreaDetails{}).Error; err != nil {
			panic(err)
		}
		//创建外键
		db.Model(&Declaration{}).
			AddForeignKey("contract_no", "contract(contract_no)", "RESTRICT", "RESTRICT")
	}
}

func (ad *AreaDetails) Create() error {
	return db.Create(&ad).Error
}

type ResultAreaDetailsBody struct {
	ContractNo         string `json:"contract_no"`
	OperationalDetails string `json:"operational_details"`
	OperationalArea    string `json:"operational_area"`
	RemainingArea      string `json:"remaining_area"`
}

func GetAreaDetailPageByContextNo(contractNo string, page uint, pageSize uint) (data *PaginationQ, err error) {
	q := PaginationQ{
		Page:     page,
		PageSize: pageSize,
		Data:     &[]ResultAreaDetailsBody{},
	}

	sql := db.Model(AreaDetails{}).Where("contract_no = ?", contractNo)

	data, err = q.SearchAll(sql)
	return
}
