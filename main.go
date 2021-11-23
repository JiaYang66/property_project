package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"relocate/model"
	"relocate/router"
	"relocate/util/conf"
	"relocate/util/gredis"
	"relocate/util/logging"
	"relocate/util/times"
	"relocate/util/validator"
	"time"
)

func init() {
	logging.Setup()
	conf.Setup()
	model.Setup()
	gredis.Setup()
	validator.Setup()
}

// @title relocate
// @version 1.0
// @description 回迁平台
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	logging.Info("设置服务器时区为东八区")
	timeLocal := time.FixedZone("CST", 8*3600)
	time.Local = timeLocal
	logging.Info("当前时间:", times.ToStr())
	gin.SetMode(conf.Data.Server.RunMode)
	httpPort := fmt.Sprintf(":%d", conf.Data.Server.HttpPort)
	server := &http.Server{
		Addr:           httpPort,
		Handler:        router.InitRouter(),
		ReadTimeout:    conf.Data.Server.ReadTimeout,
		WriteTimeout:   conf.Data.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	logging.Info("开始监听服务", httpPort)
	if err := server.ListenAndServe(); err != nil {
		logging.Fatal("启动服务失败", err)
	}
}
