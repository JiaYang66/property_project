package declaration_service

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"log"
	"relocate/api"
	"relocate/model"
	"relocate/util/errors"
	"relocate/util/times"
	"strconv"
	"strings"
)

//添加申报表
//func AddDeclaration(body api.AddDeclarationBody, operator string, ok bool) error {
//	//获取日志记录员
//	record := operator
//	//获取当前分期
//	stagingId, err := model.GetNowStagingConfig()
//	if err != nil {
//		return err
//	}
//	if stagingId == 0 {
//		return errors.BadError("配置未初始化")
//	}
//	//获取合同单
//	contract, err := model.FindContractById(body.ContractNo)
//	if err != nil {
//		if err == gorm.ErrRecordNotFound {
//			return errors.BadError("该合同单不存在")
//		}
//		return err
//	}
//	if !contract.CanDeclare {
//		return errors.BadError("该合同单无法申报")
//	}
//	//判断操作者权限
//	if ok {
//		admin, err := model.GetAdminInfo(operator)
//		if err != nil {
//			if err == nil {
//				return errors.BadError("不存在该管理员")
//			}
//			return err
//		}
//		operator = admin.AdminSignname
//	} else {
//		user, err := model.GetUserInfo(operator)
//		if err != nil {
//			if err == gorm.ErrRecordNotFound {
//				return errors.BadError("不存在该用户")
//			}
//		}
//		if user.UserStatus == 0 || user.UserStatus == 1 {
//			return errors.BadError("用户校验未通过")
//		}
//		operator = user.Name
//		log.Println(contract.CardNumber)
//		log.Println(user.IdNumber)
//		if !strings.Contains(contract.CardNumber, user.IdNumber) {
//			return errors.BadError("该合同单不属于你，无法申报")
//		}
//	}
//	//获取该合同号在当前分期数下的有效申报
//	list, err := model.GetActiveStateDeclaration(body.ContractNo, stagingId)
//	if err != nil {
//		return err
//	}
//	//获取可选套数
//	number, err := model.GetHuxingOptionalConfig()
//	if err != nil {
//		return err
//	}
//	//判断是否超出当前期数可选套数的数量
//	if uint(len(list)) >= number {
//		return errors.BadError("不可申报，当前期数可选套数你已达标")
//	}
//	//获取户型
//	huxing := new(model.Huxing)
//	huxing, err = huxing.FindHuxingByID(body.HuxingID)
//	if err != nil {
//		return err
//	}
//	//判断是否超过该户型的最大可申报套数
//	var num uint = 0
//	for _, h := range list {
//		if h.DeclarationHuxingID == body.HuxingID {
//			num++
//		}
//	}
//	if num >= huxing.Maximum {
//		return errors.BadError("你申请该户型已达到该户型最多可申请的套数，不可再申请该户型")
//	}
//	//获取现场时间段
//	t := new(model.Time)
//	t, err = t.FindTimeByID(body.TimeID)
//	if err != nil {
//		if err == gorm.ErrRecordNotFound {
//			return errors.BadError("该现场确认时间段不存在")
//		}
//		return err
//	}
//	if t.SelectedNum == t.OptionalNum {
//		return errors.BadError("当前时间段已满")
//	}
//	////获取剩余拆迁面积
//	//lq_area := (*contract).RemainingHQArea
//	////获取户型面积
//	//huxingAreaDecimal, err := decimal.NewFromString(huxing.Area)
//	//if err != nil {
//	//	return err
//	//}
//	//huxing_area, _ := huxingAreaDecimal.Float64()
//	////判断剩余拆迁面积是否大于户型面积
//	//if lq_area < huxing_area {
//	//	return errors.BadError("所选的户型面积大于拆迁剩余面积")
//	//}
//	////计算
//	//remainingHQArea := decimal.NewFromFloat(lq_area)
//	//huxingArea := decimal.NewFromFloat(huxing_area)
//	//area, _ := remainingHQArea.Sub(huxingArea).Float64()
//	//contract.RemainingHQArea = area
//	declaration := model.Declaration{
//		DeclarationHuxingID:   huxing.ID,
//		DeclarationHuxingNo:   huxing.HuxingNo,
//		DeclarationBuildingNo: huxing.BuildingNo,
//		DeclarationAreaShow:   huxing.AreaShow,
//		DeclarationArea:       huxing.Area,
//		TimeID:                body.TimeID,
//		TimeName:              t.Name,
//		StagingID:             stagingId,
//		ContractNo:            body.ContractNo,
//		DeclarationStatus:     model.DeclarationStatus(0),
//		ActiveState:           true,
//		WinningStatus:         model.WinningStatus(-1),
//		Operator:              operator,
//		Printer:               "",
//		Trustee:               body.Trustee,
//		TrusteeCardNumber:     body.TrusteeCardNumber,
//		TrusteePhoneNumber:    body.TrusteePhoneNumber,
//		TrusteeRelationship:   body.TrusteeRelationship,
//	}
//	if err := declaration.Create(); err != nil {
//		return err
//	}
//	//更新现场时间段的已选人数
//	t.SelectedNum++
//	if err := t.Update(); err != nil {
//		return err
//	}
//	//更改合同单数据
//	if err := model.UpdateContract(contract); err != nil {
//		return err
//	}
//	////创建面积明细表
//	//areaDetail := model.AreaDetails{
//	//	ContractNo:         contract.ContractNo,
//	//	OperationalDetails: "申报",
//	//	OperationalArea:    "-" + huxing.Area + "㎡",
//	//	RemainingArea:      remainingHQArea.Sub(huxingArea).String() + "㎡",
//	//}
//	//err = areaDetail.Create()
//	//if err != nil {
//	//	return err
//	//}
//	//管理员添加申报创建日志
//	if ok {
//		data, err := json.Marshal(body)
//		if err != nil {
//			return err
//		}
//		//获取当前分期的名称 -- 创建日志
//		staging, _ := model.FindStagingById(stagingId)
//		logging := model.Logging{
//			Username:    record,
//			StagingName: staging.StagingName,
//			Operation:   contract.ContractNo + "申报成功",
//			Details:     string(data),
//		}
//		logging.Create()
//	}
//	return nil
//}

