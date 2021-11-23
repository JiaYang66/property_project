package random

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"strings"
	"time"
)

const (
	ACCESSKEYID  = "***"
	ACCESSSECRET = "***"
	SIGNNAME     = "***"
	TMPLATECODE  = "***"
	SUCCESS      = "***"
)

//生成指定长度的随机字符串
func String(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//生成指定长度的随机数字
func Code(length int) string {
	var container string
	for i := 0; i < length; i++ {
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		container += fmt.Sprintf("%01v", rnd.Int31n(10))
	}
	return container
}

//生成UUID
//is是否去除- 默认去除
func UUID(is ...bool) string {
	if len(is) > 0 && is[0] {
		return fmt.Sprintf("%s", uuid.Must(uuid.NewV4(), nil))
	}
	return strings.Replace(fmt.Sprintf("%s", uuid.Must(uuid.NewV4(), nil)), "-", "", -1)
}

//生成指定范围内随机值
//[min,max)
func RangeNum(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(max-min) + min
	return randNum
}

// 发送短信
func SendCode(phoneNumber, code string) bool {
	response := SendMessage(ACCESSKEYID, ACCESSSECRET, SIGNNAME, phoneNumber, TMPLATECODE, code)
	fmt.Println(response.Code)
	return response.Code == SUCCESS
}

func SendMessage(accessKeyId, accessSecret, signName, phoneNum, templateCode string, code string) *dysmsapi.SendSmsResponse {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", accessKeyId, accessSecret)

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = phoneNum
	request.SignName = signName
	request.TemplateCode = templateCode
	params := make(map[string]interface{})
	params["code"] = code

	jsonData, err := json.Marshal(&params)
	if err != nil {
		panic(err)
	}
	request.TemplateParam = string(jsonData)

	response, err := client.SendSms(request)
	if err != nil {
		fmt.Print(err.Error())
	}

	return response
}
