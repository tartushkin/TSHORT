package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/tartushkin/TSHORT.git/internal/config/db"
	"github.com/tartushkin/TSHORT.git/internal/repository"
)

func BenchmarkGetUserURLsHandler(b *testing.B) {
	// Создаём экземпляр Echo
	h := testCreate()
	err := mockConn(h)
	if err != nil {
		h.Short.Logger.Error("db: ошибка - ", err.Error())
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	req.AddCookie(&http.Cookie{
		Name:     "user_id",
		Value:    "a7fa515e.fb64b57c209853ea0fb0fb0f2adb85aa6f1ef726d18c38c44d530036b5e7282f",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		_ = h.getMyShortURL(c)
	}
}
func mockConn(h *Handlers) error {
	if DNS != "" {
		conn, err := db.NewConnection(DNS)
		if err != nil {
			h.Short.Logger.Info("db: не удалось подключилиться к DB, используем другое хранилище")
			return err
		} else {
			h.Short.DNS = DNS
			h.Short.Logger.Info("db: успешно подключились к DB")
			h.Short.Conn = conn
			h.Short.Repo = repository.NewRepository(conn)
		}

	}
	return nil
}