//添加申报表,2021-3-26改
func AddDeclarationNew(body api.AddDeclarationBody, operator string, ok bool) error {
	//获取日志记录员
	record := operator
	//获取当前分期
	stagingId, err := model.GetNowStagingConfig()
	if err != nil {
		return err
	}
	if stagingId == 0 {
		return errors.BadError("配置未初始化")
	}
	//获取当前轮
	rounds, err := model.GetNowRoundsConfig()
	if err != nil {
		return err
	}
	if rounds == 0 {
		return errors.BadError("配置未初始化")
	}
	//获取合同单
	contract, err := model.FindContractById(body.ContractNo)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.BadError("该合同单不存在")
		}
		return err
	}
	if !contract.CanDeclare {
		return errors.BadError("该合同单无法申报")
	}
	//判断操作者权限
	if ok {
		admin, err := model.GetAdminInfo(operator)
		if err != nil {
			if err == nil {
				return errors.BadError("不存在该管理员")
			}
			return err
		}
		operator = admin.AdminSignname
	} else {
		user, err := model.GetUserInfo(operator)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.BadError("不存在该用户")
			}
		}
		if user.UserStatus == 0 || user.UserStatus == 1 {
			return errors.BadError("用户校验未通过")
		}
		operator = user.Name
		log.Println(contract.CardNumber)
		log.Println(user.IdNumber)
		if !strings.Contains(contract.CardNumber, user.IdNumber) {
			return errors.BadError("该合同单不属于你，无法申报")
		}
	}
	//获取该合同号在当前分期数下的有效申报
	list, err := model.GetActiveStateDeclaration(body.ContractNo, stagingId, rounds)
	if err != nil {
		return err
	}
	//获取可选套数 返回json数据
	hxConfig, err := model.GetHuxingGroupingOptionalConfig()
	if err != nil {
		return err
	}
	for _, configJson := range hxConfig {
		var groupingNum = 0
		var is = false
		split := strings.Split(configJson.IDs, ",")
		for _, s := range split {
			// 数据库里面的
			for _, declaration := range list {
				if strconv.Itoa(int(declaration.DeclarationHuxingID)) == s {
					groupingNum++
				}
				if strconv.Itoa(int(body.HuxingID)) == s {
					is = true
				}
			}
		}
		log.Println(groupingNum, configJson.IDs, configJson.Num)
		if groupingNum >= configJson.Num && is {
			return errors.BadError("不可申报，当前期数可选套数已达标")
		}
	}
	//获取户型
	huxing := new(model.Huxing)
	huxing.ID = body.HuxingID
	err = huxing.FindHuxingByID()
	if err != nil {
		return err
	}
	//判断是否超过该户型的最大可申报套数
	var num uint = 0
	for _, h := range list {
		if h.DeclarationHuxingID == body.HuxingID {
			num++
		}
	}
	if num >= huxing.Maximum {
		return errors.BadError("你申请该户型已达到该户型最多可申请的套数，不可再申请该户型")
	}
	//获取现场时间段
	t := new(model.Time)
	t.ID = body.TimeID
	err = t.FindTimeByID()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.BadError("该现场确认时间段不存在")
		}
		return err
	}
	if t.SelectedNum == t.OptionalNum {
		return errors.BadError("当前时间段已满")
	}
	////获取剩余拆迁面积
	//lq_area := (*contract).RemainingHQArea
	////获取户型面积
	//huxingAreaDecimal, err := decimal.NewFromString(huxing.Area)
	//if err != nil {
	//	return err
	//}
	//huxing_area, _ := huxingAreaDecimal.Float64()
	////判断剩余拆迁面积是否大于户型面积
	//if lq_area < huxing_area {
	//	return errors.BadError("所选的户型面积大于拆迁剩余面积")
	//}
	////计算
	//remainingHQArea := decimal.NewFromFloat(lq_area)
	//huxingArea := decimal.NewFromFloat(huxing_area)
	//area, _ := remainingHQArea.Sub(huxingArea).Float64()
	//contract.RemainingHQArea = area
	declaration := model.Declaration{
		DeclarationHuxingID:   huxing.ID,
		DeclarationHuxingNo:   huxing.HuxingNo,
		DeclarationBuildingNo: huxing.BuildingNo,
		DeclarationAreaShow:   huxing.AreaShow,
		DeclarationArea:       huxing.Area,
		TimeID:                body.TimeID,
		TimeName:              t.Name,
		StagingID:             stagingId,
		Rounds:                rounds,
		ContractNo:            body.ContractNo,
		DeclarationStatus:     model.DeclarationStatus(0),
		ActiveState:           true,
		WinningStatus:         model.WinningStatus(-1),
		Operator:              operator,
		Printer:               "",
		Trustee:               body.Trustee,
		TrusteeCardNumber:     body.TrusteeCardNumber,
		TrusteePhoneNumber:    body.TrusteePhoneNumber,
		TrusteeRelationship:   body.TrusteeRelationship,
	}
	if err := declaration.Create(); err != nil {
		return err
	}
	//更新现场时间段的已选人数
	t.SelectedNum++
	if err := t.Update(); err != nil {
		return err
	}
	//更改合同单数据
	if err := model.UpdateContract(contract); err != nil {
		return err
	}
	////创建面积明细表
	//areaDetail := model.AreaDetails{
	//	ContractNo:         contract.ContractNo,
	//	OperationalDetails: "申报",
	//	OperationalArea:    "-" + huxing.Area + "㎡",
	//	RemainingArea:      remainingHQArea.Sub(huxingArea).String() + "㎡",
	//}
	//err = areaDetail.Create()
	//if err != nil {
	//	return err
	//}
	//管理员添加申报创建日志
	if ok {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		//获取当前分期的名称 -- 创建日志
		staging, _ := model.FindStagingById(stagingId)
		logging := model.Logging{
			Username:    record,
			StagingName: staging.StagingName,
			Operation:   contract.ContractNo + "申报成功",
			Details:     string(data),
		}
		logging.Create()
	}
	return nil
}

