package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tartushkin/TSHORT.git/internal/model"
)

func (h *Handlers) oldPostURLHandler(ctx echo.Context) error {
	var shortURL string
	body, err := h.getBody(ctx)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	userID, err := h.getUserID(ctx)
	if err != nil {
		h.Short.Logger.Error("ошибка при работе с телом запроса: " + err.Error())
		return ctx.String(http.StatusUnauthorized, err.Error())
	}
	list, err := h.Short.ReaderBody(body, userID, model.Text)
	if err != nil {
		if strings.HasPrefix(err.Error(), model.ERRCONFLICT) {
			parts := strings.Split(err.Error(), "-")
			return ctx.String(http.StatusConflict, parts[1])
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
		h.Short.Logger.Error("getRedirectHandler.err - отсутствует алиас")
		return ctx.JSON(http.StatusBadRequest, "Требуется алиас")
	}
	h.Short.Logger.Info("getRedirectHandler.info - полученный алиас: " + alias)
	// Извлекаем алиас
	originalURL, err := h.Short.GetAliasName(alias)
	if err != nil {
		h.Short.Logger.Error("getRedirectHandler.err - возникла ошбка при получениии оригинального URL: " + err.Error())
		if strings.Contains(err.Error(), "URL deleted") {
			return ctx.NoContent(http.StatusGone)
			//return ctx.JSON(http.StatusGone, err.Error())
		}
		return ctx.JSON(http.StatusNotFound, err.Error())
	}
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
		return ctx.JSON(http.StatusBadRequest, "Content-Type не соответсвует ожидаемому: application/json")
	}
	body, err := h.getBody(ctx)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	userID, err := h.getUserID(ctx)
	if err != nil {
		h.Short.Logger.Error("ошибка при работе с телом запроса: " + err.Error())
		return ctx.String(http.StatusUnauthorized, err.Error())
	}
	defer ctx.Request().Body.Close()
	res := model.PostURLHandlerResponse{}
	listURL, err := h.Short.ReaderBody(body, userID, model.One)
	if err != nil {
		res.ErrMsg = err.Error()
		if strings.HasPrefix(res.ErrMsg, model.ERRCONFLICT) {
			return ctx.JSON(http.StatusConflict, res)
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
	body, err := h.getBody(ctx)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	userID, err := h.getUserID(ctx)
	if err != nil {
		h.Short.Logger.Error("ошибка при работе с телом запроса: " + err.Error())
		return ctx.String(http.StatusUnauthorized, err.Error())
	}
	defer ctx.Request().Body.Close()

	listURL, err := h.Short.ReaderBody(body, userID, model.List)
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
	userID, err := h.getUserID(ctx)
	if err != nil {
		h.Short.Logger.Error("ошибка при работе с телом запроса: " + err.Error())
		return ctx.String(http.StatusUnauthorized, err.Error())
	}
	userIDList, err := h.Short.GetUserURL(userID)
	if err != nil {
		if he, ok := err.(*echo.HTTPError); ok {
			return he
		}
		h.Short.Logger.Error("Ошибка при работе с телом запроса: " + err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
		//return ctx.JSON(http.StatusNoContent, err.Error())
	}
	if len(userID) == 0 {
		return ctx.NoContent(http.StatusNoContent)
	}
	h.Short.Logger.Info("Возвращаем список URL: ", userID)
	return ctx.JSON(http.StatusOK, userIDList)

}

func (h *Handlers) deleteURL(ctx echo.Context) error {

	deleteList := []string{}
	if err := ctx.Bind(&deleteList); err != nil {
		h.Short.Logger.Error("ошибка при работе с телом запроса: " + err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "error:"+err.Error())
	}
	userID, err := h.getUserID(ctx)
	if err != nil {
		h.Short.Logger.Error("ошибка при работе с телом запроса: " + err.Error())
		return ctx.String(http.StatusUnauthorized, err.Error())
	}
	err = h.Short.DeleteUserURL(userID, deleteList)
	if err != nil {
		h.Short.Logger.Error("Ошибка при отметке url на удаление: " + err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusAccepted)
}
