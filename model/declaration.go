package model

import (
	"github.com/jinzhu/gorm"
	"relocate/api"
	"relocate/util/errors"
	"relocate/util/logging"
)

type DeclarationStatus int

const (
	DeclarationOngoing DeclarationStatus = iota
	DeclarationConfirmed
)

func (d DeclarationStatus) String() string {
	switch d {
	case DeclarationOngoing:
		return "进行中"
	case DeclarationConfirmed:
		return "已确认"
	default:
		return "unknown"
	}
}

type WinningStatus int

const (
	WinningNo WinningStatus = iota
	WinningYes
)

func (w WinningStatus) String() string {
	switch w {
	case WinningNo:
		return "不中"
	case WinningYes:
		return "中签"
	default:
		return "unknown"
	}
}

//type HuxingData struct {
//	HuxingNo string `json:"huxing_no"`
//	Area     string `json:"area"`
//	Quantity uint   `json:"quantity"`
//}
//
//type DeclarationHuxingData struct {
//	HuxingData []HuxingData `json:"huxing_data"`
//}

// 申报表
type Declaration struct {
	Model
	//DeclarationHuxingData string            `json:"declaration_huxing_data" gorm:"not null;comment:'申报户型数据-存储json格式文本'"`
	TimeID                uint              `json:"time_id" gorm:"not null;comment:'时段id-外键'"`
	TimeName              string            `json:"time_name" gorm:"not null;comment:'冗余 时段表述'"`
	StagingID             uint              `json:"staging_id" gorm:"not null;comment:'分期数ID-外键'"`
	Rounds                uint              `json:"rounds" gorm:"null;comment:'轮数'"`
	ContractNo            string            `json:"contract_no" gorm:"not null;comment:'合同号-外键'"`
	DeclarationStatus     DeclarationStatus `json:"declaration_status" gorm:"not null;comment:'申报状态'"`
	ActiveState           bool              `json:"active_state" gorm:"not null;comment:'有效状态-是否作废'"`
	WinningStatus         WinningStatus     `json:"winning_status" gorm:"null;comment:'中签状态'"`
	DeclarationHuxingID   uint              `json:"declaration_huxing_id" gorm:"not null;comment:'申报户型ID'"`
	DeclarationHuxingNo   string            `json:"declaration_huxing_no" gorm:"not null;comment:'申报户型'"`
	DeclarationBuildingNo string            `json:"declaration_huxing_building_no" gorm:"not null;comment:'申报户型栋号'"`
	DeclarationAreaShow   string            `json:"declaration_area_show" gorm:"not null;comment:'申报户型面积显示'"`
	DeclarationArea       string            `json:"declaration_area" gorm:"not null;comment:'申报面积㎡'"`
	Trustee               string            `json:"trustee" gorm:"null;comment:'受托人'"`
	TrusteeCardNumber     string            `json:"trustee_card_number" gorm:"null;comment:'受托人身份证号码'"`
	TrusteePhoneNumber    string            `json:"trustee_phone_number" gorm:"null;comment:'受托人手机号码'"`
	TrusteeRelationship   string            `json:"trustee_relationship" gorm:"null;comment:'受托人关系'"`
	Operator              string            `json:"operator" gorm:"not null;comment:'操作人 管理员姓名、登录人姓名'"`
	Printer               string            `json:"printer" gorm:"null;comment:'申报表打印管理员姓名'"`
}

func (d Declaration) TableName() string {
	return "declaration"
}

func initDeclaration() {
	if !db.HasTable(&Declaration{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&Declaration{}).Error; err != nil {
			panic(err)
		}
		//创建外键
		db.Model(&Declaration{}).
			AddForeignKey("time_id", "time(id)", "RESTRICT", "RESTRICT")
		db.Model(&Declaration{}).
			AddForeignKey("staging_id", "staging(id)", "RESTRICT", "RESTRICT")
		db.Model(&Declaration{}).
			AddForeignKey("contract_no", "contract(contract_no)", "RESTRICT", "RESTRICT")
		db.Model(&Declaration{}).
			AddForeignKey("declaration_huxing_id", "huxing(id)", "RESTRICT", "RESTRICT")
	}
}

