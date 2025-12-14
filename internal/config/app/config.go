package config

import (
	"flag"
	"os"
)

type Config struct {
	Port            string
	Address         string
	FileStoragePath string
	DNS             string
	ParamDelete     int
}

// NewConfig - создание конфигурации приложения
func NewConfig() *Config {
	cfg := Config{}
	flag.StringVar(&cfg.Port, "a", ":8080", "порт сервиса")
	flag.StringVar(&cfg.Address, "b", "http://localhost:8080", "базовый адрес результирующего сокращённого URL")
	flag.StringVar(&cfg.FileStoragePath, "c", "./StorageURL.TXT", "путь для файла хранения URL")
	flag.StringVar(&cfg.DNS, "d", "postgres://postgres:12345678@localhost:5432/myDB?sslmode=disable", "cтрока с адресом подключения к БД")
	flag.IntVar(&cfg.ParamDelete, "t", 20, "частота запуска очистки от помеченных на удаление URL")

	flag.Parse()

	if runAddr, exists := os.LookupEnv("SERVER_ADDRESS"); exists && runAddr != "" {
		cfg.Port = runAddr
	}

	if baseURL, exists := os.LookupEnv("BASE_URL"); exists && baseURL != "" {
		cfg.Address = baseURL
	}

	if fileStoragePath, exists := os.LookupEnv("FILE_STORAGE_PATH"); exists && fileStoragePath != "" {
		cfg.FileStoragePath = fileStoragePath
	}
	if db, exists := os.LookupEnv("DATABASE_DSN"); exists && db != "" {
		cfg.DNS = db
	}

	return &cfg
}
