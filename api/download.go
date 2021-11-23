package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"regexp"
	"relocate/util/errors"
	"relocate/util/files"
)

var imgRegexp = `(.*)\.(jpg|bmp|gif|ico|pcx|jpeg|tif|png|raw|tga)$`

func Download(c *gin.Context, path string) error {
	if !files.IsExist(path) {
		return errors.BadError("文件不存在")
	}
	if !regexp.MustCompile(imgRegexp).
		MatchString(filepath.Ext(path)) {
		//如果不是图片格式，则设置为下载服务
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(path)))
		c.Writer.Header().Add("Content-Type", "application/octet-stream")
	}
	c.File(path)
	return nil
}
