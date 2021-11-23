package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"path"
	"relocate/util/app"
	"relocate/util/conf"
	"relocate/util/logging"
	"time"
)

type respLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w respLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w respLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

var (
	logFilePath = "./log"
	logFileName = "system.log"
)

// 日志记录到文件
func LoggerToFile() gin.HandlerFunc {
	// 设置日志级别
	if conf.Data.Server.RunMode == gin.DebugMode {
		logging.Logger.SetLevel(logrus.DebugLevel)
	} else {
		logging.Logger.SetLevel(logrus.InfoLevel)
	}
	// 设置是否输出到文件
	if conf.Data.Server.Log {
		// 日志文件
		fileName := path.Join(logFilePath, logFileName)
		// 设置日志切割
		logWriter, _ := rotatelogs.New(
			// 分割后的文件名称
			fileName+".%Y%m%d.log",
			// 生成软链，指向最新日志文件
			rotatelogs.WithLinkName(fileName),
			// 设置最大保存时间(7天)
			rotatelogs.WithMaxAge(7*24*time.Hour),
			// 设置日志切割时间间隔(1天)
			rotatelogs.WithRotationTime(24*time.Hour),
		)
		writeMap := lfshook.WriterMap{
			logrus.InfoLevel:  logWriter,
			logrus.FatalLevel: logWriter,
			logrus.DebugLevel: logWriter,
			logrus.WarnLevel:  logWriter,
			logrus.ErrorLevel: logWriter,
			logrus.PanicLevel: logWriter,
		}
		// 新增钩子
		logging.Logger.AddHook(lfshook.NewHook(writeMap, &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		}))
	}
	return func(c *gin.Context) {
		bodyLogWriter := &respLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 状态码
		statusCode := c.Writer.Status()
		// 日志格式
		entry := logging.Logger.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"uri":        c.Request.RequestURI,
			"ip":         c.ClientIP(),
			"statusCode": statusCode,
			"costTime":   endTime.Sub(startTime),
		})
		// 返回数据
		resp := app.Response{}
		if bodyLogWriter.body.String() != "" {
			_ = json.Unmarshal(bodyLogWriter.body.Bytes(), &resp)
		}
		var result string
		if resp.Code != 0 {
			result = fmt.Sprintf("返回数据: %s", bodyLogWriter.body.String())
		}
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			if statusCode > 499 {
				entry.Error(result)
			} else if statusCode > 399 {
				entry.Warn(result)
			} else {
				entry.Info(result)
			}
		}
	}
}

// 日志记录到 DB
func LoggerToDB() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 ES
func LoggerToES() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 MQ
func LoggerToMQ() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
