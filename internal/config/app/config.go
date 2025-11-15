package config

import (
	"flag"
	"os"
)

type Config struct {
	Port    string
	Address string
}

// NewConfig - создание конфигурации приложения
func NewConfig() *Config {
	cfg := Config{}
	flag.StringVar(&cfg.Port, "a", ":8080", "порт сервиса")
	flag.StringVar(&cfg.Address, "b", "http://localhost:8080", "базовый адрес результирующего сокращённого URL")
	flag.Parse()

	runAddr := os.Getenv("SERVER_ADDRESS")
	if runAddr != "" {
		cfg.Port = runAddr
	}
	baseURL := os.Getenv("BASE_URL")
	if runAddr != "" {
		cfg.Address = baseURL
	}
	return &cfg
}
