package model

import "relocate/util/logging"

// 管理员表
type Admin struct {
	Model
	AdminName     string `json:"admin_name" gorm:"unique;not null;comment:'管理员账号-唯一'"`
	AdminSignname string `json:"admin_signname" gorm:"not null;comment:'管理员姓名'"`
	Password      string `json:"-" gorm:"not null;comment:'密码'"`
	Salt          string `json:"-" gorm:"not null;comment:'混淆盐'"`
	Visitors      bool   `json:"visitors" gorm:"comment:'1为访客'"`
}

func (a Admin) TableName() string {
	return "admin"
}

func initAdmin() {
	if !db.HasTable(&Admin{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&Admin{}).Error; err != nil {
			panic(err)
		}
		admin := Admin{
			AdminName:     "admin",
			AdminSignname: "管理员",
			Password:      "e7b3e907df88943aff31c0007feecdb6", //123456
			Salt:          "ABCDEF",
		}
		admin.Create()
		admin1 := Admin{
			AdminName:     "admin1",
			AdminSignname: "管理员1",
			Password:      "e7b3e907df88943aff31c0007feecdb6", //123456
			Salt:          "ABCDEF",
		}
		admin1.Create()
		admin2 := Admin{
			AdminName:     "admin2",
			AdminSignname: "管理员2",
			Password:      "e7b3e907df88943aff31c0007feecdb6", //123456
			Salt:          "ABCDEF",
		}
		admin2.Create()
		admin3 := Admin{
			AdminName:     "admin3",
			AdminSignname: "管理员3",
			Password:      "e7b3e907df88943aff31c0007feecdb6", //123456
			Salt:          "ABCDEF",
		}
		admin3.Create()
		admin4 := Admin{
			AdminName:     "admin4",
			AdminSignname: "管理员4",
			Password:      "e7b3e907df88943aff31c0007feecdb6", //123456
			Salt:          "ABCDEF",
		}
		admin4.Create()
		admin5 := Admin{
			AdminName:     "admin5",
			AdminSignname: "管理员5",
			Password:      "e7b3e907df88943aff31c0007feecdb6", //123456
			Salt:          "ABCDEF",
		}
		admin5.Create()
	}
}

func (a *Admin) Create() error {
	return db.Create(&a).Error
}

// 通过管理员账号获取信息
func GetAdminInfo(adminName string) (*Admin, error) {
	var admin Admin
	if err := db.Where("admin_name = ?", adminName).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

// 通过管理员账号获取信息
func GetAllAdminInfo() (interface{}, error) {
	var adminList []Admin
	if err := db.Select("id,admin_name,admin_signname").Find(&adminList).Error; err != nil {
		return nil, err
	}
	return adminList, nil
}
func (a *Admin) UpdateUser() error {
	sql := db.Model(a).Where("admin_name = ?", a.AdminName).Updates(&a)
	rowsAffected := sql.RowsAffected
	logging.Infof("更新影响的记录数%d", rowsAffected)
	logging.Infoln(sql.Error)
	return sql.Error
}
