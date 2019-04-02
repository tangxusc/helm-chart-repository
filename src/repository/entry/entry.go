package entry

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"repository/config"
	"repository/domain"
	"repository/event"
)

func init() {
	event.Subscribe(100, event.Handlers{
		"*domain.ChartCreated": handlerChartCreated,
		"*domain.FileUploaded": handlerFileUploaded,
		"*domain.ChartDeleted": handlerChartDeleted,
	}, "chartEntry")
}

func handlerChartDeleted(event interface{}) {
	deleted := event.(*domain.ChartDeleted)
	path := getEntryFilePath(deleted.ChartName)
	versions := MustLoadChartVersionsByFile(path)
	for key, value := range versions {
		if value.Name == deleted.ChartName && value.Version == deleted.Version {
			versions = append(versions[:key], versions[key+1:]...)
			break
		}
	}
	writeYaml(versions, deleted.ChartName, path)
}

func handlerChartCreated(event interface{}) {
	created := event.(*domain.ChartCreated)

	path := getEntryFilePath(created.Name)
	versions := MustLoadChartVersionsByFile(path)

	versions = append(versions, created.ChartVersion)
	writeYaml(versions, created.Name, path)
}

func writeYaml(versions []*domain.ChartVersion, chartName string, path string) {
	out, err := yaml.Marshal(versions)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ChartName": chartName,
		}).Errorf("marshal entry instance error, %s", err.Error())
		panic(err)
	}
	err = checkDirExist(chartName)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path, out, os.ModePerm)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ChartName": chartName,
		}).Errorf("write entry instance error, %s", err.Error())
		panic(err)
	}
}

func getEntryFilePath(entryName string) string {
	return filepath.Join(config.Config.DataDir, entryName, config.Config.EntryFileName)
}

func handlerFileUploaded(event interface{}) {
	uploaded := event.(*domain.FileUploaded)
	err := checkDirExist(uploaded.ChartName)
	if err != nil {
		panic(err)
	}
	join := filepath.Join(config.Config.DataDir, uploaded.ChartName, uploaded.NewFileName)
	file, err := os.OpenFile(join, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()
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
	defer file.Close()
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

func MustLoadChartVersionsByFile(filepath string) []*domain.ChartVersion {
	versions := make([]*domain.ChartVersion, 0)
	file, e := os.Open(filepath)
	defer func() {
		file.Sync()
		file.Close()
	}()
	if e != nil && os.IsNotExist(e) {
		return versions
	}
	if e != nil {
		logrus.WithFields(logrus.Fields{
			"ChartPath": filepath,
		}).Errorf("open entry instance error, %s", e.Error())
		panic(e)
	}
	bytes, e := ioutil.ReadAll(file)
	if e != nil {
		logrus.WithFields(logrus.Fields{
			"ChartPath": filepath,
		}).Errorf("read entry instance error, %s", e.Error())
		panic(e)
	}
	e = yaml.Unmarshal(bytes, &versions)
	if e != nil {
		logrus.WithFields(logrus.Fields{
			"ChartPath": filepath,
		}).Errorf("unMarshal entry instance error, %s", e.Error())
		panic(e)
	}
	return versions
}

func MustLoadChartVersionByName(chartName string) []*domain.ChartVersion {
	path := getEntryFilePath(chartName)
	return MustLoadChartVersionsByFile(path)
}

func MustLoadChartVersion(chartName, version string) *domain.ChartVersion {
	path := getEntryFilePath(chartName)
	versionsByFile := MustLoadChartVersionsByFile(path)
	for index, value := range versionsByFile {
		if value.Version == version {
			return versionsByFile[index]
		}
	}
	return nil
}
