package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPostHandler(t *testing.T) {
	//short := &service.Short{
	//	CacheURL: make(map[string]string),
	//}
	//short.PathStorage = configPath
	//handlers := &Handlers{Short: short}

	//file, err := short.NewFile()
	//if err != nil {
	//t.Fatalf("Ошибка при формировании файла: %v", err)
	//}
	//short.File = file

	// Создаём экземпляр Echo
	e := echo.New()

	// 3. Валидный запрос
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://exampl.com"))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Вызываем postHandler
	if assert.NoError(t, TestHandlers.oldPostURLHandler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "text/plain; charset=UTF-8", rec.Header().Get("Content-Type"))

		// Проверяем, что ответ содержит сокращённый URL
		body := rec.Body.String()
		assert.NotEmpty(t, body, "Ожидали не пустой ответ")
	}
}
