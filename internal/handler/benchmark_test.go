package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func BenchmarkGetUserURLsHandler(b *testing.B) {
	// Создаём экземпляр Echo
	h := testCreate()
	// Создаём экземпляр Echo
	e := echo.New()

	// 3. Валидный запрос
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://exampl.com"))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()
	req.AddCookie(&http.Cookie{
		Name:     "user_id",
		Value:    "a7fa515e.fb64b57c209853ea0fb0fb0f2adb85aa6f1ef726d18c38c44d530036b5e7282f",
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Используйте true, если используете HTTPS
	})
	c := e.NewContext(req, rec)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		c = e.NewContext(req, rec)

		_ = h.oldPostURLHandler(c)
	}
}