//更改申报表申报状态
func ChangeDeclarationStatus(declarationID uint, status int, operator string) error {
	//获取日志记录员
	record := operator
	admin, err := model.GetAdminInfo(operator)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.BadError("该管理员不存在")
		}
	}
	operator = admin.AdminSignname
	declaration, err := model.FindDeclarationByID(declarationID)
	if err != nil {
		return err
	}
	//如果是待抽签状态就不可以改
	if declaration.WinningStatus.String() != model.WinningStatus(-1).String() {
		return errors.BadError("该申报已录入结果，无法更改申报状态")
	}
	if model.DeclarationStatus(status) == declaration.DeclarationStatus {
		return nil
	}
	//获取合同单
	contract, err := model.FindContractById(declaration.ContractNo)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.BadError("该合同单不存在")
		}
		return err
	}
	//获取剩余拆迁面积
	lq_area := (*contract).RemainingHQArea
	//获取户型面积
	huxingAreaDecimal, err := decimal.NewFromString(declaration.DeclarationArea)
	if err != nil {
		return err
	}
	huxing_area, _ := huxingAreaDecimal.Float64()
	//如果更改为已确认，则减少合同单的剩余面积，添加面积明细记录
	log.Println(status)
	if status == 1 {
		//判断剩余拆迁面积是否小于户型面积
		if lq_area < huxing_area {
			return errors.BadError("所选的户型面积大于拆迁剩余面积，无法确认")
		}
		//计算 --减
		remainingHQArea := decimal.NewFromFloat(lq_area)
		huxingArea := decimal.NewFromFloat(huxing_area)
		area, _ := remainingHQArea.Sub(huxingArea).Float64()
		//更新合同单数据
		contract.RemainingHQArea = area
		if err := model.UpdateContract(contract); err != nil {
			return err
		}
		//创建面积明细表
		areaDetail := model.AreaDetails{
			ContractNo:         contract.ContractNo,
			OperationalDetails: "申报",
			OperationalArea:    "-" + declaration.DeclarationArea + "㎡",
			RemainingArea:      remainingHQArea.Sub(huxingArea).String() + "㎡",
		}
		err = areaDetail.Create()
		if err != nil {
			return err
		}
	} else {
		//计算 --加
		remainingHQArea := decimal.NewFromFloat(lq_area)
		huxingArea := decimal.NewFromFloat(huxing_area)
		area, _ := remainingHQArea.Add(huxingArea).Float64()
		//更新合同单数据
		contract.RemainingHQArea = area
		if err := model.UpdateContract(contract); err != nil {
			return err
		}
		//创建面积明细表
		areaDetail := model.AreaDetails{
			ContractNo:         contract.ContractNo,
			OperationalDetails: "申报被驳回",
			OperationalArea:    "+" + declaration.DeclarationArea + "㎡",
			RemainingArea:      remainingHQArea.Add(huxingArea).String() + "㎡",
		}
		err = areaDetail.Create()
		if err != nil {
			return err
		}
	}
	declaration.DeclarationStatus = model.DeclarationStatus(status)
	declaration.Operator = operator
	err = declaration.Update()
	data, err := json.Marshal(map[string]interface{}{
		"declarationID": declarationID,
		"status":        model.DeclarationStatus(status),
	})
	if err != nil {
		return err
	}
	staging, err := model.FindStagingById(declaration.StagingID)
	if err != nil {
		return err
	}
	logging := model.Logging{
		Username:    record,
		StagingName: staging.StagingName,
		Operation:   "更改合同" + declaration.ContractNo + "申报表的申报状态，改为" + model.DeclarationStatus(status).String(),
		Details:     string(data),
	}
	logging.Create()
	return nil
}

