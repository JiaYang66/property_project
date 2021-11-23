package contract_service

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io"
	"relocate/api"
	"relocate/model"
	"relocate/service/staging_service"
	"relocate/util/convert"
	"relocate/util/errors"
	"relocate/util/excel"
	"strconv"
	"strings"
)

func UpdateHouseWriteOffStatusList(ContractNoList []string, HouseWriteOff bool, operator, data string) (err error) {
	err = model.UpdateHouseWriteOffStatusList(ContractNoList, HouseWriteOff)
	if err == nil {
		id, _ := model.GetNowStagingConfig()
		staging, _ := staging_service.GetStagingInfoById(id)
		logging := model.Logging{
			Username:    operator,
			StagingName: staging.StagingName,
			Operation:   "批量修改合同号的注销状态",
			Details:     data,
		}
		logging.Create()
	}
	return err
}

func GetStagingById(stagingID uint) bool {
	_, err := staging_service.GetStagingInfoById(stagingID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
	}
	return true
}

func AddContract(body api.AddContractBody, operator, data string) error {
	cardNumber := strings.Replace(body.CardNumber, "，", ",", -1)
	cardNumber = strings.Replace(cardNumber, "（", "(", -1)
	cardNumber = strings.Replace(cardNumber, "）", ")", -1)
	if body.Registration == "" {
		body.Registration = "无证"
	}
	if body.Proprietor == "" {
		body.Proprietor = "无证"
	}
	if body.CollectiveLandPropertyCertificate == "" {
		body.CollectiveLandPropertyCertificate = "无证"
	}
	contract := &model.Contract{
		ContractNo:                        body.ContractNo,
		SocialCategory:                    body.SocialCategory,
		Peoples:                           body.Peoples,
		CardNumber:                        cardNumber,
		HouseNumber:                       body.HouseNumber,
		OldAddress:                        body.OldAddress,
		PhoneNumber1:                      body.PhoneNumber1,
		PhoneNumber2:                      body.PhoneNumber2,
		DateOfSigning:                     body.DateOfSigning,
		DateOfDelivery:                    body.DateOfDelivery,
		Signatory:                         body.Signatory,
		Registration:                      body.Registration,
		InitialHQArea:                     body.InitialHQArea,
		RemainingHQArea:                   body.InitialHQArea,
		IsDelivery:                        body.IsDelivery,
		CanDeclare:                        body.CanDeclare,
		StagingID:                         body.StagingID,
		HouseWriteOff:                     body.HouseWriteOff,
		Proprietor:                        body.Proprietor,
		Desc:                              body.Desc,
		CollectiveLandPropertyCertificate: body.CollectiveLandPropertyCertificate,
		ChangeMethod:                      body.ChangeMethod,
	}
	contract.Single = staging_service.JudgmentSingle(cardNumber)
	if err := contract.Create(); err != nil {
		return err
	}
	id, _ := model.GetNowStagingConfig()
	staging, _ := staging_service.GetStagingInfoById(id)
	logging := model.Logging{
		Username:    operator,
		StagingName: staging.StagingName,
		Operation:   "添加合同号：" + body.ContractNo,
		Details:     data,
	}
	logging.Create()
	return nil
}

func UpdateCardNumber(contractNo, cardNumber, operator, data string) error {
	if err := model.UpdateCardNumber(contractNo, cardNumber); err != nil {
		return err
	}
	id, _ := model.GetNowStagingConfig()
	staging, _ := staging_service.GetStagingInfoById(id)
	logging := model.Logging{
		Username:    operator,
		StagingName: staging.StagingName,
		Operation:   "合同号：" + contractNo + "添加身份证",
		Details:     data,
	}
	logging.Create()
	return nil
}

func SupplementContract(r io.Reader) (int, error) {
	contracts, err := parseExcel(r)
	if err != nil {
		return 0, err
	}
	for _, contract := range *contracts {
		_ = model.UpdateTargetPlacementAreaAndTemporaryRelocationArea(contract.ContractNo, contract.TargetPlacementArea, contract.TemporaryRelocationArea)
	}
	return len(*contracts), nil
}

func parseExcel(r io.Reader) (*[]model.Contract, error) {
	var contracts []model.Contract
	exFile, err := excel.Open(r)
	if err != nil {
		return nil, err
	}
	results, err := exFile.GetSheetData("668355") //表名记号
	if err != nil {
		return nil, err
	}
	resultsLen := len(results)
	indexColumn := make(map[string]int)
	var initialLine int
	for initialLine = 0; initialLine < resultsLen; initialLine++ {
		for index, column := range results[initialLine] {
			//如果找到合同号关键字，则将此行记录识别为头部标题
			switch column {
			case "拆迁安置补偿协议合同号":
				indexColumn["ContractNo"] = index
				break
			case "指标安置面积":
				indexColumn["TargetPlacementArea"] = index
				break
			case "计算临迁费面积":
				indexColumn["TemporaryRelocationArea"] = index
				break
			}
		}
		if len(indexColumn) > 0 {
			initialLine++
			break
		}
	}
	if initialLine >= resultsLen {
		return nil, errors.BadError("无效数据")
	}
	for i := initialLine; i < resultsLen; i++ {
		if len(results[i]) > 0 {
			targetPlacementArea, _ := convert.StrToFloat64(results[i][indexColumn["TargetPlacementArea"]], 4)
			temporaryRelocationArea, _ := convert.StrToFloat64(results[i][indexColumn["TemporaryRelocationArea"]], 4)
			contracts = append(contracts, model.Contract{
				ContractNo:              results[i][indexColumn["ContractNo"]],
				TargetPlacementArea:     targetPlacementArea,
				TemporaryRelocationArea: temporaryRelocationArea,
			})
		}
	}
	return &contracts, nil
}

