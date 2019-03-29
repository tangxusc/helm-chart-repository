package controller

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/sirupsen/logrus"
	"repository/config"
	"repository/httpserver"
)

func init() {
	httpserver.AddRegister(func(app *iris.Application) {
		app.Get("/index.yaml", indexHandler)
	})
}

/**
TODO:直接发送文件还是序列化repository缓存的IndexFile好一些呢?
*/
func indexHandler(ctx context.Context) {
	logrus.Debugf("DataDir is :%s", config.Config.DataDir)
	filename := "index.yaml"
	err := ctx.SendFile(config.Config.DataDir+filename, filename)
	if err != nil {
		logrus.Error(err.Error())
	}
}