//更改申报表内容
//func UpdateDeclaration(body api.UpdateDeclaration, operator string, ok bool) error {
//	//获取日志记录员
//	record := operator
//	//获取申报表
//	declaration, err := model.FindDeclarationByID(body.DeclarationID)
//	if err != nil {
//		return err
//	}
//	//判断申报表是否已作废
//	if !declaration.ActiveState {
//		return errors.BadError("该申报已失效，无法修改数据")
//	}
//	//判断是报表是否已确认
//	if declaration.DeclarationStatus == model.DeclarationStatus(1) {
//		return errors.BadError("该申报已确认，无法更改数据")
//	}
//	//获取合同单
//	contract, err := model.FindContractById(declaration.ContractNo)
//	if err != nil {
//		if err == gorm.ErrRecordNotFound {
//			return errors.BadError("该合同单不存在")
//		}
//		return err
//	}
//	//判断操作者权限
//	if ok {
//		admin, err := model.GetAdminInfo(operator)
//		if err != nil {
//			if err == nil {
//				return errors.BadError("不存在该管理员")
//			}
//			return err
//		}
//		operator = admin.AdminSignname
//	} else {
//		user, err := model.GetUserInfo(operator)
//		if err != nil {
//			if err == gorm.ErrRecordNotFound {
//				return errors.BadError("不存在该用户")
//			}
//		}
//		if user.UserStatus == 0 || user.UserStatus == 1 {
//			return errors.BadError("用户校验未通过")
//		}
//		operator = user.Name
//		log.Println(contract.CardNumber)
//		log.Println(user.IdNumber)
//		if !strings.Contains(contract.CardNumber, user.IdNumber) {
//			return errors.BadError("不存在该合同单")
//		}
//	}
//	var description string
//	//判断是否更改户型
//	if body.HuxingID != declaration.DeclarationHuxingID {
//		//通过参数传过来的户型id，查找户型
//		newHuxing := new(model.Huxing)
//		newHuxing, err = newHuxing.FindHuxingByID(body.HuxingID)
//		if err != nil {
//			if err == gorm.ErrRecordNotFound {
//				return errors.BadError("该户型不存在")
//			}
//			return err
//		}
//		//获取用户的有效申报
//		list, err := model.GetActiveStateDeclaration(declaration.ContractNo, declaration.StagingID)
//		if err != nil {
//			return err
//		}
//		//获取已申报该户型套数
//		var num uint = 0
//		for _, l := range list {
//			if l.DeclarationHuxingID == body.HuxingID {
//				num++
//			}
//		}
//		//判断是否超出该户型的最大可申请套数
//		if num >= newHuxing.Maximum {
//			return errors.BadError("该户型你已超出户型的最大申报套数，更改失败")
//		}
//		////获取原先的户型面积
//		//areaHu1Decimal, err := decimal.NewFromString(declarationHuxing.Area)
//		//if err != nil {
//		//	return err
//		//}
//		////获取最新的户型面积
//		//areaHu2Decimal, err := decimal.NewFromString(newHuxing.Area)
//		//if err != nil {
//		//	return err
//		//}
//		////获取合同单的剩余面积
//		//hqAreaDecimal := decimal.NewFromFloat(contract.RemainingHQArea)
//		//if err != nil {
//		//	return err
//		//}
//		////判断剩余面积是否小于新申请的户型面积
//		//hqAreaDecimal = hqAreaDecimal.Add(areaHu1Decimal)
//		//if hqAreaDecimal.LessThan(areaHu2Decimal) {
//		//	return errors.BadError("新申报的户型面积大于剩余面积，修改失败")
//		//}
//		////添加面积详细
//		//areaDetail := model.AreaDetails{
//		//	ContractNo:         contract.ContractNo,
//		//	OperationalDetails: "申报表更换户型 -- 原户型：" + declarationHuxing.HuxingNo,
//		//	OperationalArea:    "+" + declarationHuxing.Area + "㎡",
//		//	RemainingArea:      hqAreaDecimal.String() + "㎡",
//		//}
//		//if err := areaDetail.Create(); err != nil {
//		//	return err
//		//}
//		////添加面积详细
//		//hqAreaDecimal = hqAreaDecimal.Sub(areaHu2Decimal)
//		//areaDetail = model.AreaDetails{
//		//	ContractNo:         contract.ContractNo,
//		//	OperationalDetails: "申报表更换户型 -- 最新户型：" + newHuxing.HuxingNo,
//		//	OperationalArea:    "-" + newHuxing.Area + "㎡",
//		//	RemainingArea:      hqAreaDecimal.String() + "㎡",
//		//}
//		//if err := areaDetail.Create(); err != nil {
//		//	return err
//		//}
//		//remainingArea, _ := hqAreaDecimal.Float64()
//		//contract.RemainingHQArea = remainingArea
//		//更新申报表的申请户型
//		declaration.DeclarationArea = newHuxing.Area
//		declaration.DeclarationHuxingNo = newHuxing.HuxingNo
//		declaration.DeclarationHuxingID = newHuxing.ID
//		declaration.DeclarationBuildingNo = newHuxing.BuildingNo
//		declaration.DeclarationAreaShow = newHuxing.AreaShow
//		//更改操作者
//		declaration.Operator = operator
//		if err := declaration.Update(); err != nil {
//			return err
//		}
//		////更新合同表
//		//if err := model.UpdateContract(contract); err != nil {
//		//	return err
//		//}
//		if ok {
//			description += "更改合同单" + declaration.ContractNo +
//				"的申报表的申请户型，更改" + declaration.DeclarationHuxingNo +
//				"为" + newHuxing.HuxingNo + "。"
//		}
//	}
//	//判断受托人信息是否发生改变
//	if body.Trustee != declaration.Trustee ||
//		body.TrusteeRelationship != declaration.TrusteeRelationship ||
//		body.TrusteeCardNumber != declaration.TrusteeCardNumber ||
//		body.TrusteePhoneNumber != declaration.TrusteePhoneNumber {
//		declaration.TrusteePhoneNumber = body.TrusteePhoneNumber
//		declaration.TrusteeCardNumber = body.TrusteeCardNumber
//		declaration.Trustee = body.Trustee
//		declaration.TrusteeRelationship = body.TrusteeRelationship
//		//更改操作者
//		declaration.Operator = operator
//		if err := declaration.Update(); err != nil {
//			return err
//		}
//		if ok {
//			description += "更改合同单" + declaration.ContractNo + "的申报表的受托人信息。"
//		}
//	}
//	//判断是否更新数据并且是管理员更改，更新创建日志
//	if description != "" && ok {
//		data, err := json.Marshal(body)
//		if err != nil {
//			return err
//		}
//		staging, _ := model.FindStagingById(declaration.StagingID)
//		logging := model.Logging{
//			Username:    record,
//			StagingName: staging.StagingName,
//			Operation:   description,
//			Details:     string(data),
//		}
//		logging.Create()
//	}
//	return nil
//}

