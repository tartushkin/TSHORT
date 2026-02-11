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

	LocalAuditPath string
	AuditPath      string

	ParamDelete int
	SecretKey   string

	RunProfile bool
	runProfile string
}

// NewConfig - создание конфигурации приложения
func NewConfig() *Config {
	cfg := Config{}
	// обязательные
	flag.StringVar(&cfg.Port, "a", ":8080", "порт сервиса")
	flag.StringVar(&cfg.Address, "b", "http://localhost:8080", "базовый адрес результирующего сокращённого URL")
	flag.StringVar(&cfg.FileStoragePath, "c", "./StorageURL.TXT", "путь для файла хранения URL")
	flag.StringVar(&cfg.DNS, "d", "postgres://postgres:12345678@localhost:5432/myDB?sslmode=disable", "cтрока с адресом подключения к БД")

	flag.StringVar(&cfg.LocalAuditPath, "audit-file", "./Audit.TXT", "cтрока с адресом подключения к локальному аудит файлу")
	flag.StringVar(&cfg.AuditPath, "audit-url", "", "cтрока с адресом подключения к внешнему аудит")

	//индивидуальные
	flag.IntVar(&cfg.ParamDelete, "t", 20, "частота запуска очистки от помеченных на удаление URL")
	flag.StringVar(&cfg.SecretKey, "k", "tort-secret-key", "ключ")
	flag.StringVar(&cfg.runProfile, "p", "", "флаг необходимоти профилирования сервиса")

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

	//audit
	if local, exists := os.LookupEnv("AUDIT_FILE"); exists && local != "" {
		cfg.LocalAuditPath = local
	}
	if auditURL, exists := os.LookupEnv("AUDIT_URL"); exists && auditURL != "" {
		cfg.AuditPath = auditURL
	}
	//profile

	if cfg.runProfile == "" {
		cfg.RunProfile = false
	} else {
		cfg.RunProfile = true
	}

	return &cfg
}
