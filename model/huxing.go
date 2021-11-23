package model

// 户型
type Huxing struct {
	Model
	StagingID  uint   `json:"staging_id" gorm:"not null;comment:'分期数ID'"`
	BuildingNo string `json:"building_no" gorm:"not null;comment:'栋号'"`
	HuxingNo   string `json:"huxing_no" gorm:"not null;comment:'户型'"`
	Area       string `json:"area" gorm:"not null;comment:'面积㎡'"`
	AreaShow   string `json:"area_show" gorm:"not null;comment:'面积显示'"`
	Quantity   uint   `json:"quantity" gorm:"not null;comment:'套数'"`
	Maximum    uint   `json:"maximum" gorm:"not null;comment:'最多可选'"`
	Rounds     uint   `json:"rounds" gorm:"null;comment:'轮数'"`
}

func (h Huxing) TableName() string {
	return "huxing"
}

func initHuxing() {
	if !db.HasTable(&Huxing{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&Huxing{}).Error; err != nil {
			panic(err)
		}
		db.Model(&Declaration{}).
			AddForeignKey("staging_id", "staging(id)", "RESTRICT", "RESTRICT")
	}
}

type FindHuxingAndSelected struct {
	Id         uint   `json:"id"`
	StagingID  uint   `json:"staging_id"`
	BuildingNo string `json:"building_no"`
	HuxingNo   string `json:"huxing_no"`
	AreaShow   string `json:"area_show"`
	Area       string `json:"area"`
	Quantity   uint   `json:"quantity"`
	Maximum    uint   `json:"maximum"`
	Selected   uint   `json:"selected"`
	Rounds     uint   `json:"rounds"`
}

func GetAllOptionalHuxing() ([]Huxing, error) {
	var huxing []Huxing
	err := db.Where("maximum > 0").Find(&huxing).Error
	if err != nil {
		return nil, err
	}
	return huxing, nil
}

func FindAllOptionalHuxing(stagingID uint) (data interface{}, err error) {
	var date []FindHuxingAndSelected
	if err := db.Table("huxing h").
		Select("h.*,(select sum(1) from declaration d where d.declaration_huxing_id = h.id AND active_state != false) as selected").
		Where("h.deleted_at is null and h.maximum > 0 and staging_id = ?", stagingID).Having("selected < maximum").Find(&date).Error; err != nil {
		return nil, err
	}

	return date, nil
}

func FindAllHuxing(stagingID uint) (data interface{}, err error) {
	var date []FindHuxingAndSelected
	if err := db.Table("huxing h").
		Select("h.*,(select sum(1) from declaration d where d.declaration_huxing_id = h.id AND active_state != false) as selected").
		Where("h.deleted_at is null and h.maximum > 0 and staging_id = ?", stagingID).Find(&date).Error; err != nil {
		return nil, err
	}

	return date, nil
}

func (h *Huxing) FindHuxingByID() (err error) {
	err = db.Where("ID = ?", h.ID).First(&h).Error
	return
}

func DelteHuxingById(id uint) error {
	err := db.Model(&Huxing{}).Where("id = ?", id).Delete(Huxing{}).Error
	return err
}

func AddHuxing(huxing Huxing) error {
	return db.Create(&huxing).Error
}

func ExistsHuxingById(id uint) bool {
	var h Huxing
	db.Where("ID = ?", id).First(&h)
	return h.ID > 0
}

func UpdateHuxing(huxing Huxing) error {
	return db.Debug().Model(&Huxing{}).Updates(&huxing).Error
}

func GetAllOptionHuxing() (total uint, err error) {
	err = db.Model(Huxing{}).Select("Count(maximum) as total").First(&total).Error
	return
}
