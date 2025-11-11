package handler

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/tartushkin/TSHORT.git/internal/service"
)

type Handlers struct {
	Short      *service.Short // внутриняя логика приложения
	httpServer *echo.Echo
}
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func NewHandlers(short *service.Short) *Handlers {
	return &Handlers{Short: short}
}

// StartHTTP - инициализация и запуск сервера
func (h *Handlers) StartHTTP(ctx context.Context, httpPort string) error {
	h.httpServer = echo.New()
	h.httpServer.Use(middleware.Logger()) //в билиотеке уже есть middleware для логирования запрсов
	h.httpServer.Use(middleware.Recover())
	h.httpServer.Use(GzipMiddleware)

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

func GzipMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// проверяем что это не редирект
		if c.Response().Status == http.StatusTemporaryRedirect ||
			c.Response().Status == http.StatusPermanentRedirect {
			return next(c)
		}

		gz, err := gzip.NewWriterLevel(c.Response().Writer, gzip.BestSpeed)
		if err != nil {
			return err
		}
		defer gz.Close()

		// обёртываем оригинальный ResponseWriter
		gzw := gzipResponseWriter{Writer: gz, ResponseWriter: c.Response().Writer}

		// Заменяем ResponseWriter в контексте Echo
		c.Response().Writer = gzw
		c.Response().Header().Set(echo.HeaderContentEncoding, "gzip")

		return next(c)
	}
}

func (w gzipResponseWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w gzipResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}
