package config

import "flag"

type Configure struct {
	Port    string
	Address string
}

// NewConfig - создание конфигурации приложения
func NewConfig() *Configure {
	cfg := Configure{}
	flag.StringVar(&cfg.Port, "a", ":8080", "порт сервиса")
	flag.StringVar(&cfg.Address, "b", "http://127.0.0.1", "базовый адрес результирующего сокращённого URL")
	flag.Parse()
	return &cfg
}
