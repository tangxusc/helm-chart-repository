package controller

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"repository/config"
	"repository/domain"
	"repository/event"
	"repository/httpserver"
	"repository/repository/entry"
	"repository/repository/index"
	"strings"
)

func init() {
	httpserver.AddRegister(func(app *iris.Application) {
		app.Post("/chart", createHandler)
		app.Delete("/chart/{chartName:string}/{version:string}", deleteHandler)
		app.Get("/chart/{chartName:string}/{version:string}/info", infoHandler)
		app.Get("/chart/{chartName:string}/{version:string}/download", downloadHandler)
		app.Get("/chart/{chartName:string}", listChartVersionsHandler)
		app.Get("/chart", listChartHandler)
	})
}

func downloadHandler(ctx context.Context) {
	chartName := ctx.Params().Get("chartName")
	version := ctx.Params().Get("version")
	filename := fmt.Sprintf("%s-%s", chartName, version)
	resultPath := filepath.Join(config.Config.DataDir, chartName, filename)
	err := ctx.SendFile(resultPath, filename)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"chartName": chartName,
			"version":   version,
			"path":      resultPath,
		}).Error("download chart error")
	}
}

func listChartHandler(ctx context.Context) {
	pageSize, _ := ctx.URLParamInt("pageSize")
	if pageSize <= 0 {
		pageSize = 20
	}
	pageNum, _ := ctx.URLParamInt("pageNum")
	if pageNum <= 0 {
		pageNum = 1
	}
	charts := getChartInfos()
	total := len(charts)
	startIndex := pageSize * (pageNum - 1)
	if startIndex > total {
		startIndex = total
	}
	endIndex := startIndex + pageSize
	if endIndex > total {
		endIndex = total
	}
	pageTotal := total / pageSize
	if total%pageSize > 0 {
		pageTotal = pageTotal + 1
	}
	result := charts[startIndex:endIndex]
	ctx.JSON(&ChartList{
		pageTotal,
		pageNum,
		result,
	})
}

func getChartInfos() []*ChartInfo {
	charts := index.GetCharts()
	result := make([]*ChartInfo, len(charts))
	var i = 0
	for _, value := range charts {
		item := value[0]
		result[i] = &ChartInfo{
			item.Name,
			item.Description,
			item.Engine,
			item.Home,
			item.Icon,
			strings.Join(item.URLs, ","),
		}
		i++
	}
	return result
}

type ChartInfo struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Engine      string `json:"engine,omitempty"`
	Home        string `json:"home,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Version     string `json:"version,omitempty"`
}
type ChartList struct {
	PageTotal int          `json:"pageTotal"`
	PageNum   int          `json:"pageNum"`
	List      []*ChartInfo `json:"list"`
}

func listChartVersionsHandler(ctx context.Context) {
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
	ctx.JSON("success")
}

func createHandler(ctx context.Context) {
	chart := &domain.ChartVersion{}
	err := ctx.ReadForm(chart)
	if err != nil {
		logrus.Errorf("read create chart form error, %s", err.Error())
		panic(err)
	}
	file, header, err := ctx.FormFile("tarFile")
	if err != nil {
		logrus.Errorf("read chart file error, %s", err.Error())
		panic(err)
	}
	filename := fmt.Sprintf("%s-%s", chart.Name, chart.Version)
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
	ctx.JSON("success")
}