//更改申报表打印员姓名
func Printing(declarationID uint, operator string) error {
	admin, err := model.GetAdminInfo(operator)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.BadError("该管理员不存在")
		}
	}
	operator = admin.AdminSignname
	declaration, err := model.FindDeclarationByID(declarationID)
	if err != nil {
		return err
	}
	//if declaration.ActiveState {
	//	return errors.BadError("该申报已失效")
	//}
	//if declaration.DeclarationStatus.String() == model.DeclarationStatus(1).String() {
	//	return errors.BadError("该申报表未确认无法打印")
	//}
	declaration.Printer = operator
	if err := declaration.Update(); err != nil {
		return err
	}
	//创建日志记录打印人名字
	staging, err := model.FindStagingById(declaration.StagingID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.BadError("该分期数不存在")
		}
		return err
	}
	logging := model.Logging{
		Username:    operator,
		StagingName: staging.StagingName,
		Details:     "打印合同单" + declaration.ContractNo + "的申报表",
	}
	logging.Create()
	return nil
}

//生成Excel
func GenerateExcelNew(declarationID uint, adminName string) (*excelize.File, string, error) {
	declaration, err := model.FindDeclarationByID(declarationID)
	if err != nil {
		return nil, "", err
	}
	if declaration.DeclarationStatus != model.DeclarationConfirmed {
		return nil, "", errors.BadError("只有确认后的申报表，才可以打印申报表")
	}
	contract, err := model.FindContractById(declaration.ContractNo)
	if err != nil {
		return nil, "", err
	}
	huxingList, err := model.GetAllOptionalHuxing()
	if err != nil {
		return nil, "", err
	}
	//往Excel写入信息
	file, err := excelize.OpenFile("./excel/茅岗路以西城中村改造项目西和元片第一轮回迁安置房户型摸查表.xlsx")
	if err != nil {
		return nil, "", err
	}
	sheetName := "户型摸查表"
	file.SetCellStr(sheetName, "C2", contract.ContractNo)
	file.SetCellStr(sheetName, "G2", contract.SocialCategory)
	file.SetCellStr(sheetName, "C3", contract.Peoples)
	file.SetCellStr(sheetName, "F3", contract.CardNumber)
	file.SetCellStr(sheetName, "C5", contract.PhoneNumber1+"\n"+contract.PhoneNumber2)
	if declaration.Trustee != "" {
		file.SetCellStr(sheetName, "C4", declaration.Trustee)
		file.SetCellStr(sheetName, "F4", declaration.TrusteeCardNumber)
	} else {
		file.SetCellStr(sheetName, "C4", "空")
		file.SetCellStr(sheetName, "F4", "空")
	}
	file.SetCellStr(sheetName, "F5", contract.OldAddress)
	file.SetCellStr(sheetName, "C6", contract.HouseNumber)
	file.SetCellStr(sheetName, "F6", fmt.Sprintf("%.4f ㎡", contract.InitialHQArea))

	//循环申报户型列表
	index := 8
	for i, huxing := range huxingList {
		if i != len(huxingList)-1 {
			file.DuplicateRow(sheetName, index)
		}
		file.SetCellStr(sheetName, fmt.Sprintf("B%d", index), fmt.Sprintf("%d", i+1))
		file.SetCellStr(sheetName, fmt.Sprintf("C%d", index), huxing.AreaShow)
		file.SetCellStr(sheetName, fmt.Sprintf("D%d", index), huxing.HuxingNo)
		file.SetCellStr(sheetName, fmt.Sprintf("E%d", index), huxing.BuildingNo)
		file.SetCellStr(sheetName, fmt.Sprintf("F%d", index), fmt.Sprintf("%d", huxing.Quantity))
		if huxing.ID == declaration.DeclarationHuxingID {
			file.SetCellStr(sheetName, fmt.Sprintf("G%d", index), fmt.Sprintf("(  1  )"))
		} else {
			file.SetCellStr(sheetName, fmt.Sprintf("G%d", index), fmt.Sprintf("(  ×  )"))
		}
		index++
	}
	//合并单元格
	file.MergeCell(sheetName, "A7", fmt.Sprintf("A%d", index-1))
	file.MergeCell(sheetName, fmt.Sprintf("H%d", index-3), fmt.Sprintf("H%d", index-1))
	file.SetCellStr(sheetName, fmt.Sprintf("H%d", index-3), fmt.Sprintf("每宗被拆迁房屋仅能选择一套大户型\n（即三选一）"))

	declarationAreaDecimal, err := decimal.NewFromString(declaration.DeclarationArea)
	if err != nil {
		return nil, "", err
	}
	remainingHQAreaDecimal := decimal.NewFromFloat(contract.RemainingHQArea)
	//获得最终剩余面积
	remainingHQArea, _ := remainingHQAreaDecimal.Add(declarationAreaDecimal).Float64()

	//	file.SetCellStr(sheetName, "D8", fmt.Sprintf("%.4f ㎡", contract.InitialHQArea))
	file.SetCellStr(sheetName, fmt.Sprintf("D%d", index), fmt.Sprintf("%.4f ㎡", remainingHQArea))

	//取消另存为，直接流形式传输
	fileName := fmt.Sprintf("%s-%d-%s.xlsx", contract.ContractNo, declaration.ID, times.ToStr())
	//path := fmt.Sprintf("./excel/%s", fileName)
	//err = file.SaveAs(path)
	//if err != nil {
	//	return nil, "", err
	//}
	//获取当前分期的名称 -- 创建日志
	staging, err := model.FindStagingById(declaration.StagingID)
	logging := model.Logging{
		Username:    adminName,
		StagingName: staging.StagingName,
		Details:     "打印合同单" + declaration.ContractNo + "的申报表",
	}
	logging.Create()
	return file, fileName, nil
}

