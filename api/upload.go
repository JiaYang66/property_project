package api

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"path"
	"relocate/util/errors"
	"relocate/util/files"
	"relocate/util/logging"
	"strings"
)

const (
	ACCESSKEYID     = "***"
	ACCESSSECRET    = "***"
	ENDPOINT        = "***"
	BUCKEDNAME      = "***"
	OBJECTKEYIDCARD = "***"
	ROOT            = "***"
)

// 支持的上传文件后缀名
var defaultSupportExtNames = []string{".jpg", ".jpeg", ".png", ".ico", ".svg", ".bmp", ".gif"}

func isSupportExtNames(extName string, supportExtNames []string) bool {
	for i := 0; i < len(supportExtNames); i++ {
		if supportExtNames[i] == extName {
			return true
		}
	}
	return false
}

func Upload(c *gin.Context, name string, sizeM int, supportExtNames ...[]string) (*multipart.FileHeader, error) {
	fileHeader, err := c.FormFile(name)
	if err != nil {
		return nil, err
	}
	extname := strings.ToLower(path.Ext(fileHeader.Filename))
	var supportedFileTypes []string
	if len(supportExtNames) > 0 {
		supportedFileTypes = supportExtNames[0]
	} else {
		supportedFileTypes = defaultSupportExtNames
	}
	if !isSupportExtNames(extname, supportedFileTypes) {
		return nil, errors.BadError(fmt.Sprintf("不支持的文件类型,请上传%s类型文件", supportedFileTypes))
	}
	if fileHeader.Size > int64(sizeM*1024*1024) {
		return nil, errors.BadError(fmt.Sprintf("文件大小不能超过%dM", sizeM))
	}
	return fileHeader, nil
}

func ReadFile(fileHeader *multipart.FileHeader) ([]byte, error) {
	f, err := fileHeader.Open()
	defer f.Close()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

const defaultFolder = "./upload/"

//path为路径+文件名
func SaveFile(c *gin.Context, fileHeader *multipart.FileHeader, path ...string) error {
	var disk string
	if len(path) > 0 {
		disk = path[0]
	} else {
		disk = defaultFolder + fileHeader.Filename
	}
	if err := files.MkdirAll(disk); err != nil {
		return err
	}
	return c.SaveUploadedFile(fileHeader, disk)
}

//上传图片
func UploadImg(file, suffix, phoneNumber, name string) (string, error) {
	flag := true
	for _, val := range defaultSupportExtNames {
		if val == suffix {
			flag = false
			break
		}
	}
	if flag {
		return "", errors.BadError(fmt.Sprintf("不支持的文件类型"))
	}

	client, err := oss.New(ENDPOINT, ACCESSKEYID, ACCESSSECRET)
	if err != nil {
		return "", errors.BadError(fmt.Sprintf("上传图片连接失败"))
	}

	// 获取存储空间。
	bucket, err := client.Bucket(BUCKEDNAME)
	if err != nil {
		return "", errors.BadError(fmt.Sprintf("出现异常"))
	}

	data, err := base64.StdEncoding.DecodeString(file)
	if err != nil {
		return "", errors.BadError(fmt.Sprintf("解码错误"))
	}
	image := OBJECTKEYIDCARD + phoneNumber + "/" + name + suffix
	err = bucket.PutObject(image, bytes.NewBuffer(data))
	if err != nil {
		return "", errors.BadError(fmt.Sprintf("上传失败"))
	}
	logging.Info("上传成功")
	return ROOT + image, nil
}
