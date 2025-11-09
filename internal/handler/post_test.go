package handler

import (
	"bytes"
	"flag"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/tartushkin/TSHORT.git/internal/service"
)

func TestPostHandler(t *testing.T) {
	flag.StringVar(&configPath, "c", "./StorageURL.TXT", "путь для файла хранения URL")
	short := &service.Short{
		CacheURL: make(map[string]string),
	}
	short.PathStorage = configPath
	handlers := &Handlers{Short: short}
	file, err := short.NewFile()
	if err != nil {
		t.Fatalf("Ошибка при формировании файла: %v", err)
	}
	short.File = file

	// Создаём экземпляр Echo
	e := echo.New()

	// 3. Валидный запрос
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://example.com"))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Вызываем postHandler
	if assert.NoError(t, handlers.oldPostURLHandler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "text/plain; charset=UTF-8", rec.Header().Get("Content-Type"))

		// Проверяем, что ответ содержит сокращённый URL
		body := rec.Body.String()
		assert.NotEmpty(t, body, "Ожидали не пустой ответ")
	}
}
