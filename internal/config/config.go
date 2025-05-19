package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	HTTPServer      `yaml:"server"`
	Datasource      `yaml:"datasource"`
	WeatherProvider `yaml:"weather-provider"`
}

type HTTPServer struct {
	Host string `yaml:"host"`
	Port string `yaml:"port" env-default:"8080"`
}

type Datasource struct {
	Url string `yaml:"url" env:"DATABASE_URL"`
}

type WeatherProvider struct {
	Url string `yaml:"url"`
	Key string `yaml:"key"`
}

func ReadConfig(configPath string) *Config {
	if configPath == "" {
		log.Fatal("configPath is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
