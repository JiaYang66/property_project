package user_service

import (
	"encoding/json"
	"relocate/api"
	"relocate/model"
	"relocate/util/errors"
	"relocate/util/gredis"
	"relocate/util/logging"
	"relocate/util/random"
	"relocate/util/sign"
	"relocate/util/times"
	"strconv"
)

func CreateUser(username, password string) (err error) {
	user, err := model.GetUserInfo(username)
	// 稀释盐
	salt := "ABCDEF"
	if user != nil {
		return errors.BadError("手机号码已被使用")
	}
	user = &model.User{
		PhoneNumber: username,
		Password:    sign.EncodeMD5(password + salt),
		Salt:        salt,
		UserStatus:  model.NotVerifiedUser,
		CreatedAt:   times.JsonTime{},
	}
	return user.Create()
}

func UpdateUser(idCardType int, phoneNumber, name, idCard, imagesA, imagesB, suffixA, suffixB string) error {
	idtype := model.MainlandId
	if idCardType == 1 {
		idtype = model.HongKongId
	} else if idCardType == 2 {
		idtype = model.HuZhaoId
	}
	formattedA, err := api.UploadImg(imagesA, suffixA, phoneNumber, "positive_images")
	if err != nil {
		logging.Info(err)
		return errors.BadError("正面照上传出错")
	}
	formattedB, err := api.UploadImg(imagesB, suffixB, phoneNumber, "negative_images")
	if err != nil {
		logging.Info(err)
		return errors.BadError("反面照上传出错")
	}
	user := &model.User{
		IdNumberType:  idtype,
		PhoneNumber:   phoneNumber,
		Name:          name,
		IdNumber:      idCard,
		PositiveImage: formattedA,
		NegativeImage: formattedB,
	}
	return user.UpdateUser()
}

func UpdateStatus(phoneNumber string) error {
	user := &model.User{
		PhoneNumber: phoneNumber,
		UserStatus:  model.AutomaticMatchingUser,
	}
	return user.UpdateUser()
}

func UserForgetSendCode(phone string) error {
	_, err := model.GetUserInfo(phone)
	if err != nil {
		return errors.BadError("该用户不存在")
	}

	code := random.Code(6)
	if !random.SendCode(phone, code) {
		return errors.BadError("发送验证码错误")
	}

	user := make(map[string]string)
	user["user-forget-"+phone] = code
	gredis.Set(user, 200)
	return nil
}

func UpdateUserPassword(phoneNumber string, password string, code string) error {
	_, err := model.GetUserInfo(phoneNumber)
	if err != nil {
		return errors.BadError("该用户不存在")
	}

	realyCode, _ := gredis.Get("user-forget-" + phoneNumber)
	if realyCode != code { //用户没有获取验证码或者验证码过期
		return errors.BadError("验证码错误")
	}
	// 验证码正确
	user := model.User{
		PhoneNumber: phoneNumber,
		Password:    password,
	}

	user.UpdateUser()
	return nil
}

func QueryContractByAccountId(IdNumber string) ([]model.Contract, error) {
	// 查看此id是否存在
	if !model.ExistsUserByIdNumber(IdNumber) {
		return nil, errors.BadError("不存在此用户")
	}

	// 查看是否有此合同单
	if !model.ExistsContractByCardNumber(IdNumber) {
		return nil, errors.BadError("改用户没有合同单")
	}
	// 一个人可能有多个合同单
	contracts, err := model.FindContractByCardNumber(IdNumber)
	if err != nil {
		return nil, errors.BadError(err.Error())
	}

	return contracts, nil
}

func QueryResultByAccountId(IdNumber string) ([]model.ResultResultByAccountId, error) {
	// 查看此id是否存在
	if !model.ExistsUserByIdNumber(IdNumber) {
		return nil, errors.BadError("不存在此用户")
	}

	// 查看是否有此瑶珠结果
	if !model.ExistResultByCardNumber(IdNumber) {
		return nil, errors.BadError("改用户没有摇珠结果")
	}
	// 一个人可能有多个摇珠结果
	results, err := model.FindResultByCardNumber(IdNumber)
	if err != nil {
		return nil, errors.BadError(err.Error())
	}

	return results, nil
}

func GenerateToken(adminName, password string) (interface{}, error) {
	user, err := model.FindUserByPhone(adminName)

	if err != nil {
		return "", errors.BadError("用户账号不存在")
	}
	//校验密码
	if user.Password != sign.EncodeMD5(password+user.Salt) {
		return "", errors.BadError("密码错误")
	}
	//生成jwt-token
	token, err := sign.GenerateToken(user.PhoneNumber, user.PhoneNumber, sign.UserClaimsType)
	if err != nil {
		return "", err
	}

	return map[string]interface{}{"token": token, "signName": user.Name}, nil
}

func UserQueryByStatusOrFilterName(filterStatus string, filterName string, page uint, pageSize uint) (data *model.PaginationQ, err error) {

	if filterStatus != "" {
		// status 是否在范围内
		status, err := strconv.ParseInt(filterStatus, 10, 8)
		if err != nil {
			return nil, errors.BadError("status错误")
		}

		_, ok := model.AllUserStatus[model.UserStatus(status)]
		if !ok {
			return nil, errors.BadError("status不存在")
		}
	}

	data, err = model.UserQueryByStatusOrFilterName(filterStatus, filterName, page, pageSize)
	if err != nil {
		return nil, err
	}
	return data, nil

}

func UpdateUserStatus(phone string, statue bool) error {
	if !model.ExistsUserByPhoneNumber(phone) {
		return errors.New("不存在此用户")
	}
	var states model.UserStatus
	if statue == false {
		// 没传或者不通过
		var user = model.User{
			PhoneNumber: phone,
			UserStatus:  model.AutomaticMatchingUser,
		}
		states = model.AutomaticMatchingUser
		err := user.Update()
		if err != nil {
			return err
		}
		return nil
	} else if statue == true {
		// 通过
		var user = model.User{
			PhoneNumber: phone,
			UserStatus:  model.VerifiedUser,
		}
		states = model.VerifiedUser
		err := user.Update()
		if err != nil {
			return err
		}
		return nil
	}
	user, err := model.FindUserByPhone(phone)
	if err != nil {
		return err
	}
	marshal, err := json.Marshal(user)
	if err != nil {
		return err
	}
	logging := model.Logging{
		Username:    phone,
		StagingName: user.Name,
		Operation:   "修改用户状态" + user.PhoneNumber + "->" + states.String(),
		Details:     string(marshal),
	}
	logging.Create()
	return nil

}
