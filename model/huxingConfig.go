package model

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
)

const (
	HuxingGroupingOptionalConfig string = "HuxingGroupingOptionalConfig"
)

// 户型全局配置表
type HuXingConfig struct {
	Key   string `json:"key" gorm:"unique;not null;comment:'key'"`
	Value string `json:"value" gorm:"type:text;not null;comment:'value'"`
}

type HuxingGroupingOptionalConfigJson struct {
	IDs string `json:"ids"`
	Num int    `json:"num"`
}

func (hc HuXingConfig) TableName() string {
	return "huxing_config"
}

func initHuXingConfig() {
	if !db.HasTable(&HuXingConfig{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&HuXingConfig{}).Error; err != nil {
			panic(err)
		}
		var hcs = []HuxingGroupingOptionalConfigJson{
			{IDs: "", Num: 0},
		}
		hcss, _ := json.Marshal(&hcs)
		huxingGroupingOptionalConfig := HuXingConfig{
			Key:   HuxingGroupingOptionalConfig,
			Value: string(hcss),
		}
		huxingGroupingOptionalConfig.Create()
	}
}

func (hc *HuXingConfig) Create() error {
	return db.Create(&hc).Error
}

func SettingHuxingGroupingOptionalConfig(h []HuxingGroupingOptionalConfigJson) error {
	hs, err := json.Marshal(&h)
	if err != nil {
		return err
	}
	if err := db.Model(&HuXingConfig{}).Where("`key` = ?", HuxingGroupingOptionalConfig).
		Update("value", string(hs)).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}

func GetHuxingGroupingOptionalConfig() ([]HuxingGroupingOptionalConfigJson, error) {
	var config HuXingConfig
	if err := db.Where("`key` = ?", HuxingGroupingOptionalConfig).
		First(&config).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	var h []HuxingGroupingOptionalConfigJson
	err := json.Unmarshal([]byte(config.Value), &h)
	return h, err
}
