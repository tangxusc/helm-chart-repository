package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"repository/config"
	"repository/domain"
	"repository/event"
	"repository/repository/index"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	//logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{})

	//加载配置
	config.LoadConfig()
	index.InitIndexFile()
	go event.Listen()

	chart := &domain.ChartVersion{}
	chart.Name = "test"
	chart.Version = "0.1"
	event.Send(&domain.ChartCreated{
		ChartVersion: chart,
		FileName:     "test.tar.gz",
	})
	chart2 := &domain.ChartVersion{}
	chart2.Name = "test"
	chart2.Version = "0.2"
	event.Send(&domain.ChartCreated{
		ChartVersion: chart2,
		FileName:     "test.tar.gz",
	})

	time.Sleep(2 * time.Second)
}

func TestReadTar(t *testing.T) {
	file, e := os.Open("testdata/apache-4.1.0.tgz")
	if e != nil {
		panic(e.Error())
	}
	defer file.Close()
	gzReader, e := gzip.NewReader(file)
	if e != nil {
		panic(e)
	}
	defer gzReader.Close()
	reader := tar.NewReader(gzReader)

	for {
		header, e := reader.Next()
		if e != nil && e == io.EOF {
			break
		}
		if e != nil {
			panic(e)
		}

		fmt.Println("===============文件名称:", header.Name, "===================")
		compile, e := regexp.Compile("([a-zA-Z]*)/?Chart.yaml")
		if e != nil {
			panic(e)
		}
		matchString := compile.MatchString(header.Name)
		fmt.Println(header.Name, matchString)
		allString := compile.FindStringSubmatch(header.Name)
		fmt.Println(allString, "================")
		if matchString {
			bytes, e := ioutil.ReadAll(reader)
			if e != nil {
				panic(e)
			}
			fmt.Println("==============", header.Name, ",内容===========")
			fmt.Println(string(bytes))
		}

	}

}
