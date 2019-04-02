package controller

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"repository/httpserver"
)

func init() {
	httpserver.AddRegister(func(app *iris.Application) {
		app.Get("/health", func(ctx context.Context) {
			ctx.Text("OK")
		})
	})
}
