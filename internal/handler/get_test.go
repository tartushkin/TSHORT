package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	config "github.com/tartushkin/TSHORT.git/internal/config/app"
	"github.com/tartushkin/TSHORT.git/internal/service"
)

func TestGetHandler(t *testing.T) {
	// Инициализация

	ctx := context.Background()
	cfg := config.NewConfig()
	lg := logrus.New()
	short, err := service.Create(ctx, lg, cfg)
	if err != nil {
		t.Fatalf("Ошибка инициализации сервиса: %v", err)
	}
	handlers := &Handlers{Short: short}

	// Создаём экземпляр Echo
	e := echo.New()

	// 1. Создаём сокращённый URL через postHandler
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://example.com"))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Вызываем postHandler
	err = handlers.oldPostURLHandler(c)
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
	req = httptest.NewRequest(http.MethodGet, shortURL, nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	parts := strings.Split(shortURL, "/")
	alias := parts[len(parts)-1]
	c.SetParamValues(alias)

	// Вызываем getHandler
	err = handlers.getRedirectHandler(c)
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
