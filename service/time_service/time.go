package time_service

import (
	"github.com/jinzhu/gorm"
	"relocate/model"
	"relocate/util/errors"
)

//查找时段
func getTime(id uint) (*model.Time, error) {
	var t model.Time
	t.ID = id
	err := t.FindTimeByID()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.BadError("时段不存在")
		}
		return nil, err
	}
	return &t, nil
}

func GetTimeByStagingId(id uint) ([]model.Time, error) {
	times, err := model.QueryTimeByStagingID(id)

	if err != nil {
		return nil, err
	}

	return times, nil
}

func SaveUser(name string, stagingId uint, optionalNum uint) error {
	if model.ExistsStagingById(stagingId) {
		if err := model.SaveTime(name, stagingId, optionalNum); err == nil {
			return nil
		} else {
			return errors.New("添加失败")
		}

	} else {
		return errors.New("不存在此期数")
	}
}

func UpdateTime(id uint, name string, option uint) error {
	if id == 0 {
		return errors.New("id 不能为空")
	}

	//查看是否存在此id
	if model.ExistsTime(id) {
		err := model.UpdateTime(id, name, option)
		return err
	}
	return errors.New("不存在此id")
}

func DeleteById(id uint) error {
	if model.ExistsTime(id) {
		err := model.DeleteTime(id)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("不存在此id")
}
