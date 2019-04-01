package index

import (
	"fmt"
	"repository/config"
	"repository/repository/entry"
	"testing"
)

func TestAll(t *testing.T) {
	config.Config.DataDir = "../testdata/data/"
	indexFile, e := LoadIndexFile()
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
	versions, e := entry.LoadChartVersionsByFile("../testdata/data/apache/entry.yaml")
	if e != nil {
		panic(e)
	}
	fmt.Println(versions[0])
}
