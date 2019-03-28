package main

import (
	"github.com/sirupsen/logrus"
	"repository/config"
	"repository/httpserver"
	_ "repository/httpserver/controller"
	"repository/httpserver/metrics"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{})
	//加载配置
	config.LoadConfig()
	httpserver.InitServer()
	go metrics.Listen()
	//扫描存储库
	//启动http服务器
	httpserver.Start()
}
