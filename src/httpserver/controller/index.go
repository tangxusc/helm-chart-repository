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

func indexHandler(ctx context.Context) {
	err := ctx.SendFile(config.Config.DataDir+"index.yaml", "chart index")
	if err != nil {
		logrus.Error(err.Error())
	}
}
