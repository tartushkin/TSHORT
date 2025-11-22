package handler

import (
	"bytes"
	"context"
	"flag"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tartushkin/TSHORT.git/internal/config/db"
	"github.com/tartushkin/TSHORT.git/internal/repository"
	"github.com/tartushkin/TSHORT.git/internal/service"
)

// var configPath, DNS string
// var TestHandlers *Handlers
// var configPath string
var DNS string

//	func TestMain(m *testing.M) {
//		lg := logrus.New()
//		ctx := context.Background()
//		cfg := cfg.NewConfig()
//
//		sh, err := sr.Create(ctx, lg, cfg)
//		if err != nil {
//			log.Fatalf("Ошибка создания сервиса: %v", err)
//		}
//		defer sh.Close()
//		TestHandlers = &Handlers{Short: sh}
//
//		// Запуск тестов
//		code := m.Run()
//
//		os.Exit(code)
//	}
func testCreate() *Handlers {
	if DNS == "" {
		flag.StringVar(&DNS, "d", "host=localhost port=5432 user=postgres password=12345678 dbname=myDB sslmode=disable", "cтрока с адресом подключения к БД")
	}
	short := &service.Short{
		CacheURL: make(map[string]string),
		Logger:   logrus.New(),
	}
	//short.PathStorage = configPath
	short.DNS = DNS
	short.Ctx = context.Background()
	conn, err := db.NewConnection(short.DNS)
	if err != nil {
		panic(err)
	}
	short.Logger.Info("db: успешно подключились к DB")

	short.Repo = repository.NewRepository(conn)

	TestHandlers := &Handlers{Short: short}
	//file, err := TestHandlers.Short.NewFile()
	//if err != nil {
	//	short.Logger.Fatalf("Ошибка при формировании файла: %v", err)
	//}
	//TestHandlers.Short.File = file
	return TestHandlers
}
func TestGetHandler(t *testing.T) {
	// Инициализация
	//flag.StringVar(&port, "a", ":8080", "порт сервиса")
	//flag.StringVar(&address, "b", "http://localhost:8080", "базовый адрес результирующего сокращённого URL")
	//flag.StringVar(&configPath, "c", "./StorageURL.TXT", "путь для файла хранения URL")
	//flag.StringVar(&DNS, "d", "host=localhost port=5432 user=postgres password=12345678 dbname=myDB sslmode=disable", "cтрока с адресом подключения к БД")
	//
	//short := &service.Short{
	//	CacheURL: make(map[string]string),
	//	Logger:   logrus.New(),
	//}
	////short.PathStorage = configPath
	//short.DNS = DNS
	//short.Ctx = context.Background()
	//
	//conn, err := db.NewConnection(short.DNS)
	//if err != nil {
	//	panic(err)
	//}
	//short.Logger.Info("db: успешно подключились к DB")
	//
	//short.Repo = repository.NewRepository(conn)
	//
	//TestHandlers := &Handlers{Short: short}
	//file, err := TestHandlers.Short.NewFile()
	//if err != nil {
	//	t.Fatalf("Ошибка при формировании файла: %v", err)
	//}
	//TestHandlers.Short.File = file
	// Создаём экземпляр Echo
	h := testCreate()
	e := echo.New()

	// 1. Создаём сокращённый URL через postHandler
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://example.com"))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Вызываем postHandler
	err := h.oldPostURLHandler(c)
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
	err = h.getRedirectHandler(c)
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
