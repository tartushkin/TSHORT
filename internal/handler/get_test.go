package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/tartushkin/TSHORT.git/internal/service"
)

func TestGetHandler(t *testing.T) {
	short := &service.Short{
		CacheURL: make(map[string]string),
	}
	handlers := &Handlers{Short: short}
	// даем новый url
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://example.com"))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()
	handlers.postHandler(w, req)
	URL := w.Body.String()

	//1. Не Get - запрос
	req = httptest.NewRequest(http.MethodPost, "/", nil)
	w = httptest.NewRecorder()
	handlers.getHandler(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Ожидали %d, получичли %d", http.StatusMethodNotAllowed, w.Code)
	}
	// Создаём маршрутизатор
	r := mux.NewRouter()
	r.HandleFunc("/{id}", handlers.getHandler).Methods("GET")
	// 3:Валидный запрос  
	parts := strings.Split(URL, "/")
	URL = parts[3]
	req = httptest.NewRequest(http.MethodGet, "/"+URL, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// Проверяем, что ответ содержит сокращённый URL
	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Ожидали статус %d, получли %d", http.StatusTemporaryRedirect, w.Code)
	}
	if w.Header().Get("Location") != "https://example.com" {
		t.Errorf("Ожидали %s, получили %s", "https://example.com", w.Header().Get("Location"))
	}
}
