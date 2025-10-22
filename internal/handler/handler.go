package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (h *Handlers) postHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что это POST-запрос
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Возникал ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	contentType := r.Header.Get("Content-Type")
	if contentType != "text/plain" {
		http.Error(w, "Content-Type must be text/plain", http.StatusBadRequest)
		return
	}

	aliaseURL := h.Short.SetAliaseName(string(body))

	shortURL := fmt.Sprintf("http://localhost:8080/%s", aliaseURL)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte(shortURL))

}
func (h *Handlers) getHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "text/plain" {
		http.Error(w, "Content-Type must be text/plain", http.StatusBadRequest)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
	url, err := h.Short.GetAliasName(parts[1])
	fmt.Println("шо тут у нас", url)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)

}
