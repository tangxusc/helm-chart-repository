package controller

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"path/filepath"
	"repository/event"
	"repository/httpserver"
	"repository/repository/domain"
	"strconv"
	"time"
)

func init() {
	httpserver.AddRegister(func(app *iris.Application) {
		app.Post("/chart", createHandler)
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
	event.Send(&FileUploaded{
		File:        &file,
		ChartName:   chart.Name,
		FileName:    header.Filename,
		NewFileName: filename,
	})
	event.Send(&ChartCreated{
		ChartVersion: chart,
		FileName:     filename,
	})
}

type FileUploaded struct {
	File        *multipart.File
	ChartName   string
	FileName    string
	NewFileName string
}

type ChartCreated struct {
	*domain.ChartVersion
	FileName string
}
