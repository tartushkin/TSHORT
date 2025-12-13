package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tartushkin/TSHORT.git/internal/model"
)

func (h *Handlers) oldPostURLHandler(ctx echo.Context) error {
	var shortURL string
	list, err := h.Short.ReaderBody(ctx, model.Text)
	if err != nil {
		if strings.HasPrefix(err.Error(), model.ERRCONFLICT) {
			parts := strings.Split(err.Error(), "-")
			return ctx.JSON(http.StatusConflict, parts[1])
		}
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	for _, couple := range list {
		shortURL = couple.ShortURL
	}
	return ctx.String(http.StatusCreated, shortURL)

}

func (h *Handlers) getRedirectHandler(ctx echo.Context) error {
	// Получаем URL из параметров запроса
	alias := ctx.Param("id")
	if alias == "" {
		fmt.Println("1")
		return ctx.String(http.StatusBadRequest, "Требуется алиас")
	}
	fmt.Println("2")
	// Извлекаем алиас
	originalURL, err := h.Short.GetAliasName(alias)
	fmt.Println("че тут", originalURL)
	if err != nil {
		fmt.Println("3")
		return ctx.String(http.StatusNotFound, err.Error())
	}
	fmt.Println("originalURL тут -", originalURL)
	h.Short.Logger.Info("HTTP.Response - возвращаем полный URL по алиасу: " + alias + " - " + originalURL)
	ctx.Redirect(http.StatusTemporaryRedirect, originalURL)
	for key, values := range ctx.Response().Header() {
		for _, value := range values {
			h.Short.Logger.Info("HTTP.headers - " + key + ":" + value)
		}
	}
	return nil
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
		if strings.HasPrefix(res.ErrMsg, model.ERRCONFLICT) {
			parts := strings.Split(err.Error(), "-")
			return ctx.JSON(http.StatusConflict, parts[1])
		}
		return ctx.JSON(http.StatusInternalServerError, res)
	}
	for _, couple := range listURL {
		res.Result = couple.ShortURL
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

	// Проверяем, что Content-Type равен "application/json"
	if contentType != "application/json" {
		h.Short.Logger.Error("Content-Type не соответсвует ожидаемому: application/json")
		return ctx.JSON(http.StatusBadRequest, "Content-Type не соответсвует ожидаемому: application/json")
	}
	defer ctx.Request().Body.Close()

	listURL, err := h.Short.ReaderBody(ctx, model.List)
	if err != nil {
		h.Short.Logger.Error("Ошибка при работе с телом запроса: " + err.Error())
		if err.Error() == model.CONFLICT {
			return ctx.JSON(http.StatusConflict, err.Error())
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, listURL)
}

func (h *Handlers) getMyShortURL(ctx echo.Context) error {
	// Получаем значение заголовка Content-Type
	contentType := ctx.Request().Header.Get("Content-Type")

	// Проверяем, что Content-Type равен "application/json"
	if contentType != "application/json" {
		return ctx.String(http.StatusBadRequest, "Content-Type не соответсвует ожидаемому: application/json")
	}
	defer ctx.Request().Body.Close()

	userID, err := h.Short.GetUserURL(ctx)
	if err != nil {
		h.Short.Logger.Error("Ошибка при работе с телом запроса: " + err.Error())
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, userID)

}

func (h *Handlers) deleteURL(ctx echo.Context) error {
	// Получаем значение заголовка Content-Type
	contentType := ctx.Request().Header.Get("Content-Type")

	// Проверяем, что Content-Type равен "application/json"
	if contentType != "application/json" {
		return ctx.String(http.StatusBadRequest, "Content-Type не соответсвует ожидаемому: application/json")
	}
	defer ctx.Request().Body.Close()

	deleteList := []string{}
	if err := ctx.Bind(&deleteList); err != nil {
		h.Short.Logger.Error("ошибка при работе с телом запроса: " + err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "error:"+err.Error())
	}
	h.Short.DeleteUserURL(ctx, deleteList)

	return ctx.JSON(http.StatusAccepted, "")
}
