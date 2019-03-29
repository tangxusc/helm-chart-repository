package repository

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"repository/config"
	"repository/event"
	"repository/repository/domain"
	"time"
)

func InitIndexFile() {
	_, err := os.Open(getIndexFilePath())
	if err != nil {
		notExist := os.IsNotExist(err)
		if !notExist {
			logrus.Fatalf("open index.yaml error,%s", err.Error())
		}
	}

	indexFile, err := AllIndexFile()
	if err != nil {
		logrus.Fatalf("get index.yaml error,%s", err.Error())
	}
	err = WriteIndexFile(indexFile)
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

func WriteIndexFile(indexFile *domain.IndexFile) error {
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

func AllIndexFile() (*domain.IndexFile, error) {
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
