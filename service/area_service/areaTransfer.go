/**
* @Author: SAmwake
* @Date: 2021/11/12 17:17
 */
package area_service

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"relocate/model"
	"relocate/util/errors"
	"time"
)

func AreaTransfer(contractNo string, targetContractNo string, area float64, areaUnitPrice float64) error {

	// 面积转移表创建条件
	// 合同必须存在
	// 上传面积不能为负数
	// 合同是否满足可申报条件
	// 该合同下是否有其他面积转移没有处理
	if !model.ExistsContractByContractOn(contractNo) {
		return errors.New("合同号不存在")
	}
	if !model.ExistsContractByContractOn(targetContractNo) {
		return errors.New("目标合同号不存在")
	}
	if area < 0 || area > 9999 {
		return errors.New("请传入正确面积参数")
	}
	if areaUnitPrice < 0 {
		return errors.New("请传入正确面积单价参数")
	}
	if model.IsAreaContractNotProcessed(contractNo) {
		return errors.New("该合同尚有未处理面积转移合同")
	}

	preData, err := model.FindContractById(contractNo)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("没有找到此用户")
		}
		return err
	}
	if preData.CanDeclare == false {
		return errors.New("该合同尚未能申报")
	}
	preUser, err := model.FindUserByPhone(preData.PhoneNumber1)
	if err != nil {
		return err
	}

	afterData, err := model.FindContractById(targetContractNo)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("没有找到此用户")
		}
		return err
	}
	if afterData.CanDeclare == false {
		return errors.New("目标合同尚未能申报")
	}
	afterUser, err := model.FindUserByPhone(afterData.PhoneNumber1)
	if err != nil {
		return err
	}

	// 万一手机号相同
	if preData.PhoneNumber1 == afterData.PhoneNumber1 ||
		preData.PhoneNumber1 == afterData.PhoneNumber2 ||
		preData.PhoneNumber2 == afterData.PhoneNumber1 ||
		preData.PhoneNumber2 == afterData.PhoneNumber2 {
		return errors.New("目标合同尚未能申报")
	}

	areaContract := model.AreaTransfer{
		PhoneNumber:       preUser.PhoneNumber,
		Name:              preUser.Name,
		ContractNo:        contractNo,
		Status:            model.UnderWag,
		TargetContractNo:  targetContractNo,
		TargetPhoneNumber: afterUser.PhoneNumber,
		TargetName:        afterUser.Name,
		Area:              area,
		DateOfSigning:     time.Now().String(),
		AreaUnitPrice:     areaUnitPrice,
	}

	err = areaContract.Create()
	if err != nil {
		return errors.New("创建面积合同失败")
	}

	return nil
}

