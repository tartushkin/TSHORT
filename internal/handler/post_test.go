package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	config "github.com/tartushkin/TSHORT.git/internal/config/app"
	"github.com/tartushkin/TSHORT.git/internal/service"
)

func TestPostHandler(t *testing.T) {
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
