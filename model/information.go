package model

import (
	"relocate/util/logging"
)

type InformationStatus int

const (
	InformationPending InformationStatus = iota
	InformationPublished
)

func (is InformationStatus) String() string {
	switch is {
	case InformationPending:
		return "待发布中"
	case InformationPublished:
		return "已发布"
	default:
		return "unknown"
	}
}

// 资讯表
type Information struct {
	Model
	Title   string            `json:"title" gorm:"not null;comment:'标题'"`
	Content string            `json:"content" gorm:"not null;type:text(65535);comment:'内容'"`
	Picture string            `json:"picture" gorm:"not null;comment:'图片'"`
	Status  InformationStatus `json:"status" gorm:"not null;comment:'资讯状态'"`
}

func (i Information) TableName() string {
	return "information"
}

func initInformation() {
	if !db.HasTable(&Information{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&Information{}).Error; err != nil {
			panic(err)
		}
	}
}

func (i *Information) Create() error {
	return db.Create(&i).Error
}

func (i *Information) Delete() error {
	return db.Delete(&i).Error
}

func FindInformationByID(id uint) (*Information, error) {
	var information Information
	err := db.Where("ID = ?", id).First(&information).Error
	if err != nil {
		return nil, err
	}
	return &information, nil
}

func (i *Information) UpdateStatus() error {
	sql := db.Save(&i)
	logging.Infof("影响的行数为%d", sql.RowsAffected)
	logging.Infoln(sql.Error)
	return sql.Error
}

func GetInformations(page, pageSize uint) (data *PaginationQ, err error) {
	q := PaginationQ{
		PageSize: pageSize,
		Page:     page,
		Data:     &[]Information{},
	}
	return q.SearchAll(db.Model(&Information{}).Order("created_at desc"))
}

func GetAllPublicInformation(page uint, pageSize uint) (*PaginationQ, error) {
	q := PaginationQ{
		PageSize: pageSize,
		Page:     page,
		Data:     &[]Information{},
	}

	data, err := q.SearchAll(db.Model(&Information{}).Where("status = ?", 1))
	return data, err
}

func UpdateInformationById(info Information) error {
	err := db.Model(&Information{}).Updates(info).Error
	return err
}
