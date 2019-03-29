package repository

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"repository/config"
	"repository/event"
	"repository/httpserver/controller"
	"repository/repository/domain"
	"time"
)

var indexFile *domain.IndexFile
var eventChan = make(chan interface{}, 1000)

func init() {
	event.RegisterChannel(eventChan)
}

func Listen() {
	for {
		ok := true
		var evt interface{}
		select {
		case evt, ok = <-eventChan:
			logrus.WithFields(logrus.Fields{
				"event": evt,
				"ok":    ok,
			}).Debug("handler Event")
		}
		if !ok {
			break
		}
		handlerEvent(evt)
	}
}

func handlerEvent(event interface{}) {
	switch event.(type) {
	case *controller.ChartCreated:
		created := event.(*controller.ChartCreated)
		handlerChartCreated(created)
	case *controller.FileUploaded:
		upload := event.(*controller.FileUploaded)
		handlerFileUploaded(upload)
	}
}

func handlerFileUploaded(uploaded *controller.FileUploaded) {
	err := checkDirExist(uploaded)
	if err != nil {
		panic(err)
	}
	join := filepath.Join(config.Config.DataDir, uploaded.ChartName, uploaded.NewFileName)
	file, err := os.Open(join)
	if err != nil {
		logrus.Errorf("openfile dir %s not is dir, %s", join, err.Error())
		panic(err)
	}
	_, err = io.Copy(file, *uploaded.File)
	if err != nil {
		logrus.Errorf("chart dir %s not is dir, %s", join, err.Error())
		panic(err)
	}
}

func checkDirExist(uploaded *controller.FileUploaded) error {
	join := filepath.Join(config.Config.DataDir, uploaded.ChartName)
	fileInfo, err := os.Stat(join)
	//检查文件夹是否存在
	if err != nil && os.IsNotExist(err) {
		err := os.Mkdir(join, os.ModePerm)
		if err != nil {
			logrus.Errorf("create chart dir %s error, %s", join, err.Error())
		}
		err = nil
	}
	if err != nil {
		logrus.Errorf("open chart dir %s error, %s", join, err.Error())
		return err
	}
	//检查是否为文件夹
	if !fileInfo.IsDir() {
		logrus.Errorf("chart dir %s not is dir, %s", join, err.Error())
		return err
	}
	return err
}

/**
处理 chart创建事件
*/
func handlerChartCreated(created *controller.ChartCreated) {
	indexFile.Generated = time.Now()
	//TODO:加上下载前缀
	created.URLs = []string{created.FileName}

	versions, ok := indexFile.Entries[created.Name]
	if !ok {
		versions = make([]*domain.ChartVersion, 1)
	}
	versions = append(versions, created.ChartVersion)
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
	matches, e := filepath.Glob(filepath.Join(config.Config.DataDir, "**/entry.yaml"))
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
		charts, e := loadChartVersionsByFile(value)
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

func loadChartVersionsByFile(filepath string) ([]*domain.ChartVersion, error) {
	versions := make([]*domain.ChartVersion, 0)
	file, e := os.Open(filepath)
	if e != nil {
		return nil, e
	}
	bytes, e := ioutil.ReadAll(file)
	if e != nil {
		return nil, e
	}
	e = yaml.Unmarshal(bytes, &versions)
	if e != nil {
		return nil, e
	}
	return versions, nil
}