func (d *Declaration) Create() error {
	return db.Create(&d).Error
}

type DeclarationData struct {
	ID                    int               `json:"id"`
	TimeID                uint              `json:"time_id"`
	TimeName              string            `json:"time_name"`
	ContractNo            string            `json:"contract_no"`
	StagingID             uint              `json:"staging_id"`
	StagingName           string            `json:"staging_name"`
	Rounds                uint              `json:"rounds"`
	DeclarationArea       string            `json:"declaration_area"`
	DeclarationHuxingID   uint              `json:"declaration_huxing_id"`
	DeclarationHuxingNo   string            `json:"declaration_huxing_no"`
	DeclarationBuildingNo string            `json:"declaration_huxing_building_no"`
	DeclarationAreaShow   string            `json:"declaration_area_show"`
	ActiveState           bool              `json:"active_state"`
	WinningStatus         WinningStatus     `json:"winning_status"`
	DeclarationStatus     DeclarationStatus `json:"declaration_status"`
	Peoples               string            `json:"peoples"`
	OldAddress            string            `json:"old_address"`
	CardNumber            string            `json:"card_number"`
	HouseNumber           string            `json:"house_number"`
	PhoneNumber1          string            `json:"phone_number1"`
	PhoneNumber2          string            `json:"phone_number2"`
	Trustee               string            `json:"trustee"`
	TrusteeCardNumber     string            `json:"trustee_card_number"`
	TrusteePhoneNumber    string            `json:"trustee_phone_number"`
	TrusteeRelationship   string            `json:"trustee_relationship"`
	SocialCategory        string            `json:"social_category"`
}

func GetAllDeclaration() (data *PaginationQ, err error) {
	q := PaginationQ{
		Data: &[]DeclarationData{},
	}
	return q.SearchAll(
		db.Table("declaration").Select("declaration.*," +
			"contract.*," +
			"staging.*",
		).Joins(
			"join contract on contract.contract_no = declaration.contract_no",
		).Joins(
			"join staging on staging.id = declaration.staging_id",
		).Order("declaration.created_at desc"),
	)
}

type Trustee struct {
	Trustee             string `json:"trustee"`
	TrusteeCardNumber   string `json:"trustee_card_number"`
	TrusteePhoneNumber  string `json:"trustee_phone_number"`
	TrusteeRelationship string `json:"trustee_relationship"`
}

type Detail struct {
	ID                    int               `json:"id"`
	Operator              string            `json:"operator"`
	ContractNo            string            `json:"contract_no"`
	Peoples               string            `json:"peoples"`
	CardNumber            string            `json:"card_number"`
	OldAddress            string            `json:"old_address"`
	InitialHQArea         string            `json:"initial_hq_area"`
	RemainingHQArea       string            `json:"remaining_hq_area"`
	DeclarationArea       string            `json:"declaration_area"`
	DeclarationHuxingID   uint              `json:"declaration_huxing_id"`
	DeclarationHuxingNo   string            `json:"declaration_huxing_no"`
	DeclarationBuildingNo string            `json:"declaration_huxing_building_no"`
	DeclarationAreaShow   string            `json:"declaration_area_show"`
	DeclarationStatus     DeclarationStatus `json:"declaration_status"`
	WinningStatus         WinningStatus     `json:"winning_status"`
	ActiveState           bool              `json:"active_state"`
	TimeID                uint              `json:"time_id"`
	TimeName              string            `json:"time_name"`
	SocialCategory        string            `json:"social_category"`
	HouseNumber           string            `json:"house_number"`
	Trustee
}

type ContractInfo struct {
	ContractNo      string `json:"contract_no"`
	Peoples         string `json:"peoples"`
	CardNumber      string `json:"card_number"`
	OldAddress      string `json:"old_address"`
	InitialHQArea   string `json:"initial_hq_area"`
	RemainingHQArea string `json:"remaining_hq_area"`
}

