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
	h.httpServer.Use(middleware.Logger()) //в билиотеке уже есть middleware для логирования запрсов
	h.httpServer.Use(middleware.Recover())
	h.httpServer.Use(middleware.Gzip()) //в билиотеке уже есть middleware для сжатия

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
