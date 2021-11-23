package model

// 分期，如：西华一期
type Staging struct {
	Model
	StagingName string `json:"staging_name" gorm:"not null;comment:'分期名称'"`
}

func (s Staging) TableName() string {
	return "staging"
}

func initStaging() {
	if !db.HasTable(&Staging{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&Staging{}).Error; err != nil {
			panic(err)
		}
	}
}

func FindStagingById(id uint) (*Staging, error) {
	var staging Staging
	if err := db.Where("id = ?", id).First(&staging).Error; err != nil {
		return nil, err
	}
	return &staging, nil
}

type ContractCount struct {
	Id          uint   `json:"id"`
	StagingName string `json:"staging_name"`
	Count       int    `json:"count"`
}

func GetStagingContractCount(page, pageSize uint) (data *PaginationQ, err error) {
	q := PaginationQ{
		PageSize: pageSize,
		Page:     page,
		Data:     &[]ContractCount{},
	}
	return q.SearchAll(db.Table("staging s").Select("s.*,(select sum(1) from contract c where c.staging_id = s.id) as count"))
}

func ExistsStagingById(stagingId uint) bool {
	var stag Staging
	db.Select("id").Where("id = ?", stagingId).First(&stag)

	return stag.ID > 0
}

type GetStagingPageBody struct {
	Id          uint   `json:"id"`
	StagingName string `json:"staging_name"`
	CreatedAt   string `json:"created_at"`
}

func GetStagingPage(page uint, pageSize uint) (q *PaginationQ, err error) {
	q = &PaginationQ{
		PageSize: pageSize,
		Page:     page,
		Data:     &[]GetStagingPageBody{},
	}
	return q.SearchAll(db.Table("staging").Select("*"))
}

func SaveStaging(stagingName string) error {
	var staging Staging

	staging.StagingName = stagingName

	err := db.Create(&staging).Error

	return err
}
func (S *Staging) Update() error {
	return db.Model(Staging{}).Where("id = ?", S.ID).Updates(S).Error
}
func (s *Staging) Get() error {
	return db.Where("id = ?", s.ID).First(&s).Error
}