type DeclarationItem struct {
	ID                    int               `json:"id"`
	Operator              string            `json:"operator"`
	Peoples               string            `json:"peoples"`
	DeclarationArea       string            `json:"declaration_area"`
	DeclarationHuxingID   uint              `json:"declaration_huxing_id"`
	DeclarationHuxingNo   string            `json:"declaration_huxing_no"`
	DeclarationBuildingNo string            `json:"declaration_huxing_building_no"`
	DeclarationAreaShow   string            `json:"declaration_area_show"`
	StagingID             uint              `json:"staging_id"`
	StagingName           string            `json:"staging_name"`
	DeclarationStatus     DeclarationStatus `json:"declaration_status"`
	WinningStatus         WinningStatus     `json:"winning_status"`
	ActiveState           bool              `json:"active_state"`
}

func GetActiveStateDeclaration(contractNo string, stagingID, rounds uint) (data []Declaration, err error) {
	if err = db.Model(&Declaration{}).Where("active_state = ? "+
		"and contract_no = ? "+
		"and staging_id = ? "+
		"and rounds = ?",
		true,
		contractNo,
		stagingID, rounds).Find(&data).Error; err != nil {
		return
	}
	return
}

func FindDeclarationByID(declarationID uint) (declaration *Declaration, err error) {
	declaration = new(Declaration)
	err = db.Model(&Declaration{}).Where("id = ?", declarationID).First(&declaration).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.BadError("该申报不存在")
		}
		return nil, err
	}
	return
}