//生成Excel
func GenerateExcel(declarationID uint, adminName string) (*excelize.File, string, error) {
	declaration, err := model.FindDeclarationByID(declarationID)
	if err != nil {
		return nil, "", err
	}
	if declaration.DeclarationStatus != model.DeclarationConfirmed {
		return nil, "", errors.BadError("只有确认后的申报表，才可以打印申报表")
	}
	contract, err := model.FindContractById(declaration.ContractNo)
	if err != nil {
		return nil, "", err
	}
	huxingList, err := model.GetAllOptionalHuxing()
	if err != nil {
		return nil, "", err
	}
	//往Excel写入信息
	file, err := excelize.OpenFile("./excel/申报表-西华片一期-20201213.xlsx")
	if err != nil {
		return nil, "", err
	}
	sheetName := "版2"
	file.SetCellStr(sheetName, "C2", contract.ContractNo)
	file.SetCellStr(sheetName, "G2", contract.SocialCategory)
	file.SetCellStr(sheetName, "C3", contract.Peoples)
	file.SetCellStr(sheetName, "F3", contract.CardNumber)
	file.SetCellStr(sheetName, "H3", contract.PhoneNumber1+"\n"+contract.PhoneNumber2)
	if declaration.Trustee != "" {
		file.SetCellStr(sheetName, "C4", declaration.Trustee)
		file.SetCellStr(sheetName, "F4", declaration.TrusteeCardNumber)
		file.SetCellStr(sheetName, "H4", declaration.TrusteePhoneNumber)
	} else {
		file.SetCellStr(sheetName, "C4", "空")
		file.SetCellStr(sheetName, "F4", "空")
		file.SetCellStr(sheetName, "H4", "空")
	}
	file.SetCellStr(sheetName, "C5", contract.OldAddress)
	file.SetCellStr(sheetName, "C6", contract.HouseNumber)

	// 更新1218
	declarationAreaDecimal, err := decimal.NewFromString(declaration.DeclarationArea)
	if err != nil {
		return nil, "", err
	}
	remainingHQAreaDecimal := decimal.NewFromFloat(contract.RemainingHQArea)
	//获得最终剩余面积
	remainingHQArea, _ := remainingHQAreaDecimal.Add(declarationAreaDecimal).Float64()

	//	file.SetCellStr(sheetName, "D8", fmt.Sprintf("%.4f ㎡", contract.InitialHQArea))
	file.SetCellStr(sheetName, "D8", fmt.Sprintf("%.4f ㎡", remainingHQArea))
	//循环申报户型列表
	index := 10
	for i, huxing := range huxingList {
		if i != len(huxingList)-1 {
			file.DuplicateRow(sheetName, index)
		}
		file.SetCellStr(sheetName, fmt.Sprintf("C%d", index), fmt.Sprintf("%d", i+1))
		file.SetCellStr(sheetName, fmt.Sprintf("D%d", index), huxing.BuildingNo)
		file.SetCellStr(sheetName, fmt.Sprintf("E%d", index), huxing.HuxingNo)
		file.SetCellStr(sheetName, fmt.Sprintf("F%d", index), huxing.AreaShow)
		if huxing.ID == declaration.DeclarationHuxingID {
			file.SetCellStr(sheetName, fmt.Sprintf("H%d", index), fmt.Sprintf("(  1  )"))
		} else {
			file.SetCellStr(sheetName, fmt.Sprintf("H%d", index), fmt.Sprintf("(  ×  )"))
		}
		index++
	}
	//合并单元格
	file.MergeCell(sheetName, "A9", fmt.Sprintf("B%d", index-1))
	file.SetCellStr(sheetName, fmt.Sprintf("D%d", index), fmt.Sprintf("%s", declaration.DeclarationAreaShow))
	file.SetCellStr(sheetName, fmt.Sprintf("D%d", index+1), adminName)
	//取消另存为，直接流形式传输
	fileName := fmt.Sprintf("%s-%d-%s.xlsx", contract.ContractNo, declaration.ID, times.ToStr())
	//path := fmt.Sprintf("./excel/%s", fileName)
	//err = file.SaveAs(path)
	//if err != nil {
	//	return nil, "", err
	//}
	//获取当前分期的名称 -- 创建日志
	staging, err := model.FindStagingById(declaration.StagingID)
	logging := model.Logging{
		Username:    adminName,
		StagingName: staging.StagingName,
		Details:     "打印合同单" + declaration.ContractNo + "的申报表",
	}
	logging.Create()
	return file, fileName, nil
}

