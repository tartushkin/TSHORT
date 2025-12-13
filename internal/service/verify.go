package service

import (
	"strings"

	"github.com/labstack/echo/v4"
)

// getUserID - извлекает ID пользователя из куки
func (s *Short) GetUserID(c echo.Context) (string, error) {

	if c.Request().URL.Path != "/api/user/urls" && c.Request().URL.Path != "/api/shorten/batch" {
		return "testID", nil
	}

	cookie, err := c.Cookie("user_id")
	if err != nil {
		return "", err
	}

	// Разбиваем значение куки на части
	parts := strings.Split(cookie.Value, ".")
	if len(parts) != 2 {
		return "", err
	}

	// Извлекаем user_id
	userID := parts[0]
	return userID, nil
}