func QueryContractByFuzzyNameAndStagingId(stagingId uint, filterName string, page uint, pageSize uint) (data *model.PaginationQ, err error) {
	// 判断stagingId是否存在
	if !model.ExistsStagingById(stagingId) {
		return nil, errors.BadError("该期数不存在")
	}

	data, err = model.QueryContractByFuzzyNameAndStagingId(stagingId, filterName, page, pageSize)
	return
}

func QueryDeclarationCountByContractOn(contractOn string) (result map[string]interface{}, err error) {
	// 判断此 contractOn 是否存在
	if !model.ExistsContractByContractOn(contractOn) {
		return nil, errors.BadError("不存在此合同")
	}

	count, active, err := model.QueryDeclarationCountByContractOn(contractOn)

	if err != nil {
		return nil, err
	}
	result = make(map[string]interface{}, 0)
	result["result_count"] = count
	result["declaration_active"] = active

	return

}

func SetCanDeclareByStagingId(stagingId uint, status bool, operator string) error {

	// 先查询该分期所有合同 (can_declare)
	contracts, err := model.QueryAllContractByStagingId(stagingId)
	if err != nil {
		return err
	}
	if len(contracts) == 0 {
		return errors.New("没有找到该期数的合同")
	}
	changeSuccess := ""
	// 循环判断该合同是否符合条件
	for _, contract := range contracts {
		c, err := model.FindContractById(contract.ContractNo)
		if err != nil {
			return errors.New("查询合同出错:" + c.ContractNo)
		}
		// 判断合同状态是否一样
		if c.CanDeclare == status {
			continue
		} else {
			if c.CanDeclare == true {
				// 要修改成不可申报
				c.CanDeclare = false
				err := model.UpdateContract(c)
				if err != nil {
					return errors.New("合同修改状态失败:" + c.ContractNo)
				}
			} else {
				// 修改成可申报 c
				err := c.IsCanDeclareAble()
				if err != nil {
					return err
				}
				c.CanDeclare = true
				err = model.UpdateContract(c)
				if err != nil {
					return errors.New("合同修改状态失败:" + c.ContractNo)
				}
			}
		}
		// 修改成功
		changeSuccess += c.ContractNo + " "
	}

	if changeSuccess == "" {
		return errors.New("没有可需要改变状态的合同")
	} else {
		if err != nil {
			return err
		}
		// 注意添加logging
		id, _ := model.GetNowStagingConfig()
		if err != nil {
			return err
		}
		staging, _ := staging_service.GetStagingInfoById(id)
		logging := model.Logging{
			Username:    operator,
			StagingName: staging.StagingName,
			Operation:   "根据期数修改合同申报状态：" + string(stagingId),
			Details:     changeSuccess + "===>修改状态",
		}
		logging.Create()
	}

	return nil
}

// UpdateContract 根据合同号修改信息
func UpdateContract(body api.ContractUpdateBody, operator string) error {
	if !model.ExistsContractByContractOn(body.ContractNo) {
		return errors.New("不存在此合同号")
	}

	if !model.ExistsStagingById(body.StagingID) {
		return errors.New("不存在此期数")
	}
	// update message
	err2 := model.UpdateContractByContractOn(body)
	if err2 != nil {
		return err2
	}
	// 注意添加logging
	id, _ := model.GetNowStagingConfig()
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	staging, _ := staging_service.GetStagingInfoById(id)
	logging := model.Logging{
		Username:    operator,
		StagingName: staging.StagingName,
		Operation:   "更新合同数据：" + string(body.ContractNo),
		Details:     string(data),
	}
	logging.Create()
	return nil
}

func UpdateContractCanDeclareByContractOn(contractList []string, canDeclare bool, operator string) error {
	NotExistsContract := ""
	// 判断所有合同号是否都存在
	for _, s := range contractList {
		if !model.ExistsContractByContractOn(s) {
			NotExistsContract += " ->" + s
		}
	}
	if NotExistsContract != "" {
		return errors.New("不存在合同号" + NotExistsContract)
	}

	changeSuccess := ""
	// 循环判断该合同是否符合条件
	for _, contract := range contractList {
		c, err := model.FindContractById(contract)
		if err != nil {
			return errors.New("查询合同出错:" + c.ContractNo)
		}
		// 判断合同状态是否一样
		if c.CanDeclare == canDeclare {
			continue
		} else {
			if c.CanDeclare == true {
				// 要修改成不可申报
				c.CanDeclare = false
				err := model.UpdateContract(c)
				if err != nil {
					return errors.New("合同修改状态失败:" + c.ContractNo)
				}
			} else {
				// 修改成可申报 c
				err := c.IsCanDeclareAble()
				if err != nil {
					return err
				}
				c.CanDeclare = true
				err = model.UpdateContract(c)
				if err != nil {
					return errors.New("合同修改状态失败:" + c.ContractNo)
				}
			}
		}
		// 修改成功
		changeSuccess += c.ContractNo + " "
	}

	// 注意添加logging
	id, _ := model.GetNowStagingConfig()
	staging, _ := staging_service.GetStagingInfoById(id)
	logging := model.Logging{
		Username:    operator,
		StagingName: staging.StagingName,
		Operation:   "修改多合同号的是否可申报：" + changeSuccess,
		Details:     changeSuccess + "是否可申报修改为" + strconv.FormatBool(canDeclare),
	}
	logging.Create()
	return nil
}
