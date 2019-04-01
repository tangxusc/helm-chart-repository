package main

import (
	"github.com/sirupsen/logrus"
	"repository/config"
	"repository/event"
	"repository/httpserver"
	_ "repository/httpserver/controller"
	_ "repository/httpserver/metrics"
	_ "repository/repository/entry"
	"repository/repository/index"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{})
	//加载配置
	config.LoadConfig()
	//注册http server
	httpserver.InitServer()
	//初始化index.yaml
	index.InitIndexFile()
	go event.Listen()
	//启动http服务器
	httpserver.Start()
}
