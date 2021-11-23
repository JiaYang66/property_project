package validator

import (
	chinese "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"regexp"
	"strings"
)

var (
	Validate *validator.Validate
	trans    ut.Translator
)

func Setup() {
	// 中文翻译
	zh := chinese.New()
	uni := ut.New(zh, zh)
	trans, _ = uni.GetTranslator("zh")
	Validate = validator.New()
	// 验证器注册翻译器
	zhTranslations.RegisterDefaultTranslations(Validate, trans)
	// 自定义验证方法
	Validate.RegisterValidation("checkUsername", checkUsername)
	// 根据自定义的标记注册翻译
	Validate.RegisterTranslation("checkUsername", trans, func(ut ut.Translator) error {
		return ut.Add("checkUsername", "{0}必须是由字母开头的4-16位字母和数字组成的字符串", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("checkUsername", fe.Field())
		return t
	})

	// 自定义验证方法
	Validate.RegisterValidation("checkMobile", checkMobile)
	// 根据自定义的标记注册翻译
	Validate.RegisterTranslation("checkMobile", trans, func(ut ut.Translator) error {
		return ut.Add("checkMobile", "请输入正确的手机号码", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("checkMobile", fe.Field())
		return t
	})

	Validate.RegisterValidation("checkCardNumber", checkCardNumber)
	Validate.RegisterTranslation("checkCardNumber", trans, func(ut ut.Translator) error {
		return ut.Add("checkCardNumber", "请输入正确的身份证号码", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("checkCardNumber", fe.Field())
		return t
	})
}

func Translate(errs validator.ValidationErrors) string {
	var errList []string
	for _, e := range errs {
		errList = append(errList, e.Translate(trans))
	}
	return strings.Join(errList, "|")
}

func checkUsername(fl validator.FieldLevel) bool {
	return VerifyUsernameFormat(fl.Field().String())
}

func VerifyUsernameFormat(username string) bool {
	if ok, _ := regexp.MatchString(`^[a-zA-Z]{1}[a-zA-Z0-9]{3,15}$`, username); !ok {
		return false
	}
	return true
}
func checkMobile(fl validator.FieldLevel) bool {
	return VerifyMobileFormat(fl.Field().String())
}

func VerifyMobileFormat(mobile string) bool {
	if ok, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`,
		mobile); !ok {
		return false
	}
	return true
}

func CheckCardNumber(idNumberType int, cardNumber string) bool {
	match := `^(\d{17})([0-9]|X|x)$` //大陆18位身份证
	if idNumberType == 1 {
		match = `^[A-Z]\d{6}\([\dA]\)$`
	} else if idNumberType == 2 {
		match = `^[a-zA-Z0-9]{3,21}$`
	}
	if ok, _ := regexp.MatchString(match, cardNumber); !ok {
		return false
	}
	return true
}

func checkCardNumber(fl validator.FieldLevel) bool {
	return VerifyCardNumberFormat(fl.Field().String())
}

func VerifyCardNumberFormat(cardNumber string) bool {
	if ok, _ := regexp.MatchString(`^(\d{17})([0-9]|X|x)$`, cardNumber); ok {
		return true
	}
	if ok, _ := regexp.MatchString(`^[A-Z]\d{6}\([\dA]\)$`, cardNumber); ok {
		return true
	}
	if ok, _ := regexp.MatchString(`^[a-zA-Z0-9]{3,21}$`, cardNumber); ok {
		return true
	}
	return false
}

func CheckManyCardNumber(cardNumbers, s string) bool {
	split := strings.Split(cardNumbers, s)
	for _, val := range split {
		if len(val) == 18 {
			if ok, _ := regexp.MatchString(`^(\d{17})([0-9]|X|x)$`, val); !ok {
				return false
			}
		} else if len(val) == 9 {
			if ok, _ := regexp.MatchString(`^[a-zA-Z0-9]{3,21}$`, val); !ok {
				return false
			}
		} else if len(val) == 10 {
			ok1, _ := regexp.MatchString(`^[A-Z]\d{6}（[\dA]）$`, val)
			ok2, _ := regexp.MatchString(`^[A-Z]\d{6}\([\dA]\)$`, val)
			if !ok1 && !ok2 {
				return false
			}
		} else {
			return false
		}
	}
	return true
}
