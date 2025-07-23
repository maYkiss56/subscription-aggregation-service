package config

import (
	"log"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP struct {
		Host         string        `yaml:"host"`
		Port         string        `yaml:"port"`
		ReadTimeout  time.Duration `yaml:"read_timeout"`
		WriteTimeout time.Duration `yaml:"write_timeout"`
		Network      string        `yaml:"network"`
	} `yaml:"http"`

	Postgres struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Database string `yaml:"database"`
		SSLMode  string `yaml:"sslmode"`
		PoolSize int    `yaml:"pool_size"`
	} `yaml:"postgres"`
}

var (
	configPath = "config/config.local.yaml"
	instance   *Config
	once       sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		err := cleanenv.ReadConfig(configPath, instance)
		if err != nil {
			log.Fatalf("read config error: %v", err)
		}
	})

	return instance
}
