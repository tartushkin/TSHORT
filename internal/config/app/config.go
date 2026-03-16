// Package config отвечает за загрузку и инициализацию конфигурации приложения.
//
// Конфигурация может быть задана:
//   - через флаги командной строки
//   - через переменные окружения (имеют приоритет)
//
// Пример использования:
//
//	cfg := config.NewConfig()
//	fmt.Println(cfg.Port)
package config

import (
	"flag"
	"os"
)

const (
	trueV  = "true"
	falseV = "false"
)

// Config хранит параметры конфигурации приложения.
type Config struct {
	// Port — порт, на котором запускается HTTP-сервер.
	// Флаг: -a, переменная: SERVER_ADDRESS
	Port string
	// Address — базовый URL для генерации сокращённых ссылок.
	// Флаг: -b, переменная: BASE_URL
	Address string
	// FileStoragePath — путь к файлу для хранения URL (если используется файловое хранилище).
	// Флаг: -c, переменная: FILE_STORAGE_PATH

	FileStoragePath string

	// DNS — строка подключения к PostgreSQL.
	// Флаг: -d, переменная: DATABASE_DSN
	DNS string

	// LocalAuditPath — путь к локальному файлу аудита.
	// Флаг: -audit-file, переменная: AUDIT_FILE
	LocalAuditPath string
	// AuditPath — URL внешнего сервиса аудита.
	// Флаг: -audit-url, переменная: AUDIT_URL
	AuditPath string
	// ParamDelete — интервал (в секундах) проверки помеченных на удаление URL.
	// Флаг: -t
	ParamDelete int
	// SecretKey — ключ для подписи сессий и cookies.
	// Флаг: -k
	SecretKey string
	// RunProfile — включает профилирование CPU и памяти.
	// Флаг: -p (если указан — true)
	RunProfile bool

	runProfile bool

	TLSconn  bool
	СertFile string
	KeyFile  string
}

// NewConfig - создание конфигурации приложения.
//
// Значения устанавливаются в порядке приоритета:
// 1. Переменные окружения (если заданы)
// 2. Флаги командной строки
// 3. Значения по умолчанию
//
// Возвращает указатель на инициализированный *Config.
func NewConfig() *Config {
	cfg := Config{}
	// обязательные
	flag.StringVar(&cfg.Port, "a", ":8080", "порт сервиса")
	flag.StringVar(&cfg.Address, "b", "http://localhost:8080", "базовый адрес результирующего сокращённого URL")
	flag.StringVar(&cfg.FileStoragePath, "c", "./StorageURL.TXT", "путь для файла хранения URL")
	flag.StringVar(&cfg.DNS, "d", "postgres://postgres:12345678@localhost:5432/myDB?sslmode=disable", "cтрока с адресом подключения к БД")
	flag.BoolVar(&cfg.TLSconn, "s", false, " возможность включения HTTPS в веб-сервере")

	// аудит
	flag.StringVar(&cfg.LocalAuditPath, "audit-file", "./Audit.TXT", "cтрока с адресом подключения к локальному аудит файлу")
	flag.StringVar(&cfg.AuditPath, "audit-url", "", "cтрока с адресом подключения к внешнему аудит")

	// индивидуальные
	flag.IntVar(&cfg.ParamDelete, "t", 20, "частота запуска очистки от помеченных на удаление URL")
	flag.StringVar(&cfg.SecretKey, "k", "tort-secret-key", "ключ")
	flag.BoolVar(&cfg.runProfile, "p", false, "флаг необходимоти профилирования сервиса")

	flag.Parse()
	// переменные окружения (имеют приоритет)
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
	if sec, exists := os.LookupEnv("ENABLE_HTTPS"); exists && sec != "" {
		if sec == trueV {
			cfg.TLSconn = true
		} else {
			cfg.TLSconn = false
		}
	}

	// Переопределяем пути к сертификату и ключу из переменных окружения, если они заданы
	if envCertFile, exists := os.LookupEnv("CERT_FILE"); exists {
		cfg.СertFile = envCertFile
	}
	if envKeyFile, exists := os.LookupEnv("KEY_FILE"); exists {
		cfg.KeyFile = envKeyFile
	}

	//audit
	if local, exists := os.LookupEnv("AUDIT_FILE"); exists && local != "" {
		cfg.LocalAuditPath = local
	}
	if auditURL, exists := os.LookupEnv("AUDIT_URL"); exists && auditURL != "" {
		cfg.AuditPath = auditURL
	}
	return &cfg
}
