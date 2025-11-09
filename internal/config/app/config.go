package config

import (
	"flag"
	"os"
)

type Config struct {
	Port       string
	Address    string
	StorageURL string
}

// NewConfig - создание конфигурации приложения
func NewConfig() *Config {
	cfg := Config{}
	flag.StringVar(&cfg.Port, "a", ":8080", "порт сервиса")
	flag.StringVar(&cfg.Address, "b", "http://localhost:8080", "базовый адрес результирующего сокращённого URL")
	flag.StringVar(&cfg.StorageURL, "c", "./StorageURL.TXT", "путь для файла хранения URL")
	flag.Parse()

	//runAddr := os.Getenv("SERVER_ADDRESS")
	//if runAddr != "" {
	//	cfg.Port = runAddr
	//}
	//baseURL := os.Getenv("BASE_URL")
	//if cfg.Address != "" {
	//	cfg.Address = baseURL
	//}
	//storageURL := os.Getenv("FILE_STORAGE_PATH")
	//if storageURL != "" {
	//	cfg.StorageURL = storageURL
	//}

	if runAddr, exists := os.LookupEnv("SERVER_ADDRESS"); exists && runAddr != "" {
		cfg.Port = runAddr
	}

	if baseURL, exists := os.LookupEnv("BASE_URL"); exists && baseURL != "" {
		cfg.Address = baseURL
	}

	if storageURL, exists := os.LookupEnv("FILE_STORAGE_PATH"); exists && storageURL != "" {
		cfg.StorageURL = storageURL
	}
	return &cfg
}
