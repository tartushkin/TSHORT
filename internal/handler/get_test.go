package handler

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	cfg "github.com/tartushkin/TSHORT.git/internal/config/app"
	sr "github.com/tartushkin/TSHORT.git/internal/service"
)

// var configPath, DNS string
var TestHandlers *Handlers

func TestMain(m *testing.M) {
	lg := logrus.New()
	ctx := context.Background()
	cfg := cfg.NewConfig()

	sh, err := sr.Create(ctx, lg, cfg)
	if err != nil {
		log.Fatalf("Ошибка создания сервиса: %v", err)
	}
	defer sh.Close()
	TestHandlers = &Handlers{Short: sh}
	//go TestHandlers.StartHTTP(ctx, cfg.Port)

	// Запуск тестов
	code := m.Run()

	os.Exit(code)
}

func TestGetHandler(t *testing.T) {
	// Инициализация
	//flag.StringVar(&port, "a", ":8080", "порт сервиса")
	//flag.StringVar(&address, "b", "http://localhost:8080", "базовый адрес результирующего сокращённого URL")
	//flag.StringVar(&configPath, "c", "./StorageURL.TXT", "путь для файла хранения URL")
	//flag.StringVar(&DNS, "d", "host=localhost port=5432 user=postgres password=12345678 dbname=myDB sslmode=disable", "cтрока с адресом подключения к БД")

	//short := &service.Short{
	//	CacheURL: make(map[string]string),
	//	Logger:   logrus.New(),
	//}
	//short.PathStorage = configPath
	//short.DNS = DNS
	//short.Ctx
	//lg := logrus.New()
	//ctx := context.Background()
	//cfg := cfg.NewConfig() // инициализация конфига
	//conn, _ := db.NewConnection(DNS)
	//sh, err := sr.Create(ctx, lg, cfg) // инициализация сервиса
	//if err != nil {
	//	panic(err)
	//}
	//short.Address = address
	//short.HTTPPort = port
	//Testhandlers := &Handlers{Short: sh}
	file, err := TestHandlers.Short.NewFile()
	if err != nil {
		t.Fatalf("Ошибка при формировании файла: %v", err)
	}
	TestHandlers.Short.File = file
	// Создаём экземпляр Echo
	e := echo.New()

	// 1. Создаём сокращённый URL через postHandler
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://example.com"))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Вызываем postHandler
	err = TestHandlers.oldPostURLHandler(c)
	if err != nil {
		t.Fatalf("Ошибка в postHandler: %v", err)
	}
	if rec.Code != http.StatusCreated {
		t.Fatalf("Ожидался статус %d, получен %d", http.StatusCreated, rec.Code)
	}

	shortURL := strings.TrimSpace(rec.Body.String())
	if shortURL == "" {
		t.Fatal("Сокращённый URL не должен быть пустым")
	}

	// 2. Создаём маршрут для GET-запроса
	req = httptest.NewRequest(http.MethodGet, "/:shortURL", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	parts := strings.Split(shortURL, "/")
	alias := parts[len(parts)-1]
	c.SetParamValues(alias)

	// Вызываем getHandler
	err = TestHandlers.getRedirectHandler(c)
	if err != nil {
		t.Fatalf("Ошибка в getHandler: %v", err)
	}
	if rec.Code != http.StatusTemporaryRedirect {
		t.Fatalf("Ожидался статус %d, получен %d", http.StatusTemporaryRedirect, rec.Code)
	}
	if rec.Header().Get("Location") != "https://example.com" {
		t.Fatalf("Ожидался Location %s, получен %s", "https://example.com", rec.Header().Get("Location"))
	}

}