func GetDeclarationListByID(id []uint) ([]*Declaration, error) {
	var data []*Declaration
	if err := db.Model(&Declaration{}).Where("id in (?)", id).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func GetDeclarationByID(id uint) (Declaration, error) {
	var data Declaration
	err := db.Model(&Declaration{}).Where("id = ?", id).First(&data).Error
	return data, err
}

func (d *Declaration) Update() error {
	sql := db.Model(&Declaration{}).Save(&d)
	if err := sql.Error; err != nil {
		return err
	}
	rowsAffected := sql.RowsAffected
	logging.Infof("更新影响的记录数%d", rowsAffected)
	logging.Infoln(sql.Error)
	return nil
}

func (d *Declaration) UpdateDeclarationStatus(id uint, status DeclarationStatus) error {
	sql := db.Model(&Declaration{}).Where("id = ?", id).Update("declaration_status = ?", status)
	if err := sql.Error; err != nil {
		return err
	}
	rowsAffected := sql.RowsAffected
	logging.Infof("更新影响的记录数%d", rowsAffected)
	logging.Infoln(sql.Error)
	return nil
}

func (d *Declaration) UpdateWinningStatus(id uint, status WinningStatus) error {
	sql := db.Model(&Declaration{}).Where("id = ?", id).Update("winning_status = ?", status)
	if err := sql.Error; err != nil {
		return err
	}
	rowsAffected := sql.RowsAffected
	logging.Infof("更新影响的记录数%d", rowsAffected)
	logging.Infoln(sql.Error)
	return nil
}

func SetAllDeclarationStatusBuContractOn(contracts []api.ContractOnAndCanDeclareBody, status uint) error {

	sql := db.Debug().Model(Declaration{}).Where(contracts).Update("declaration_status", status)

	if err := sql.Error; err != nil {
		return err
	}
	rowsAffected := sql.RowsAffected
	logging.Infof("更新影响的记录数%d", rowsAffected)
	logging.Infoln(sql.Error)
	return nil

}

func DeclarationDeleteAble(d Declaration) error {
	// 申报状态
	if d.DeclarationStatus == 1 {
		return errors.New("该申报已经确定，请先修改申报状态再进行修改")
	}

	// 是否是待抽签状态
	if d.WinningStatus == 0 || d.WinningStatus == 1 {
		return errors.New("该申报已经是待抽签状态，不可删除")
	}

	if d.ActiveState == true {
		return errors.New("该申报有效，不可直接删除")
	}
	return nil
}

func DeleteDeclarationById(id uint) error {
	return db.Where("id = ?", id).Delete(&Declaration{}).Error
}

func GetAllDeclarationByStagingId(id uint) ([]Declaration, error) {
	var declarations []Declaration

	err := db.Where("staging_id = ? and active_state = 1", id).Find(&declarations).Error

	return declarations, err
}

func ExistsStagingInDeclaration(id uint) bool {
	var de Declaration
	db.Where("staging_id = ?", id).First(&de)
	return de.ID != 0
}

func GetAllDeclarationFuzzy(body api.GetAllDeclarationFuzzyBody) (*PaginationQ, error) {
	q := &PaginationQ{
		PageSize: body.PageSize,
		Page:     body.Page,
		Data:     &[]Declaration{},
	}

	sql := db.Model(&Declaration{})
	if body.HuxingId != -1 {
		sql = sql.Where("declaration_huxing_id = ?", body.HuxingId)
	}
	if body.TimeId != -1 {
		sql = sql.Where("time_id = ?", body.TimeId)
	}
	if body.DeclarationStatus != -1 {
		if body.DeclarationStatus != 0 && body.DeclarationStatus != 1 {
			return nil, errors.New("申报状态错误 0->进行中 1->申报成功")
		} else {
			sql = sql.Where("declaration_status = ?", body.DeclarationStatus)
		}
	}
	if body.WinningStatus != -1 {
		if body.WinningStatus != 0 && body.WinningStatus != 1 {
			return nil, errors.New("抽签错误 0->抽签失败 1->中签")
		} else {
			sql = sql.Where("winning_status = ?", body.WinningStatus)
		}
	}
	if body.ActiveState != -1 {
		if body.ActiveState != 0 && body.ActiveState != 1 {
			return nil, errors.New("申报有效转状态 0->无效 1->有效")
		} else {
			sql = sql.Where("active_state = ?", body.ActiveState)
		}
	}

	if body.FilterName != "" {
		// 身份证号、合同号、手机号
		body.FilterName = "%" + body.FilterName + "%"
		sql = sql.Where("trustee_card_number Like ? or contract_no Like ? or trustee_phone_number Like ?",
			body.FilterName, body.FilterName, body.FilterName)
	}

	data, err := q.SearchAll(sql)
	return data, err
}

func GetDeclarationByContractNo(contractOn string) (declarations []Declaration, err error) {
	err = db.Model(Declaration{}).Where("contract_no = ?", contractOn).Find(&declarations).Error
	return
}

func GetDeclarationByContractNoAndStagingId(contractNo string, stagingId uint) (allList []Declaration, currentList []Declaration, err error) {
	// 返回两组数据
	err = db.Model(&Declaration{}).Where("contract_no = ? and staging_id = ?", contractNo, stagingId).Find(&allList).Error
	if err != nil {
		return nil, nil, err
	}
	err = db.Model(&Declaration{}).Where("contract_no = ? and staging_id = ? and active_state = 1", contractNo, stagingId).Find(&currentList).Error
	if err != nil {
		return nil, nil, err
	}

	return
}

func ExistsDeclarationById(id uint) bool {
	var declaration Declaration
	db.Model(Declaration{}).Where("id = ?", id).First(&declaration)

	return declaration.ID != 0
}

func UpdateDeclarationByTrusteeByContractNo(declaration api.UpdateDeclaration) error {
	sql := db.Model(&Declaration{})

	if declaration.DeclarationID != 0 {
		sql = sql.Where("id = ?", declaration.DeclarationID)
	}

	if declaration.HuxingID != 0 {
		sql = sql.Where("declaration_huxing_id = ?", declaration.HuxingID)
	}
	var body Declaration
	body.Trustee = declaration.Trustee
	body.TrusteePhoneNumber = declaration.TrusteePhoneNumber
	body.TrusteeCardNumber = declaration.TrusteeCardNumber
	body.TrusteeRelationship = declaration.TrusteeRelationship
	err := sql.Updates(body).Error
	return err
}

func IsDeclarationBelongToClientByPhone(id uint, phone string) (bool, error) {
	var de Declaration
	err := db.Model(&Declaration{}).Where("id = ?", id).First(&de).Error
	if err != nil {
		return false, err
	}

	data, err := FindAllContractByPhone(phone)
	if err != nil {
		return false, err
	}

	for _, contract := range data {
		if contract.ContractNo == de.ContractNo {
			return true, nil
		}
	}

	return false, nil
}
