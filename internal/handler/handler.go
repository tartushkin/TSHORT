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
	ctx.Set("original_url", string(body))
	for _, couple := range list {
		shortURL = couple.ShortURL
	}
	return ctx.String(http.StatusCreated, shortURL)

}

// getRedirectHandler выполняет редирект по короткому ключу.
//
// Если URL помечен как удалённый — возвращает 410 Gone.
// Если не найден — 404 Not Found.
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
	ctx.Set("original_url", originalURL)
	h.Short.Logger.Info("HTTP.Response - возвращаем полный URL по алиасу: " + alias + " - " + originalURL)
	ctx.Redirect(http.StatusTemporaryRedirect, originalURL)
	for key, values := range ctx.Response().Header() {
		for _, value := range values {
			h.Short.Logger.Info("HTTP.headers - " + key + ":" + value)
		}
	}
	return nil
}

// postURLHandler сокращает один URL из JSON.
//
// Ожидает:
//   - Content-Type: application/json
//   - Тело: {"url": "https://example.com"}
//
// Возвращает:
//   - 201 Created: {"result": "http://localhost:8080/abc123"}
//   - 409 Conflict: если URL уже существует
//   - 400/500: при ошибках
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
	origURL := h.Short.GetOriginalURL(res.Result)
	ctx.Set("original_url", origURL)

	return ctx.JSON(http.StatusCreated, res)

}

// testConnectionDB проверяет подключение к базе данных.
//
// Используется для health-check.
//
// Возвращает:
//   - 200 OK: если соединение есть
//   - 500 Internal Server Error: если нет
func (h *Handlers) testConnectionDB(ctx echo.Context) error {
	err := h.Short.Repo.TestConnectionDB()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, "")
}

// batchHandler обрабатывает массовое сокращение URL.
//
// Принимает массив:
//   - [{"correlation_id": "...", "original_url": "..."}, ...]
//
// Возвращает:
//   - 201 Created: [{"correlation_id": "...", "short_url": "..."}, ...]
//   - 400: если Content-Type не application/json
//   - 401: если не авторизован
//   - 500: при внутренних ошибках
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

// getMyShortURL возвращает список всех URL пользователя.
//
// Если список пуст — возвращает 204 No Content.
//
// Возвращает:
//   - 200 OK: [{...}]
//   - 204 No Content: если нет URL
//   - 401: если не авторизован
//   - 500: при внутренней ошибке
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
	if len(userIDList) == 0 {
		return ctx.NoContent(http.StatusNoContent)
	}
	h.Short.Logger.Info("Возвращаем список URL: ", userID)
	return ctx.JSON(http.StatusOK, userIDList)

}

// deleteURL помечает URL на удаление (soft delete).
//
// Принимает:
//   - Массив alias: ["abc123", "def456"]
//
// Возвращает:
//   - 202 Accepted: успешно принято
//   - 400: неверный формат тела
//   - 401: не авторизован
//   - 500: ошибка при обработке
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
		h.Short.Logger.Error("ошибка при отметке url на удаление: " + err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusAccepted)
}

func (h *Handlers) getStats(c echo.Context) error {
	if len(h.subNet) == 0 {
		return echo.NewHTTPError(http.StatusForbidden, "список доверительных подсетей пуст")
	}
	xRealIP := c.Request().Header.Get("X-Real-IP")
	if xRealIP == "" {
		return echo.NewHTTPError(http.StatusForbidden, "заголовок: X-Real-IP - отсутствует")
	}
	// Проверяем, входит ли IP-адрес в доверенную подсеть
	if !h.trustSubNet(xRealIP) {
		return echo.NewHTTPError(http.StatusForbidden, " IP не входит в доверенную подсеть")

	}
	stats, err := h.Short.GetStats()
	if err != nil {
		h.Short.Logger.Error("ошибка при попытке получить статистику: " + err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, stats)
}
