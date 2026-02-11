package service

import (
	"strings"

	"github.com/labstack/echo/v4"
)

// GetUser - извлекает ID пользователя из куки
func (s *Short) GetUser(c echo.Context) (string, error) {
	if userID, ok := c.Get("userID").(string); ok {
		return userID, nil
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
