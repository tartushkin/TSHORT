package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/tartushkin/TSHORT.git/internal/service"
)

type Handlers struct {
	Short      *service.Short // внутриняя логика приложения
	httpServer *echo.Echo
}

func NewHandlers(short *service.Short) *Handlers {
	return &Handlers{Short: short}
}

// StartHTTP - инициализация и запуск сервера
func (h *Handlers) StartHTTP(ctx context.Context, httpPort string) error {
	h.httpServer = echo.New()
	h.httpServer.Use(
		middleware.Logger(),  // Логирование
		middleware.Recover(), // Обработка паник
		avoidGzipOnLocation,  // Наш новый фильтр
		middleware.GzipWithConfig(middleware.GzipConfig{
			Level:     5,
			MinLength: 15,
		}), // Последний шаг - GZIP
	)

	h.httpServer.POST("/", h.oldPostURLHandler)
	h.httpServer.GET("/:id", h.getRedirectHandler)
	h.httpServer.POST("/api/shorten", h.postURLHandler)

	h.httpServer.Logger.Fatal(h.httpServer.Start(httpPort))

	return nil
}

// остнавка http сервера
func (h *Handlers) StopHTTP(ctx context.Context) {
	h.httpServer.Shutdown(ctx)
}

func avoidGzipOnLocation(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Проверяем, есть ли заголовок Location
		loc := c.Response().Header().Get("Location")
		if loc != "" {
			// Возвращаемся сразу, не позволяя middleware сжимать Location
			return next(c)
		}

		// Остальные заголовки оставляем неизменными
		return next(c)
	}
}
