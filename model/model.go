// 模型定义 https://learnku.com/docs/gorm/2018/models/3782
package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"relocate/util/conf"
	"relocate/util/logging"
	"relocate/util/times"
)

type Model struct {
	ID        uint            `json:"id" gorm:"primary_key;comment:'ID'"`
	CreatedAt times.JsonTime  `json:"created_at" gorm:"not null;comment:'创建时间'"`
	UpdatedAt times.JsonTime  `json:"-" gorm:"comment:'更新时间'"`
	DeletedAt *times.JsonTime `json:"-" gorm:"comment:'删除时间'"`
}

var db *gorm.DB

func Setup() {
	var err error
	logging.Info("连接", conf.Data.Database.Type)
	db, err = gorm.Open(conf.Data.Database.Type, conf.Data.Database.Url)
	if err != nil {
		logging.Fatalf("model.Setup err: %v", err)
	}
	db.LogMode(conf.Data.Server.RunMode == gin.DebugMode) //debug模式,显示sql语句
	// 全局禁用表名复数
	db.SingularTable(true)
	initDB()
	//刷新数据库中的表格，使其保持最新
	db.AutoMigrate(&Admin{}, &User{}, &Config{}, &Staging{}, &Contract{}, &UserContract{},
		&Huxing{}, &Time{}, &Declaration{}, &Result{}, &Logging{}, &Information{}, &AreaDetails{},
		&Accounting{}, &Check{})
}

func initDB() {
	initAdmin()
	initUser()
	initConfig()
	initHuXingConfig()
	initStaging()
	initContract()
	initUserContract()
	initHuxing()
	initTime()
	initDeclaration()
	initResult()
	initLogging()
	initInformation()
	initAreaDetails()
	initAreaTransfer()
	initAreaTransferDetail()
}
