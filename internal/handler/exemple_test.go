//go:build examples
// +build examples

package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tartushkin/TSHORT.git/internal/handler"
	"github.com/tartushkin/TSHORT.git/internal/model"
	"github.com/tartushkin/TSHORT.git/internal/repository"
	"github.com/tartushkin/TSHORT.git/internal/service"
)

// ExampleHandlers_postURLHandler показывает, как сократить один URL.
//
// Запрос:
//
//	POST / HTTP/1.1
//	Content-Type: application/json
//
//	{"url":"https://example.com"}
//
// Ответ:
//
//	HTTP/1.1 201 Created
//	Content-Type: application/json
//
//	{"result":"http://localhost:8080/abc123"}
func ExampleHandlers_postURLHandler() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"url":"https://example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Мок сервиса
	svc := &service.Short{
		Repo:     &repository.Repo{}, // можно заменить на мок
		Address:  "http://localhost:8080",
		Logger:   logrus.New(),
		CacheURL: make(map[string]*model.AliasFullCore),
	}
	h := handler.NewHandlers(svc)

	c := e.NewContext(req, rec)
	h.PostURLHandler(c) // убедись, что метод экспортируем (с большой буквы)

	fmt.Println(rec.Code)
	fmt.Println(rec.Header().Get("Content-Type"))
	// Output:
	// 201
	// application/json
}

// ExampleHandlers_getRedirectHandler показывает редирект по короткому ключу.
//
// Запрос:
//
//	GET /abc123 HTTP/1.1
//
// Ответ:
//
//	HTTP/1.1 307 Temporary Redirect
//	Location: https://example.com
func ExampleHandlers_getRedirectHandler() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	rec := httptest.NewRecorder()

	svc := &service.Short{
		Repo:    &repository.Repo{},
		Address: "http://localhost:8080",
		Logger:  logrus.New(),
		CacheURL: map[string]*model.AliasFullCore{
			"abc123": {
				Alias:       "abc123",
				OriginalURL: "https://example.com",
			},
		},
	}
	h := handler.NewHandlers(svc)

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc123")
	h.GetRedirectHandler(c)

	fmt.Println(rec.Code)
	fmt.Println(rec.Header().Get("Location"))
	// Output:
	// 307
	// https://example.com
}

// ExampleHandlers_batchHandler показывает массовое сокращение URL.
//
// Запрос:
//
//	POST /api/shorten/batch HTTP/1.1
//	Content-Type: application/json
//
//	[{"correlation_id":"1","original_url":"https://yandex.ru"},{"correlation_id":"2","original_url":"https://google.com"}]
//
// Ответ:
//
//	HTTP/1.1 201 Created
//	Content-Type: application/json
//
//	[{"correlation_id":"1","short_url":"http://localhost:8080/xyz"},{"correlation_id":"2","short_url":"http://localhost:8080/uvw"}]
func ExampleHandlers_batchHandler() {
	e := echo.New()
	body := []byte(`[{"correlation_id":"1","original_url":"https://yandex.ru"},{"correlation_id":"2","original_url":"https://google.com"}]`)
	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	svc := &service.Short{
		Repo:     &repository.Repo{},
		Address:  "http://localhost:8080",
		Logger:   logrus.New(),
		CacheURL: make(map[string]*model.AliasFullCore),
	}
	h := handler.NewHandlers(svc)

	c := e.NewContext(req, rec)
	h.BatchHandler(c)

	fmt.Println(rec.Code)
	fmt.Println(rec.Header().Get("Content-Type"))
	// Output:
	// 201
	// application/json
}

// ExampleHandlers_getMyShortURL показывает получение списка URL пользователя.
//
// Запрос:
//
//	GET /api/user/urls HTTP/1.1
//	Cookie: user_id=abc123
//
// Ответ:
//
//	HTTP/1.1 200 OK
//	Content-Type: application/json
//
//	[{"short_url":"http://...","original_url":"https://..."},...]
func ExampleHandlers_getMyShortURL() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	req.AddCookie(&http.Cookie{Name: "user_id", Value: "testuser"})
	rec := httptest.NewRecorder()

	svc := &service.Short{
		Repo:   &repository.Repo{},
		Logger: logrus.New(),
	}
	// Мок ответа
	svc.GetUserURL = func(userID string) ([]model.UserURLResponse, error) {
		return []model.UserURLResponse{
			{ShortURL: "http://localhost:8080/abc123", OriginalURL: "https://example.com"},
		}, nil
	}

	h := handler.NewHandlers(svc)
	c := e.NewContext(req, rec)
	h.GetMyShortURL(c)

	fmt.Println(rec.Code)
	fmt.Println(len(rec.Body.String()) > 0)
	// Output:
	// 200
	// true
}

// ExampleHandlers_deleteURL показывает пометку URL на удаление.
//
// Запрос:
//
//	DELETE /api/user/urls HTTP/1.1
//	Content-Type: application/json
//	Cookie: user_id=testuser
//
//	["abc123","def456"]
//
// Ответ:
//
//	HTTP/1.1 202 Accepted
func ExampleHandlers_deleteURL() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(`["abc123","def456"]`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "user_id", Value: "testuser"})
	rec := httptest.NewRecorder()

	svc := &service.Short{
		Repo:   &repository.Repo{},
		Logger: logrus.New(),
	}
	// Мок
	svc.DeleteUserURL = func(userID string, aliases []string) error {
		return nil
	}

	h := handler.NewHandlers(svc)
	c := e.NewContext(req, rec)
	h.DeleteURL(c)

	fmt.Println(rec.Code)
	// Output:
	// 202
}
