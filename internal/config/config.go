package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/viper"
)

type Config struct {
	AppCfg      AppConfig      `yaml:"app"`
	JWTCfg      JWT            `yaml:"jwt"`
	GrpcCfg     Grpc           `yaml:"grpc"`
	PostgresCfg PostgresConfig `yaml:"postgres"`
	MongoCfg    MongoConfig    `yaml:"mongo"`
}

type AppConfig struct {
	IP   string `yaml:"ip" env:"APP_IP"`
	Port string `yaml:"port" env:"APP_PORT"`
	DB   string `yaml:"db" env:"APP_DB"`
	Log  string `yaml:"log" env:"LOG_PATH"`
}

type JWT struct {
	Secret        string `yaml:"secret" env:"JWT_SECRET"`
	AccessCookie  string `yaml:"access_cookie" env:"ACCESS_COOKIE"`
	RefreshCookie string `yaml:"refresh_cookie" env:"REFRESH_COOKIE"`
}

type Grpc struct {
	Port    string `yaml:"port"`
	Timeout string `yaml:"timeout"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" env:"PG_HOST"`
	Port     string `yaml:"port" env:"PG_PORT"`
	Database string `yaml:"database" env:"PG_DATABASE"`
	Username string `yaml:"username" env:"PG_USERNAME"`
	Password string `yaml:"password" env:"PG_PASSWORD"`
}

type MongoConfig struct {
	Url string `yaml:"url" env:"MONGO_URL"`
}

func GetConfig() *Config {
	var instance Config

	// Проверяем, что передан путь к файлу конфигурации
	if len(os.Args) < 2 {
		// Если нет пути, пробуем .env
		viper.SetConfigType("env")
		viper.AutomaticEnv()
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Error reading config file: %w", err))
		}
	} else {
		// Получаем путь к файлу конфигурации
		configPath := os.Args[1]

		err := cleanenv.ReadConfig(configPath, &instance)
		if err != nil {
			panic(fmt.Errorf("Error reading config file: %w", err))
		}
	}

	return &instance
}
