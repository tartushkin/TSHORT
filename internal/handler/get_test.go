package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/tartushkin/TSHORT.git/internal/service"
)

func TestGetHandler(t *testing.T) {
	// Инициализация
	short := &service.Short{
		CacheURL: make(map[string]string),
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
	if assert.NoError(t, handlers.oldPostURLHandler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
	shortURL := rec.Body.String()

	// 3. Создаём маршрут для GET-запроса
	req = httptest.NewRequest(http.MethodGet, "/:"+shortURL, nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(shortURL)

	// Вызываем getHandler
	if assert.NoError(t, handlers.getRedirectHandler(c)) {
		assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
		assert.Equal(t, "https://example.com", rec.Header().Get("Location"))
	}
}
