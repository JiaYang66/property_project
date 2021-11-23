package result_service

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"relocate/model"
	"relocate/util/errors"
	"relocate/util/random"
	"relocate/util/times"
	"strconv"
)

func ExportExcel(resultDataList []model.ResultData) (*excelize.File, string, error) {
	file := excelize.NewFile()
	sheet := "Sheet1"
	file.SetCellStr(sheet, "A1", "序号")
	file.SetSheetRow(sheet, "B1", &[]interface{}{
		"结果ID",
		"分期期数ID",
		"分期期数名称",
		"合同号",
		"被拆迁人",
		"申报ID",
		"申报户型ID",
		"申报户型的栋号",
		"申报户型的户型",
		"申报面积㎡",
		"申报面积描述",
		"楼号",
		"房间号",
		"公示状态",
		"录入结果人员",
	})
	for i, resultData := range resultDataList {
		file.SetCellInt(sheet, fmt.Sprintf("A%d", i+2), i+1)
		file.SetSheetRow(sheet, fmt.Sprintf("B%d", i+2), &[]interface{}{
			resultData.ID,
			resultData.StagingID,
			resultData.StagingName,
			resultData.ContractNo,
			resultData.Peoples,
			resultData.DeclarationID,
			resultData.DeclarationHuxingID,
			resultData.DeclarationBuildingNo,
			resultData.DeclarationHuxingNo,
			resultData.DeclarationArea,
			resultData.DeclarationAreaShow,
			resultData.BuildingNo,
			resultData.RoomNo,
			resultData.PublicityStatus,
			resultData.Operator,
		})
	}
	fileName := fmt.Sprintf("摇珠表-%s-%s.xlsx", times.ToStr(), random.String(6))
	//path := fmt.Sprintf("./excel/export/%s", fileName)
	//err := file.SaveAs(path)
	//if err != nil {
	//	return nil, "", err
	//}
	return file, fileName, nil
}

func UpdateResultPublicStat(declarationId []int, publicityStatus bool) error {
	notExistsDeclarationId := ""
	for _, data := range declarationId {
		if !model.ExistResultByDeclarationId(data) {
			s := strconv.Itoa(data)
			notExistsDeclarationId += s + " "
		}
	}

	if notExistsDeclarationId != "" {
		return errors.New("不存在 " + notExistsDeclarationId + " 申报号")
	}

	// 更新
	if err := model.UpdateResultPublicStatusByDeclarationId(declarationId, publicityStatus); err != nil {
		return err
	}
	return nil

}
