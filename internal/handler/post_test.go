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

//func TestPostHandler(t *testing.T) {
//	short := &service.Short{
//		CacheURL: make(map[string]string),
//	}
//	handlers := &Handlers{Short: short}
//
//	//  1: Не POST-запрос
//	req := httptest.NewRequest(http.MethodGet, "/", nil)
//	w := httptest.NewRecorder()
//	handlers.postHandler(w, req)
//	if w.Code != http.StatusMethodNotAllowed {
//		t.Errorf("Ожидали %d, получичли %d", http.StatusMethodNotAllowed, w.Code)
//	}
//
//	//  2: Неправильный Content-Type
//	req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://example.com"))
//	req.Header.Set("Content-Type", "application/json")
//	w = httptest.NewRecorder()
//	handlers.postHandler(w, req)
//	if w.Code != http.StatusBadRequest {
//		t.Errorf("Ожидали  %d, получли %d", http.StatusBadRequest, w.Code)
//	}
//
//	// 3: Валидный запрос
//	req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://example.com"))
//	req.Header.Set("Content-Type", "text/plain")
//	w = httptest.NewRecorder()
//	handlers.postHandler(w, req)
//	if w.Code != http.StatusCreated {
//		t.Errorf("Ожидали %d, получли %d", http.StatusCreated, w.Code)
//	}
//	if w.Header().Get("Content-Type") != "text/plain" {
//		t.Errorf("Ожидали Content-Type %s, получили %s", "text/plain", w.Header().Get("Content-Type"))
//	}
//	// Проверяем, что ответ содержит сокращённый URL
//	if len(w.Body.String()) == 0 {
//		t.Error("Ожидали не пустой ответ")
//	}
//}

func TestPostHandler(t *testing.T) {
	// Инициализация
	short := &service.Short{
		CacheURL: make(map[string]string),
	}
	handlers := &Handlers{Short: short}

	// Создаём экземпляр Echo
	e := echo.New()

	// 1. Не POST-запрос
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Вызываем postHandler
	if assert.NoError(t, handlers.postURLHandler(c)) {
		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	}

	// 2. Неправильный Content-Type
	req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://example.com"))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	// Вызываем postHandler
	if assert.NoError(t, handlers.postURLHandler(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}

	// 3. Валидный запрос
	req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://example.com"))
	req.Header.Set("Content-Type", "text/plain")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	// Вызываем postHandler
	if assert.NoError(t, handlers.postURLHandler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "text/plain; charset=UTF-8", rec.Header().Get("Content-Type"))

		// Проверяем, что ответ содержит сокращённый URL
		body := rec.Body.String()
		assert.NotEmpty(t, body, "Ожидали не пустой ответ")
	}
}