func DeleteDeclarationById(id uint) error {
	data, err := model.GetDeclarationByID(id)
	if err != nil {
		return errors.New("没有此申报信息")
	}
	if err := model.DeclarationDeleteAble(data); err != nil {
		return err
	}
	// 可以删除 修改时段
	var time model.Time
	time.ID = data.TimeID
	err = time.FindTimeByID()
	if err != nil {
		return err
	}
	time.SelectedNum--
	time.Update()

	err2 := model.DeleteDeclarationById(data.ID)
	if err != nil {
		return err2
	}

	return nil
}

func DeclarationEnterResult(body api.DeclarationEnterResultBody, operator string) error {
	// 查看是否有此申报表
	data, err := model.GetDeclarationByID(body.DeclarationId)
	if err != nil {
		return errors.New("查询申报表异常")
	}

	if body.Status != 0 && body.Status != 1 {
		return errors.New("中签异常 0->不中签 1->中签")
	}
	if data.WinningStatus == model.WinningStatus(body.Status) {
		return errors.New("申报表状态与当前状态相同")
	}
	if body.Status == 0 {
		// 不中 回迁返回面积
		contract, err := model.FindContractById(data.ContractNo)
		if err != nil {
			return errors.New("查询合同异常")
		}
		//获取剩余拆迁面积
		lq_area := (*contract).RemainingHQArea
		//获取户型面积
		huxingAreaDecimal, err := decimal.NewFromString(data.DeclarationArea)
		if err != nil {
			return err
		}
		huxing_area, _ := huxingAreaDecimal.Float64()

		//计算 --加
		remainingHQArea := decimal.NewFromFloat(lq_area)
		huxingArea := decimal.NewFromFloat(huxing_area)
		area, _ := remainingHQArea.Add(huxingArea).Float64()
		//更新合同单数据
		contract.RemainingHQArea = area
		if err := model.UpdateContract(contract); err != nil {
			return err
		}
		//创建面积明细表
		areaDetail := model.AreaDetails{
			ContractNo:         contract.ContractNo,
			OperationalDetails: "未中签 回迁面积已返回",
			OperationalArea:    "+" + data.DeclarationArea + "㎡",
			RemainingArea:      remainingHQArea.Add(huxingArea).String() + "㎡",
		}
		err = areaDetail.Create()
		if err != nil {
			return err
		}
		// 修改申报表为报废状态
		data.ActiveState = false
		err = data.Update()
		if err != nil {
			return errors.New("修改申报状态失败")
		}
		data.WinningStatus = 0
		err = data.Update()
		if err != nil {
			return errors.New("修改申报表失败")
		}
	} else {
		// 中签啦
		// 申报表状态是否正常
		if data.DeclarationStatus != 1 {
			return errors.New("当前申报状态暂未确定，请先修改状态")
		}
		// 修改申报表里面户型最大可选-1 注意轮数是一起的
		huxingId := data.DeclarationHuxingID
		var huxing model.Huxing
		huxing.ID = huxingId
		err := huxing.FindHuxingByID()
		if err != nil {
			return errors.New("查找户型失败")
		}
		huxing.Maximum--
		err = model.UpdateHuxing(huxing)
		if err != nil {
			return errors.New("修改户型最大可选失败")
		}
		// 查找合同
		contract, err := model.FindContractById(data.ContractNo)
		if err != nil {
			return errors.New("查询合同失败")
		}
		// 查找时间
		var time model.Time
		time.ID = data.TimeID
		time.FindTimeByID()
		if body.BuildingOn == "" {
			body.BuildingOn = huxing.BuildingNo
		}

		// 创建摇珠结果表 result
		result := model.Result{
			DeclarationID:         data.ID,
			PublicityStatus:       true,
			BuildingNo:            body.BuildingOn,
			RoomNo:                body.RoomOn,
			Operator:              operator,
			Peoples:               contract.Peoples,
			CardNumber:            contract.CardNumber,
			ContractNo:            contract.ContractNo,
			StagingID:             contract.StagingID,
			PhoneNumber1:          contract.PhoneNumber1,
			PhoneNumber2:          contract.PhoneNumber2,
			DeclarationHuxingID:   huxing.ID,
			DeclarationHuxingNo:   huxing.HuxingNo,
			DeclarationBuildingNo: huxing.BuildingNo,
			DeclarationAreaShow:   huxing.AreaShow,
			DeclarationArea:       huxing.Area,
			TimeID:                time.ID,
			TimeName:              time.Name,
			StagingIdContract:     contract.StagingID,
		}

		err = model.SaveResult(result)
		if err != nil {
			errors.New("创建摇珠结果表失败")
		}

		// 修改申报表 中签状态
		data.WinningStatus = 1
		err = data.Update()
		if err != nil {
			return errors.New("修改申报表失败")
		}

	}

	// 注意创建logging
	staging := data.StagingID
	stagData, err := model.FindStagingById(staging)
	if err != nil {
		return errors.New("查询期数表失败")
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		return errors.New("申报表序列表失败")
	}
	logging := model.Logging{
		Username:    operator,
		StagingName: stagData.StagingName,
		Operation:   "修改申报表中签状态-> 申报表id：" + strconv.Itoa(int(data.ID)) + " 修改中签状态：" + strconv.Itoa(int(body.Status)),
		Details:     string(marshal),
	}
	logging.Create()
	return nil
}

