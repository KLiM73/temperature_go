package main

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	Telegram struct {
		Token    string `envconfig:"TG_TOKEN" required:"true"`
		Username string `envconfig:"TG_USERNAME" required:"true"`
	}
	Backend struct {
		Url         string        `envconfig:"BACKEND_URL" required:"true"`
		SendTimeout time.Duration `envconfig:"BACKEND_SEND_TIMEOUT" required:"true"`
	}
	DbFile string `envconfig:"DB_FILE" required:"true"`
}

func InitConfig() Config {
	var config Config
	envconfig.Process("", &config)

	return config
}
