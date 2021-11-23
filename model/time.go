package model

import (
	"relocate/util/logging"
)

// 时段表
type Time struct {
	Model
	Name        string `json:"name" gorm:"not null;comment:'时段名称'"`
	OptionalNum uint   `json:"optional_num" gorm:"not null;comment:'可选数'"`
	SelectedNum uint   `json:"selected_num" gorm:"not null;comment:'已选数'"`
	StagingId   uint   `json:"staging_id" gorm:"not null;comment:'分期id'"`
}

func (t Time) TableName() string {
	return "time"
}

func initTime() {
	if !db.HasTable(&Time{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&Time{}).Error; err != nil {
			panic(err)
		}
		db.Model(&Time{}).
			AddForeignKey("staging", "staging(id)", "RESTRICT", "RESTRICT")
	}
}

func (t *Time) FindTimeByID() (err error) {

	err = db.Where("ID = ?", t.ID).First(&t).Error
	return
}

func (t *Time) Update() error {
	sql := db.Model(&t).Updates(Time{
		Name:        t.Name,
		OptionalNum: t.OptionalNum,
		SelectedNum: t.SelectedNum,
	})
	logging.Infof("影响的行数为%d", sql.RowsAffected)
	logging.Infoln(sql.Error)
	return sql.Error
}

// QueryTimeByStagingID 根据分期id 获取现场确认时间
func QueryTimeByStagingID(stagingID uint) ([]Time, error) {
	var times []Time
	err := db.Where("staging_id = ?", stagingID).Find(&times).Error

	if err != nil {
		return nil, err
	}

	return times, err
}

// SaveTime 根据传入body 添加时间
func SaveTime(name string, StagingId uint, optionNum uint) error {

	err := db.Create(&Time{
		Name:        name,
		StagingId:   StagingId,
		OptionalNum: optionNum,
		SelectedNum: 0,
	}).Error

	return err
}

// ExistsTime 根据id查看时间段是否存在
func ExistsTime(id uint) bool {
	var time Time
	db.Select("id").Where("id = ?", id).First(&time)
	return time.ID > 0
}

// UpdateTime 根据参数id,name,optionalNum 跟新时间段
func UpdateTime(id uint, name string, optional uint) error {
	var time Time
	time.ID = id
	if name != "" {
		time.Name = name
	}
	if optional != 0 {
		time.OptionalNum = optional
	}
	err := db.Model(Time{}).Where("id = ?", time.ID).Updates(&time).Error

	return err
}

func DeleteTime(id uint) error {
	err := db.Where("id = ?", id).Delete(&Time{}).Error

	return err
}
