package entry

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
)

var entryEventChannel = make(chan interface{}, 1000)
var EntryFileName = "entry.yaml"

func init() {
	event.RegisterChannel(entryEventChannel)
}

func Listen() {
	ok := true
	var evt interface{}
	for {
		select {
		case evt, ok = <-entryEventChannel:
			logrus.WithFields(logrus.Fields{
				"event": evt,
				"ok":    ok,
			}).Debug("entry handler Event")
		}
		if !ok {
			break
		}
		handlerEntryEvent(evt)
	}
}

func handlerEntryEvent(event interface{}) {
	switch event.(type) {
	case *controller.ChartCreated:
		created := event.(*controller.ChartCreated)
		handlerChartCreated(created)
	case *controller.FileUploaded:
		upload := event.(*controller.FileUploaded)
		handlerFileUploaded(upload)
	}
}

func handlerChartCreated(created *controller.ChartCreated) {
	path := getEntryFilePath(created.Name)
	versions, err := LoadChartVersionsByFile(path)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ChartName": created.Name,
		}).Errorf("unMarshal entry instance error, %s", err.Error())
		panic(err)
	}
	versions = append(versions, created.ChartVersion)
	out, err := yaml.Marshal(versions)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ChartName": created.Name,
		}).Errorf("marshal entry instance error, %s", err.Error())
		panic(err)
	}
	err = checkDirExist(created.Name)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path, out, os.ModePerm)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ChartName": created.Name,
		}).Errorf("write entry instance error, %s", err.Error())
		panic(err)
	}
}

func getEntryFilePath(entryName string) string {
	return filepath.Join(config.Config.DataDir, entryName, EntryFileName)
}

func handlerFileUploaded(uploaded *controller.FileUploaded) {
	err := checkDirExist(uploaded.ChartName)
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

func checkDirExist(chartName string) error {
	join := filepath.Join(config.Config.DataDir, chartName)
	fileInfo, err := os.Stat(join)
	//检查文件夹是否存在
	if err != nil && os.IsNotExist(err) {
		direrr := os.Mkdir(join, os.ModePerm)
		if direrr != nil {
			logrus.Errorf("create chart dir %s error, %s", join, direrr.Error())
		}
		fileInfo, err = os.Stat(join)
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

func LoadChartVersionsByFile(filepath string) ([]*domain.ChartVersion, error) {
	versions := make([]*domain.ChartVersion, 0)
	file, e := os.Open(filepath)
	if e != nil && os.IsNotExist(e) {
		return versions, nil
	}
	if e != nil {
		return versions, e
	}
	bytes, e := ioutil.ReadAll(file)
	if e != nil {
		return versions, e
	}
	e = yaml.Unmarshal(bytes, &versions)
	if e != nil {
		return versions, e
	}
	return versions, nil
}
