package logging

import (
	"github.com/sirupsen/logrus"
)

var (
	Logger *logrus.Logger
)

// 初始化日志配置
func Setup() {
	Logger = logrus.New()
}

func Debug(v ...interface{}) {
	Logger.Debug(v)
}

func Debugln(v ...interface{}) {
	Logger.Debugln(v)
}

func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

func Info(v ...interface{}) {
	Logger.Info(v)
}

func Infoln(v ...interface{}) {
	Logger.Infoln(v)
}

func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

func Warn(v ...interface{}) {
	Logger.Warn(v)
}

func Warnln(v ...interface{}) {
	Logger.Warnln(v)
}

func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

func Error(v ...interface{}) {
	Logger.Error(v)
}

func Errorln(v ...interface{}) {
	Logger.Errorln(v)
}

func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

func Fatal(v ...interface{}) {
	Logger.Fatal(v)
}

func Fatalln(v ...interface{}) {
	Logger.Fatalln(v)
}

func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}
