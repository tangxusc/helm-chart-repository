package config

import (
	"context"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"github.com/heetch/confita/backend/flags"
	log "github.com/sirupsen/logrus"
)

var loader = confita.NewLoader(flags.NewBackend(),
	env.NewBackend(),
	file.NewBackend("config.json"),
	file.NewBackend("config.yaml"),
)

var Config = Configuration{
	ServerPort:    "8080",
	DataDir:       "data/",
	EntryFileName: "entry.yaml",
	Domain:        "http://localhost:8080",
}

type Configuration struct {
	ServerPort    string `config:"serverPort"`
	DataDir       string `config:"dataDir"`
	EntryFileName string `config:"entryFileName"`
	Domain        string `config:"domain"`
}

func LoadConfig() {
	err := loader.Load(context.Background(), &Config)
	if err != nil {
		log.Warnf("加载配置,%s", err.Error())
	}
}
