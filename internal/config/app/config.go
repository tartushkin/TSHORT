package config

import (
	"flag"
	"os"
)

type Config struct {
	Port       string
	Address    string
	StorageURL string
	DBHost     string
	PassConnDB string
	DBName     string
	DBUser     string
	DBPort     string
	DNS        string
}

// NewConfig - создание конфигурации приложения
func NewConfig() *Config {
	cfg := Config{}
	flag.StringVar(&cfg.Port, "a", ":8080", "порт сервиса")
	flag.StringVar(&cfg.Address, "b", "http://localhost:8080", "базовый адрес результирующего сокращённого URL")
	flag.StringVar(&cfg.StorageURL, "c", "./StorageURL.TXT", "путь для файла хранения URL")
	flag.StringVar(&cfg.DNS, "g", "host=localhost port=5432 user=postgres password=12345678 dbname=postgres sslmode=disable", "cтрока с адресом подключения к БД")
	flag.StringVar(&cfg.DBHost, "d", "localhost", "хост до БД")
	flag.StringVar(&cfg.DBPort, "t", "5432", "порт до БД")
	flag.StringVar(&cfg.PassConnDB, "e", "12345678", "пароль подлючения к БД")
	flag.StringVar(&cfg.DBName, "o", "postgres", "наименование БД")
	flag.StringVar(&cfg.DBUser, "i", "postgres", "пользователь БД")

	flag.Parse()

	if runAddr, exists := os.LookupEnv("SERVER_ADDRESS"); exists && runAddr != "" {
		cfg.Port = runAddr
	}

	if baseURL, exists := os.LookupEnv("BASE_URL"); exists && baseURL != "" {
		cfg.Address = baseURL
	}

	if storageURL, exists := os.LookupEnv("FILE_STORAGE_PATH"); exists && storageURL != "" {
		cfg.StorageURL = storageURL
	}
	if db, exists := os.LookupEnv("DATABASE_DSN"); exists && db != "" {
		cfg.DNS = db
	}
	if host, exists := os.LookupEnv("DB_HOST"); exists && host != "" {
		cfg.DBHost = host
	}
	if dbName, exists := os.LookupEnv("DB_NAME"); exists && dbName != "" {
		cfg.DBName = dbName
	}
	if pass, exists := os.LookupEnv("DB_PASS"); exists && pass != "" {
		cfg.PassConnDB = pass
	}
	if user, exists := os.LookupEnv("DB_USER"); exists && user != "" {
		cfg.DBUser = user
	}
	if port, exists := os.LookupEnv("DB_PORT"); exists && port != "" {
		cfg.DBHost = port
	}

	return &cfg
}
