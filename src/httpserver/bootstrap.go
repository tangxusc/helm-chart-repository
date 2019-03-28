package httpserver

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
	"repository/config"
)

var HttpServer *iris.Application

type Register func(application *iris.Application)

var Registers = make([]Register, 0)

func Start() {
	err := HttpServer.Run(iris.Addr(fmt.Sprintf(":%s", config.Config.ServerPort)))
	if err != nil {
		logrus.Error(err.Error())
	}
}

func InitServer() {
	HttpServer = iris.Default()
	HttpServer.Configure(iris.WithTimeFormat("2006-01-02 15:04:05"))
	fmt.Println(len(Registers))
	for _, value := range Registers {
		value(HttpServer)
	}
}

func AddRegister(router ...Register) {
	Registers = append(Registers, router...)
}
