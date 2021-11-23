package admin_service

import (
	"github.com/jinzhu/gorm"
	"relocate/model"
	"relocate/util/errors"
	"relocate/util/sign"
)

func GenerateToken(adminName, password string) (interface{}, error) {
	admin, err := model.GetAdminInfo(adminName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errors.BadError("管理员账号不存在")
		}
		return "", err
	}
	//校验密码
	if admin.Password != sign.EncodeMD5(password+admin.Salt) {
		return "", errors.BadError("密码错误")
	}
	//生成jwt-token
	token, err := sign.GenerateToken(string(admin.ID), adminName, sign.AdminClaimsType)
	if err != nil {
		return "", err
	}

	return map[string]interface{}{"token": token, "signName": admin.AdminSignname}, nil
}

func UpdateAdminPassword(adminName, password, newPassword string) error {
	admin, err := model.GetAdminInfo(adminName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.BadError("管理员账号不存在")
		}
		return err
	}
	//校验密码
	if admin.Password != sign.EncodeMD5(password+admin.Salt) {
		return errors.BadError("密码错误")
	}
	a := &model.Admin{
		AdminName: adminName,
		Password:  sign.EncodeMD5(newPassword + admin.Salt),
	}
	return a.UpdateUser()
}