func AreaContractUpdateStatus(id uint, statue bool, operator string) error {
	//  面积合同修改条件
	//  如果该合同之下有申报进行中的申报（回迁面积尚未删减） 不可转移面积

	var areaContract model.AreaTransfer
	areaContract.ID = id
	err := areaContract.Get()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("没有找到此面积合同单")
		}
		return err
	}

	option := 0
	if statue == true {
		option = 1
	}
	if int(areaContract.Status) == option {
		return errors.New("状态相同，无需修改")
	}

	if option == 0 {
		if areaContract.Status == 1 {
			return errors.New("合同已经生效，不可直接作废")
		}
	}

	ok, err := model.IsHasDeclarationNotProcessedByContractNo(areaContract.ContractNo)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}
	if ok {
		return errors.New("当前合同有申报状态进行中的申报")
	}

	contract, err := model.FindContractById(areaContract.ContractNo)
	if err != nil {
		return errors.New("查询合同失败")
	}
	targetContract, err := model.FindContractById(areaContract.TargetContractNo)
	if err != nil {
		return errors.New("查询合同失败")
	}

	// 当前合同剩余面积够不够扣
	if contract.RemainingHQArea < areaContract.Area {
		return errors.New("当前合同回迁面积小于需要转移的面积")
	}

	user, err := model.FindUserByPhone(areaContract.PhoneNumber)
	if err != nil {
		return err
	}
	targetUser, err := model.FindUserByPhone(areaContract.TargetPhoneNumber)
	if err != nil {
		return err
	}

	AreaDecimal := decimal.NewFromFloat(areaContract.Area)

	// 当前合同面积减去转移面积
	PreAreaDecimal := decimal.NewFromFloat(contract.InitialHQArea)
	truePreAreaDecimal, _ := PreAreaDecimal.Float64()
	PreAreaAfterDecimal := PreAreaDecimal.Sub(AreaDecimal)
	truePreAreaAfterDecimal, _ := PreAreaAfterDecimal.Float64()

	contract.InitialHQArea = truePreAreaAfterDecimal

	// 目标合同添加面积
	targetAreaDecimal := decimal.NewFromFloat(targetContract.InitialHQArea)
	truetargetAreaDecimal, _ := targetAreaDecimal.Float64()
	targetAreaAfterDecimal := targetAreaDecimal.Add(AreaDecimal)
	trueTargetAreaAfterDecimal, _ := targetAreaAfterDecimal.Float64()

	targetContract.InitialHQArea = trueTargetAreaAfterDecimal
	areaContract.Status = 1
	areaPriceTotal := areaContract.AreaUnitPrice * areaContract.Area

	areaTransferDetail := model.AreaTransferDetail{
		PhoneNumber:       user.PhoneNumber,
		Name:              user.Name,
		ContractNo:        areaContract.ContractNo,
		IdNumber:          user.IdNumber,
		Area:              areaContract.Area,
		PreTransferArea:   truePreAreaDecimal,
		AfterTransferArea: truePreAreaAfterDecimal,

		TargetIdNumber:          targetUser.IdNumber,
		TargetContractNo:        targetContract.ContractNo,
		TargetPhoneNumber:       targetUser.PhoneNumber,
		TargetName:              targetUser.Name,
		TargetPreTransferArea:   truetargetAreaDecimal,
		TargetAfterTransferArea: trueTargetAreaAfterDecimal,

		AreaUnitPrice:  areaContract.AreaUnitPrice,
		AreaPriceTotal: areaPriceTotal,
		Operator:       operator,
		DateOfSigning:  time.Now().String(),
	}

	err = areaTransferDetail.Create()
	if err != nil {
		return errors.New("创建面积明细失败")
	}

	err = contract.Update()
	if err != nil {
		return errors.New("合同面积更新失败")
	}

	err = targetContract.Update()
	if err != nil {
		return errors.New("合同面积更新失败")
	}

	err = areaContract.Update()
	if err != nil {
		return errors.New("面积表更新状态失败")
	}

	// 创建面积转出明细表
	areaDetial := model.AreaDetails{
		ContractNo:         contract.ContractNo,
		OperationalDetails: fmt.Sprintf("面积转出到合同：%v,面积数%v", areaContract.TargetContractNo, areaContract.Area),
		OperationalArea:    fmt.Sprintf("操作面积 -%v", areaContract.Area),
		RemainingArea:      fmt.Sprintf("剩余面积 %v", contract.InitialHQArea),
	}

	err = areaDetial.Create()
	if err != nil {
		return errors.New("创建面积详情失败")
	}

	// 创建面积转入明细表
	targetAreaDetial := model.AreaDetails{
		ContractNo:         targetContract.ContractNo,
		OperationalDetails: fmt.Sprintf("面积转入到合同：%v,面积数%v", areaContract.ContractNo, areaContract.Area),
		OperationalArea:    fmt.Sprintf("操作面积 +%v", areaContract.Area),
		RemainingArea:      fmt.Sprintf("剩余面积 %v", targetContract.InitialHQArea),
	}
	err = targetAreaDetial.Create()
	if err != nil {
		return errors.New("创建面积详情失败")
	}

	// 创建日志
	var staging model.Staging
	staging.ID = contract.StagingID
	if err = staging.Get(); err != nil {
		return err
	}
	data, err := json.Marshal(&areaTransferDetail)
	if err != nil {
		return err
	}

	logger := model.Logging{
		Username:    operator,
		StagingName: staging.StagingName,
		Operation:   fmt.Sprintf("合同号%v 转移面积共%v 到合同号%v", contract.ContractNo, areaContract.Area, targetContract.ContractNo),
		Details:     string(data),
	}

	if err := logger.Create(); err != nil {
		return err
	}

	return nil

}
