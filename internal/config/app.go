package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/Aderon1333/TPES/pkg/utils/logfacade"
)

type Config struct {
	App     AppConfig     `yaml:"app"`
	Storage StorageConfig `yaml:"storage"`
}

type AppConfig struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type StorageConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Attempts string `json:"attempts"`
}

var once sync.Once

func GetConfig(logger *logfacade.LogFacade) *Config {
	var instance *Config

	once.Do(func() {
		logger.Info("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("../../configs/tpes/config.yaml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Error(err)
		}
	})
	return instance
}
