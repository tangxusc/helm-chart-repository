package repository

import (
	"fmt"
	"repository/config"
	"testing"
)

func TestAll(t *testing.T) {
	config.Config.DataDir = "../testdata/data/"
	indexFile, e := AllIndexFile()
	if e != nil {
		panic(e)
	}
	fmt.Println(indexFile)
}

func TestInitIndexFile(t *testing.T) {
	config.Config.DataDir = "../testdata/data/"
	InitIndexFile()
}

func TestLoadChartVersionsByFile(t *testing.T) {
	versions, e := loadChartVersionsByFile("../testdata/data/apache/entry.yaml")
	if e != nil {
		panic(e)
	}
	fmt.Println(versions[0])
}
