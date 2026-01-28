package handler

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const secretKey = "tort-secret-key"

func (h *Handlers) cookieMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("user_id")

		var userID string
		if err != nil {
			// Кука отсутствует, генерируем новый ID
			id, err := generateUniqueID()
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "ошибка генерации userID")
			}
			userID = id
		} else {
			// Проверяем подлинность куки
			id, valid := verifySignedID(cookie.Value)
			if !valid {
				// Подпись недействительна, генерируем новый ID
				id, err := generateUniqueID()
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "ошибка генерации userID")
				}
				userID = id
			} else {
				userID = id
			}
		}

		// Устанавливаем новую куку
		signedID := signID(userID, h.secret)
		c.SetCookie(&http.Cookie{
			Name:     "user_id",
			Value:    signedID,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
		})

		// Сохраняем ID пользователя в контексте
		c.Set("userID", userID)

		return next(c)
	}
}

// generateUniqueID - генерирует уникальный ID
func generateUniqueID() (string, error) {
	uuid := make([]byte, 4)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(uuid), nil
}

// signID - подписывает ID
func signID(id, sk string) string {
	mac := hmac.New(sha256.New, []byte(sk))
	mac.Write([]byte(id))
	expectedMAC := mac.Sum(nil)
	return id + "." + hex.EncodeToString(expectedMAC)
}

// verifySignedID - проверяет подпись ID
func verifySignedID(signedID string) (string, bool) {
	parts := strings.Split(signedID, ".")
	if len(parts) != 2 {
		return "", false
	}

	id, receivedMAC := parts[0], parts[1]

	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(id))
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	return id, hmac.Equal([]byte(receivedMAC), []byte(expectedMAC))
}

// getUserID - извлекает ID пользователя из куки
func (h *Handlers) getUserID(c echo.Context) (string, error) {
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
