package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) postURLHandler(ctx echo.Context) error {
	if ctx.Request().Method != http.MethodPost {
		return ctx.String(http.StatusMethodNotAllowed, "Не соответствует метод запроса")
	}
	// Получаем значение заголовка Content-Type
	contentType := ctx.Request().Header.Get("Content-Type")

	// Проверяем, что Content-Type равен "text/plain"
	if contentType != "text/plain" {
		return ctx.String(http.StatusBadRequest, "Content-Type не соответсвует ожидаемому: text/plain")
	}

	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Возникал ошибка при чтении тела запроса: "+err.Error())
	}

	aliaseURL := h.Short.SetAliaseName(string(body))

	shortURL := fmt.Sprintf("http://localhost:8080/%s", aliaseURL)
	return ctx.String(http.StatusCreated, shortURL)

}

//	func (h *Handlers) postHandler(w http.ResponseWriter, r *http.Request) {
//		// Проверяем, что это POST-запрос
//		if r.Method != http.MethodPost {
//			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//			return
//		}
//
//		contentType := r.Header.Get("Content-Type")
//		if contentType != "text/plain" {
//			http.Error(w, "Content-Type must be text/plain", http.StatusBadRequest)
//			return
//		}
//		body, err := io.ReadAll(r.Body)
//		if err != nil {
//			http.Error(w, "Возникал ошибка при чтении тела запроса", http.StatusBadRequest)
//			return
//		}
//		defer r.Body.Close()
//
//		aliaseURL := h.Short.SetAliaseName(string(body))
//
//		shortURL := fmt.Sprintf("http://localhost:8080/%s", aliaseURL)
//		w.Header().Set("Content-Type", "text/plain")
//		w.WriteHeader(http.StatusCreated)
//
//		w.Write([]byte(shortURL))
//
// }

func (h *Handlers) getRedirectHandler(ctx echo.Context) error {
	if ctx.Request().Method != http.MethodGet {
		return ctx.String(http.StatusMethodNotAllowed, "Не соответствует метод запроса")
	}
	// Получаем URL из параметров запроса
	alias := ctx.Param("id")
	if alias == "" {
		return ctx.String(http.StatusBadRequest, "Требуется алиас")
	}
	// Извлекаем алиас
	originalURL, err := h.Short.GetAliasName(alias)
	if err != nil {
		return ctx.String(http.StatusNotFound, "URL не найден")
	}

	// Возвращаем редирект
	return ctx.Redirect(http.StatusTemporaryRedirect, originalURL)
}

//func (h *Handlers) getHandler(w http.ResponseWriter, r *http.Request) {
//
//	if r.Method != http.MethodGet {
//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//		return
//	}
//	parts := strings.Split(r.URL.Path, "/")
//	if len(parts) < 2 {
//		http.Error(w, "Not Found", http.StatusNotFound)
//	}
//	url, err := h.Short.GetAliasName(parts[1])
//
//	if err != nil {
//		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
//	}
//	w.Header().Set("Location", url)
//	w.WriteHeader(http.StatusTemporaryRedirect)
//
//}
