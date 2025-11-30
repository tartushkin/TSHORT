package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tartushkin/TSHORT.git/internal/model"
)

func (s *Short) ReaderBody(ctx echo.Context, jn bool) (string, error) {
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return "", ctx.String(http.StatusBadRequest, "Возникал ошибка при чтении тела запроса: "+err.Error())
	}
	defer ctx.Request().Body.Close()
	var originalURL string
	originalURL = string(body)
	if jn {
		req := model.PostURLHandlerRequest{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			errMsg := fmt.Errorf("возникла ошибка при чтении тела запроса: %w", err)
			return "", errMsg
		}
		originalURL = req.URL
	}
	aliasURL, err := s.SetAliasName(originalURL)
	if err != nil {
		return "", err
	}
	URL := fmt.Sprintf("%s/%s", s.Address, aliasURL)
	return URL, nil
}
