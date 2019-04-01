package main

import (
	"github.com/sirupsen/logrus"
	"repository/config"
	"repository/httpserver"
	_ "repository/httpserver/controller"
	"repository/httpserver/metrics"
	"repository/repository/entry"
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
	go index.Listen()
	go entry.Listen()
	//监听metrics
	go metrics.Listen()
	//启动http服务器
	httpserver.Start()
}
