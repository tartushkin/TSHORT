package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tartushkin/TSHORT.git/internal/model"
)

func (h *Handlers) oldPostURLHandler(ctx echo.Context) error {
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

	aliasURL, err := h.Short.SetAliasName(string(body))
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	shortURL := fmt.Sprintf("%s%s", h.Short.Address, aliasURL)

	return ctx.String(http.StatusCreated, shortURL)

}

func (h *Handlers) getRedirectHandler(ctx echo.Context) error {
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

func (h *Handlers) postURLHandler(ctx echo.Context) error {

	// Получаем значение заголовка Content-Type
	contentType := ctx.Request().Header.Get("Content-Type")

	// Проверяем, что Content-Type равен "text/plain"
	if contentType != "application/json" {
		return ctx.String(http.StatusBadRequest, "Content-Type не соответсвует ожидаемому: application/json")
	}
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Возникал ошибка при чтении тела запроса: "+err.Error())
	}
	defer ctx.Request().Body.Close()
	res := model.PostURLHandlerResponse{}
	req := model.PostURLHandlerRequest{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		res.ErrMsg = "Возникал ошибка при чтении тела запроса: " + err.Error()
		return ctx.JSON(http.StatusBadRequest, res)
	}

	aliasURL, err := h.Short.SetAliasName(string(req.URL))
	if err != nil {
		res.ErrMsg = err.Error()
		return ctx.JSON(http.StatusInternalServerError, res)
	}
	shortURL := fmt.Sprintf("%s%s", h.Short.Address, aliasURL)
	res.Result = shortURL

	return ctx.JSON(http.StatusCreated, res)

}
