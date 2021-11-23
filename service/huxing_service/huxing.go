package huxing_service

import (
	"github.com/jinzhu/gorm"
	"relocate/model"
	"relocate/util/errors"
)

//查找户型
func Gethuxing(id uint) (*model.Huxing, error) {
	var t model.Huxing
	t.ID = id
	err := t.FindHuxingByID()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.BadError("户型不存在")
		}
		return nil, err
	}
	return &t, nil
}

func AddHuxing(stagingId uint, buildingNo, huxingNo, area, areaShow string, quantity, maximum, rounds uint) error {
	if stagingId != 0 {
		// 判断是否存在
		if !model.ExistsStagingById(stagingId) {
			return errors.New("不存在此分期")
		}
	}

	huxing := model.Huxing{
		StagingID:  stagingId,
		BuildingNo: buildingNo,
		HuxingNo:   huxingNo,
		Area:       area,
		AreaShow:   areaShow,
		Quantity:   quantity,
		Maximum:    maximum,
		Rounds:     rounds,
	}

	return model.AddHuxing(huxing)

}

func UpdateHuxing(id uint, buildingNo, huxingNo, area, areaShow string, quantity, maximum, rounds uint) error {
	if !model.ExistsHuxingById(id) {
		return errors.New("该户型不存在")
	}

	huxing := model.Huxing{
		BuildingNo: buildingNo,
		HuxingNo:   huxingNo,
		Area:       area,
		AreaShow:   areaShow,
		Quantity:   quantity,
		Maximum:    maximum,
		Rounds:     rounds,
	}
	huxing.ID = id

	return model.UpdateHuxing(huxing)
}
