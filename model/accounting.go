package model

import (
	"encoding/json"
	"relocate/api"
)

//核算
type Accounting struct {
	ID                      uint    `json:"id" gorm:"primary_key;comment:'ID'"`
	ContractNo              string  `json:"contract_no" gorm:"not null;comment:'合同号'"`
	SocialCategory          string  `json:"social_category" gorm:"null;comment:'社别'"`
	Peoples                 string  `json:"peoples" gorm:"not null;comment:'被拆迁人(可能有多人)'"`
	HouseNumber             string  `json:"house_number" gorm:"null;comment:'房屋栋号'"`
	Desc                    string  `json:"desc" gorm:"null;comment:'备注'"`
	InitialHQArea           float64 `json:"initial_hq_area" gorm:"null;comment:'应得补偿安置总面积'"`
	TargetPlacementArea     float64 `json:"target_placement_area" gorm:"null;comment:'指标安置面积'"`
	TemporaryRelocationArea float64 `json:"temporary_relocation_area" gorm:"null;comment:'计算临迁费面积'"`

	ResultList       string       `json:"-" gorm:"type:text"`
	ResultListStruct []ResultList `json:"result_list_struct" gorm:"-"`

	// 计算字段
	PlacementOfNonTargetArea                       float64 `json:"placement_of_non_target_area" gorm:"null;comment:'安置非指标面积'"`                      // 初始回迁面积 - 指标安置面积
	NonIndexAreaRatio                              float64 `json:"non_index_area_ratio" gorm:"null;comment:'非指标面积占比'"`                              // 安置非指标面积 / 应得补偿安置总面积
	IndexAreaRatio                                 float64 `json:"index_area_ratio" gorm:"null;comment:'指标面积占比'"`                                   // 指标安置面积 / 应得补偿安置总面积
	TemporaryRelocationAreaRatioNonIndex           float64 `json:"temporary_relocation_area_ratio_non_index" gorm:"null;comment:'临迁费面积占安置非指标面积比例'"` // 计算临迁费面积 / 安置非指标面积
	RemainingResettlementArea                      float64 `json:"remaining_resettlement_area" gorm:"null;comment:'剩余应得回迁安置面积'"`                    //
	TemporaryRelocationSubPlacementOfNonTargetArea float64 `json:"temporary_relocation_sub_placement_of_non_target_area" gorm:"null;comment:'计算临迁费面积与非指标面积之差'"`
	MeasuredFloorArea                              float64 `json:"measured_floor_area" gorm:"null;comment:'中签房号实测建筑面积'"`               // 当前合同所有中签房间的面积之和
	UseTargetPlacementArea                         float64 `json:"use_target_placement_area" gorm:"null;comment:'已使用的安置指标面积'"`         // 中签房号实测建筑面积 * 指标面积占比
	UsePlacementOfNonTargetArea                    float64 `json:"use_placement_of_non_target_area" gorm:"null;comment:'已使用的安置非指标面积'"` // 中签房号实测建筑面积 * 非指标面积占比
	UseTemporaryRelocationArea                     float64 `json:"use_temporary_relocation_area" gorm:"null;comment:'已使用的计算临迁安置费面积'"`
	RemainingPlacementOfNonTargetArea              float64 `json:"remaining_placement_of_non_target_area" gorm:"null;comment:'剩余安置非指标面积'"`
	RemainingTargetPlacementArea                   float64 `json:"remaining_target_placement_area" gorm:"null;comment:'剩余安置指标面积'"`
	RemainingTemporaryRelocationArea               float64 `json:"remaining_temporary_relocation_area" gorm:"null;comment:'剩余计算临迁费面积'"`
	RemainingInitialHQArea                         float64 `json:"remaining_initial_hq_area" gorm:"null;comment:'剩余应得回迁安置总面积'"`
	AmountOfUsedArea                               float64 `json:"amount_of_used_area" gorm:"null;comment:'购买已使用指安置标面积的金额1000元/㎡'"`
}

type ResultList struct {
	ResultID            uint   `json:"result_id" gorm:"not null;comment:'结果ID'"`
	BuildingNo          string `json:"building_no" gorm:"null;comment:'楼号'"`
	RoomNo              string `json:"room_no" gorm:"null;comment:'房间号'"`
	DeclarationHuxingID uint   `json:"declaration_huxing_id" gorm:"not null;comment:'申报户型ID'"`
	Rounds              uint   `json:"rounds" gorm:"null;comment:'轮数'"`
	DeclarationArea     string `json:"declaration_area" gorm:"not null;comment:'已中签户型㎡'"`
}

func (a Accounting) TableName() string {
	return "accounting"
}

func (a *Accounting) Create() error {
	return db.Create(&a).Error
}

func (a *Accounting) Update() error {
	return db.Model(&Accounting{}).
		Where("contract_no = ?", a.ContractNo).
		Updates(&a).Error
}

func FindAccountingByContractNo(contractNo string) (*Accounting, error) {
	var a Accounting
	err := db.Model(&Accounting{}).Where("contract_no = ?", contractNo).Take(&a).Error
	return &a, err
}

func GetLikeAccounting(accountingFilterBody api.AccountingFilterBody) (data *PaginationQ, err error) {
	q := &PaginationQ{
		PageSize: accountingFilterBody.PageSize,
		Page:     accountingFilterBody.Page,
		Data:     &[]Accounting{},
	}
	args := "%" + accountingFilterBody.FilterName + "%"
	data, err = q.SearchAll(
		db.Model(&Accounting{}).Where("contract_no LIKE ? OR peoples LIKE ?", args, args),
	)
	if err != nil {
		return nil, err
	}
	switch data.Data.(type) {
	case *[]Accounting:
		for i, accounting := range *data.Data.(*[]Accounting) {
			_ = json.Unmarshal([]byte(accounting.ResultList), &((*data.Data.(*[]Accounting))[i].ResultListStruct))
		}
	}
	return
}
