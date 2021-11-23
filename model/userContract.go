package model

// 用户合同单匹配表
type UserContract struct {
	Model
	PhoneNumber string `json:"phone_number" gorm:"primary_key;comment:'手机号码-复合主键-外键'"`
	ContractNo  string `json:"contract_no" gorm:"primary_key;comment:'合同号-复合主键-外键'"`
}

func (uc UserContract) TableName() string {
	return "userContract"
}

func initUserContract() {
	if !db.HasTable(&UserContract{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&UserContract{}).Error; err != nil {
			panic(err)
		}
		//创建外键
		db.Model(&UserContract{}).
			AddForeignKey("phone_number", "user(phone_number)", "RESTRICT", "RESTRICT")
		db.Model(&UserContract{}).
			AddForeignKey("contract_no", "contract(contract_no)", "RESTRICT", "RESTRICT")
	}
}

func (uc *UserContract) Create() error {
	return db.Create(&uc).Error
}