func GetAllDeclarationFuzzy(body api.GetAllDeclarationFuzzyBody) (*model.PaginationQ, error) {
	// 判断staging是否存在
	if !model.ExistsStagingInDeclaration(body.StagingId) {
		return nil, errors.New("申报表没有当前期数的申报")
	}

	data, err := model.GetAllDeclarationFuzzy(body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func UpdateDeclarationActiveStatus(id uint, state bool, operator string) error {
	// 查询
	declaration, err := model.GetDeclarationByID(id)
	contract, err := model.FindContractById(declaration.ContractNo)

	if err != nil {
		return errors.New("查询申报表失败")
	}

	if declaration.ActiveState == state {
		return errors.New("状态一致")
	}

	if err != nil {
		return errors.New("查询合同失败")
	}
	if state == false {
		if declaration.WinningStatus != -1 {
			return errors.New("当前申报是待抽签或者已抽签状态，不可直接作废")
		}
		// 查看 申报状态是否为 1
		if declaration.DeclarationStatus == 1 {
			// 返回面积，创建面积表
			//获取剩余拆迁面积
			lq_area := (*contract).RemainingHQArea
			//获取户型面积
			huxingAreaDecimal, err := decimal.NewFromString(declaration.DeclarationArea)
			if err != nil {
				return err
			}
			huxing_area, _ := huxingAreaDecimal.Float64()
			//如果更改为已确认，则减少合同单的剩余面积，添加面积明细记录

			//计算 --加
			remainingHQArea := decimal.NewFromFloat(lq_area)
			huxingArea := decimal.NewFromFloat(huxing_area)
			area, _ := remainingHQArea.Add(huxingArea).Float64()
			//更新合同单数据
			contract.RemainingHQArea = area
			if err := model.UpdateContract(contract); err != nil {
				return err
			} //创建面积明细表
			areaDetail := model.AreaDetails{
				ContractNo:         contract.ContractNo,
				OperationalDetails: "申报被作废",
				OperationalArea:    "+" + declaration.DeclarationArea + "㎡",
				RemainingArea:      remainingHQArea.Add(huxingArea).String() + "㎡",
			}
			err = areaDetail.Create()
			if err != nil {
				return err
			}

		}

		declaration.ActiveState = false
		declaration.DeclarationStatus = model.DeclarationStatus(0)
		declaration.Operator = operator
		err = declaration.Update()
	}

	if state == true {
		// 有效
		declaration.ActiveState = true
		declaration.Operator = operator
		err = declaration.Update()
		if err != nil {
			return err
		}
	}

	data, err := json.Marshal(map[string]interface{}{
		"declarationID": declaration.ID,
		"activeStatus":  declaration.ActiveState,
	})
	staging, err := model.FindStagingById(declaration.StagingID)
	if err != nil {
		return err
	}
	logging := model.Logging{
		Username:    operator,
		StagingName: staging.StagingName,
		Operation:   "更改合同" + declaration.ContractNo + "申报表" + strconv.Itoa(int(declaration.ID)) + "是否有效" + strconv.FormatBool(state),
		Details:     string(data),
	}
	logging.Create()
	return nil

}

func UpdateDeclarationTrustee(declaration api.UpdateDeclaration, operator string, ok bool) error {

	// 首先id必须先指定
	if declaration.DeclarationID == 0 {
		return errors.BadError("请先指定申报id")
	}
	// 判断是不是管理员，管理员可以随意修改，用户只能修改自己的申报表
	if ok {
		admin, err := model.GetAdminInfo(operator)
		if err != nil {
			if err == nil {
				return errors.BadError("不存在该管理员")
			}
			return err
		}
		operator = admin.AdminSignname
		err = model.UpdateDeclarationByTrusteeByContractNo(declaration)
		if err != nil {
			return err
		}
	}

	if !ok {
		//用户
		user, err := model.GetUserInfo(operator)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.BadError("不存在该用户")
			}
		}
		// 先判断申报id属不属于这个客户
		isBelong, err := model.IsDeclarationBelongToClientByPhone(declaration.DeclarationID, user.PhoneNumber)
		if err != nil {
			return err
		}
		if isBelong {
			// 属于
			err := model.UpdateDeclarationByTrusteeByContractNo(declaration)
			if err != nil {
				return err
			}
		} else {
			// 不属于
			return errors.New("该申报不属于你，填入失败")
		}
	}
	body, err := model.GetDeclarationByID(declaration.DeclarationID)
	if err != nil {
		return err
	}
	data, err := json.Marshal(declaration)
	if err != nil {
		return err
	}
	logging := model.Logging{
		Username:    operator,
		StagingName: strconv.Itoa(int(body.StagingID)),
		Operation:   "申报表" + strconv.Itoa(int(declaration.DeclarationID)) + "添加受托人信息",
		Details:     string(data),
	}
	logging.Create()
	return nil
}
