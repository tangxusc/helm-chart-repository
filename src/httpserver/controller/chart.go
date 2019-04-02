package controller

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"repository/event"
	"repository/httpserver"
	"repository/repository/domain"
	"repository/repository/entry"
	"strconv"
	"time"
)

func init() {
	httpserver.AddRegister(func(app *iris.Application) {
		app.Post("/chart", createHandler)
		app.Delete("/chart/{chartName:string}/{version:string}", deleteHandler)
		app.Get("/chart/{chartName:string}/{version:string}/info", infoHandler)
		app.Get("/chart/{chartName:string}/{version:string}/download", infoHandler)
		app.Get("/chart/{chartName:string}/list", listHandler)
	})
}

func listHandler(ctx context.Context) {
	chartName := ctx.Params().Get("chartName")

	_, e := ctx.JSON(entry.MustLoadChartVersionByName(chartName))
	if e != nil {
		panic(e)
	}
}

func infoHandler(ctx context.Context) {
	chartName := ctx.Params().Get("chartName")
	version := ctx.Params().Get("version")

	_, e := ctx.JSON(entry.MustLoadChartVersion(chartName, version))
	if e != nil {
		panic(e)
	}
}

func deleteHandler(ctx context.Context) {
	chartName := ctx.Params().Get("chartName")
	version := ctx.Params().Get("version")

	event.Send(&domain.ChartDeleted{
		ChartName: chartName,
		Version:   version,
	})
}

func createHandler(ctx context.Context) {
	chart := &domain.ChartVersion{}
	err := ctx.ReadJSON(chart)
	if err != nil {
		logrus.Errorf("read create chart form error, %s", err.Error())
		panic(err)
	}
	file, header, err := ctx.FormFile("tarFile")
	if err != nil {
		logrus.Errorf("read chart file error, %s", err.Error())
		panic(err)
	}
	filename := strconv.Itoa(time.Now().Nanosecond()) + filepath.Ext(header.Filename)
	event.Send(&domain.FileUploaded{
		File:        &file,
		ChartName:   chart.Name,
		FileName:    header.Filename,
		NewFileName: filename,
	})
	event.Send(&domain.ChartCreated{
		ChartVersion: chart,
		FileName:     filename,
	})
}
