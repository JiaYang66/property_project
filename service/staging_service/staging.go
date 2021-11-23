package staging_service

import (
	"github.com/jinzhu/gorm"
	"io"
	"relocate/model"
	"relocate/util/convert"
	"relocate/util/errors"
	"relocate/util/excel"
	"relocate/util/times"
	"strings"
)

func GetStagingInfoById(stagingID uint) (*model.Staging, error) {
	staging, err := model.FindStagingById(stagingID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.BadError("分期期数不存在")
		}
		return nil, err
	}
	return staging, nil
}

func AddStaging(stagingName string, ok bool) error {

	if strings.TrimSpace(stagingName) == "" {
		return errors.BadError("分期名不能为空")
	}
	if ok {
		var staging model.Staging
		id, err := model.GetNowStagingConfig()
		if err != nil {
			return errors.New("未配置全局期数")
		}
		staging.ID = id
		staging.StagingName = stagingName
		err = staging.Update()
		if err != nil {
			return err
		}
		return nil
	}
	if err := model.SaveStaging(stagingName); err != nil {
		return err
	}

	return nil
}

func ImportContract(stagingID uint, r io.Reader, operator string) (int, error) {
	staging, err := GetStagingInfoById(stagingID)
	if err != nil {
		return 0, err
	}
	contracts, err := parseExcel(staging.ID, r)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, contract := range *contracts {
		if err := contract.Create(); err != nil {
			//如果合同号冲突，更改分期数id
			if strings.Contains(err.Error(), "Duplicate entry") {
				model.UpdateContractStagingID(contract.ContractNo, stagingID)
				continue
			}
			return 0, err
		}
		count++
	}
	logging := model.Logging{
		Username:    operator,
		StagingName: staging.StagingName,
		Operation:   "导入合同单文件原始数据",
	}
	logging.Create()
	return count, nil
}

//解析Excel
func parseExcel(stagingID uint, r io.Reader) (*[]model.Contract, error) {
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
			case "合同号":
				indexColumn["ContractNo"] = index
				break
			case "社别":
				indexColumn["SocialCategory"] = index
				break
			case "被拆迁人":
				indexColumn["Peoples"] = index
				break
			case "身份证号码":
				indexColumn["CardNumber"] = index
				break
			case "房屋栋号":
				indexColumn["HouseNumber"] = index
				break
			case "被拆迁房屋地址":
				indexColumn["OldAddress"] = index
				break
			case "手机号码1":
				indexColumn["PhoneNumber1"] = index
				break
			case "手机号码2":
				indexColumn["PhoneNumber2"] = index
				break
			case "签署协议日期":
				indexColumn["DateOfSigning"] = index
				break
			case "交楼日期":
				indexColumn["DateOfDelivery"] = index
				break
			case "证载产权人":
				indexColumn["Proprietor"] = index
				break
			case "签约人":
				indexColumn["Signatory"] = index
				break
			case "变更方式":
				indexColumn["ChangeMethod"] = index
				break
			case "集体土地房产证字":
				indexColumn["CollectiveLandPropertyCertificate"] = index
				break
			case "登记字号":
				indexColumn["Registration"] = index
				break
			case "剩余总回迁面积":
				indexColumn["RemainingHQArea"] = index
				indexColumn["InitialHQArea"] = index
				break
			case "是否交齐楼":
				indexColumn["IsDelivery"] = index
				break
			case "指标面积":
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
			isDelivery := false
			canDeclare := false
			if results[i][indexColumn["IsDelivery"]] == "Y" {
				isDelivery = true
				canDeclare = true
			}
			initialHQArea, _ := convert.StrToFloat64(results[i][indexColumn["InitialHQArea"]], 4)
			remainingHQArea, _ := convert.StrToFloat64(results[i][indexColumn["RemainingHQArea"]], 4)
			targetPlacementArea, _ := convert.StrToFloat64(results[i][indexColumn["TargetPlacementArea"]], 4)
			temporaryRelocationArea, _ := convert.StrToFloat64(results[i][indexColumn["TemporaryRelocationArea"]], 4)
			contracts = append(contracts, model.Contract{
				ContractNo:                        results[i][indexColumn["ContractNo"]],
				SocialCategory:                    results[i][indexColumn["SocialCategory"]],
				Peoples:                           replaceSpace(results[i][indexColumn["Peoples"]]),
				CardNumber:                        replaceSpace(results[i][indexColumn["CardNumber"]]),
				HouseNumber:                       results[i][indexColumn["HouseNumber"]],
				OldAddress:                        results[i][indexColumn["OldAddress"]],
				PhoneNumber1:                      results[i][indexColumn["PhoneNumber1"]],
				PhoneNumber2:                      results[i][indexColumn["PhoneNumber2"]],
				DateOfSigning:                     times.ConvertToFormatDay(results[i][indexColumn["DateOfSigning"]]),
				DateOfDelivery:                    times.ConvertToFormatDay(results[i][indexColumn["DateOfDelivery"]]),
				Proprietor:                        results[i][indexColumn["Proprietor"]],
				Signatory:                         results[i][indexColumn["Signatory"]],
				ChangeMethod:                      results[i][indexColumn["ChangeMethod"]],
				CollectiveLandPropertyCertificate: results[i][indexColumn["CollectiveLandPropertyCertificate"]],
				Registration:                      results[i][indexColumn["Registration"]],
				InitialHQArea:                     initialHQArea,
				RemainingHQArea:                   remainingHQArea,
				Desc:                              "备注信息",
				IsDelivery:                        isDelivery,
				HouseWriteOff:                     false,
				CanDeclare:                        canDeclare,
				Single:                            JudgmentSingle(results[i][indexColumn["CardNumber"]]),
				StagingID:                         stagingID,
				TargetPlacementArea:               targetPlacementArea,
				TemporaryRelocationArea:           temporaryRelocationArea,
			})
		}
	}
	return &contracts, nil
}

//判断是单一产权人还是共有产权人
func JudgmentSingle(cardNumber string) bool {
	return !strings.Contains(replaceSpace(cardNumber), ",")
}

func replaceSpace(str string) string {
	value := ""
	context := strings.Fields(str)
	for i, v := range context {
		if i == len(context)-1 {
			value += v
		} else {
			value += v + ","
		}
	}
	return value
}
