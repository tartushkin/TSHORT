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

	shortURL := fmt.Sprintf("%s%s", h.Short.Address, aliaseURL)

	return ctx.String(http.StatusCreated, shortURL)

}

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
