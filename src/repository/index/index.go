package index

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"repository/config"
	"repository/event"
	"repository/httpserver/controller"
	"repository/repository/domain"
	"repository/repository/entry"
	"time"
)

var indexFile *domain.IndexFile

func init() {
	event.Subscribe(100, event.Handlers{
		"*controller.ChartCreated": handlerChartCreated,
	})
}

/**
处理 chart创建事件
*/
func handlerChartCreated(event interface{}) {
	created := event.(*controller.ChartCreated)
	indexFile.Generated = time.Now()
	//TODO:加上下载前缀
	created.URLs = []string{created.FileName}

	versions, ok := indexFile.Entries[created.Name]
	if !ok {
		versions = make([]*domain.ChartVersion, 0)
	}
	versions = append(versions, created.ChartVersion)
	indexFile.Entries[created.Name] = versions
	err := WriteIndexFile()
	if err != nil {
		logrus.Errorf("handler ChartCreated event error,%s", err.Error())
	}
}

func InitIndexFile() {
	_, err := os.Open(getIndexFilePath())
	if err != nil {
		notExist := os.IsNotExist(err)
		if !notExist {
			logrus.Fatalf("open index.yaml error,%s", err.Error())
		}
	}

	indexFile, err = LoadIndexFile()
	if err != nil {
		logrus.Fatalf("get index.yaml error,%s", err.Error())
	}
	err = WriteIndexFile()
	if err != nil {
		logrus.Fatalf("write index.yaml error,%s", err.Error())
	}
	event.Send(&ChartUpdated{
		ChartTotal: len(indexFile.Entries),
	})
}

type ChartUpdated struct {
	ChartTotal int
}

func WriteIndexFile() error {
	out, err := yaml.Marshal(indexFile)
	if err != nil {
		logrus.Errorf("marshal indexFile error,%s", err.Error())
		return err
	}
	path := getIndexFilePath()
	err = ioutil.WriteFile(path, out, os.ModePerm)
	if err != nil {
		logrus.Errorf("write index.yaml error,%s", err.Error())
	}
	return err
}

func getIndexFilePath() string {
	return filepath.Join(config.Config.DataDir, "index.yaml")
}

func LoadIndexFile() (*domain.IndexFile, error) {
	matches, e := filepath.Glob(filepath.Join(config.Config.DataDir, "**/"+config.Config.EntryFileName))
	if e != nil {
		logrus.Error(e.Error())
		return nil, e
	}
	result := &domain.IndexFile{
		APIVersion: domain.APIVersionV1,
		Generated:  time.Now(),
		Entries:    make(map[string]domain.ChartVersions),
	}
	for _, value := range matches {
		charts, e := entry.LoadChartVersionsByFile(value)
		if e != nil {
			logrus.Error(e.Error())
			continue
		}
		if len(charts) > 0 {
			result.Entries[charts[0].Name] = domain.ChartVersions(charts)
		}
	}

	return result, nil
}
