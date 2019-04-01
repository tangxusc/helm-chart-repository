package main

import (
	"github.com/sirupsen/logrus"
	"repository/config"
	"repository/event"
	"repository/httpserver/controller"
	"repository/repository/domain"
	"repository/repository/entry"
	"repository/repository/index"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{})

	//加载配置
	config.LoadConfig()
	index.InitIndexFile()
	go index.Listen()
	go entry.Listen()

	chart := &domain.ChartVersion{}
	chart.Name = "test"
	chart.Version = "0.1"
	event.Send(&controller.ChartCreated{
		ChartVersion: chart,
		FileName:     "test.tar.gz",
	})
	chart2 := &domain.ChartVersion{}
	chart2.Name = "test"
	chart2.Version = "0.2"
	event.Send(&controller.ChartCreated{
		ChartVersion: chart2,
		FileName:     "test.tar.gz",
	})

	time.Sleep(2 * time.Second)
}