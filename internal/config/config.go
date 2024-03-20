package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env", env-default:"local", env-required:"true"`
	Storage    `yaml:"storage"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address string        `yaml:"address" env-default:"localhost:8080`
	Timeout time.Duration `yaml:"timeout" env-default:"4s"`
}

type Storage struct {
	User     string `yaml:"user" env-default:"postgres" env-required="true"`
	Password string `yaml:"password" env-required="true"`
	Host     string `yaml:"host" env-default:"localhost" env-required="true"`
	Port     int    `yaml:"port" env-default:"5432"`
	DBName   string `yaml:"dbname" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
