package handler

import (
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
	var shortURL string
	list, err := h.Short.ReaderBody(ctx, model.Text)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	for _, couple := range list {
		shortURL = couple.ShortUrl
	}

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
	h.Short.Logger.Info("HTTP.Response - возвращаем полный URL по алиасу: " + alias + "/" + originalURL)
	res := ctx.Redirect(http.StatusTemporaryRedirect, originalURL)

	for key, values := range ctx.Response().Header() {
		for _, value := range values {
			h.Short.Logger.Info("HTTP.headers - " + key + ":" + value)
		}
	}
	return res
}

func (h *Handlers) postURLHandler(ctx echo.Context) error {

	// Получаем значение заголовка Content-Type
	contentType := ctx.Request().Header.Get("Content-Type")

	// Проверяем, что Content-Type равен "text/plain"
	if contentType != "application/json" {
		return ctx.String(http.StatusBadRequest, "Content-Type не соответсвует ожидаемому: application/json")
	}

	defer ctx.Request().Body.Close()
	res := model.PostURLHandlerResponse{}
	listURL, err := h.Short.ReaderBody(ctx, model.One)
	if err != nil {
		res.ErrMsg = err.Error()
		//if strings.HasPrefix(res.ErrMsg, model.CONFLICT) {
		//	return ctx.JSON(http.StatusConflict, res)
		//}
		return ctx.JSON(http.StatusInternalServerError, res)
	}
	for _, couple := range listURL {
		res.Result = couple.ShortUrl
	}

	return ctx.JSON(http.StatusCreated, res)

}

func (h *Handlers) testConnectionDB(ctx echo.Context) error {
	err := h.Short.Repo.TestConnectionDB()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, "")
}

func (h *Handlers) batchHandler(ctx echo.Context) error {
	// Получаем значение заголовка Content-Type
	contentType := ctx.Request().Header.Get("Content-Type")

	// Проверяем, что Content-Type равен "text/plain"
	if contentType != "application/json" {
		return ctx.String(http.StatusBadRequest, "Content-Type не соответсвует ожидаемому: application/json")
	}
	defer ctx.Request().Body.Close()

	listURL, err := h.Short.ReaderBody(ctx, model.List)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		//return ctx.JSON(http.StatusInternalServerError, err.Error)
	}

	return ctx.JSON(http.StatusOK, listURL)
}
