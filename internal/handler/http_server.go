package handler

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/tartushkin/TSHORT.git/internal/service"
)

type Handlers struct {
	Short      *service.Short // внутриняя логика приложения
	httpServer *echo.Echo
}

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
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
	h.httpServer.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return cookieMiddleware(next)
	})
	h.httpServer.HTTPErrorHandler = func(err error, c echo.Context) {
		// Логируем ошибку
		c.Logger().Error(err)

		// Определяем код и сообщение
		code := http.StatusInternalServerError
		message := "Internal Server Error"

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = fmt.Sprintf("%v", he.Message)
		}

		// Отправляем JSON
		if !c.Response().Committed {
			c.JSON(code, map[string]string{
				"error": message,
			})
		}
	}

	h.httpServer.POST("/", h.oldPostURLHandler)
	h.httpServer.GET("/:id", h.getRedirectHandler)
	h.httpServer.POST("/api/shorten", h.postURLHandler)
	h.httpServer.GET("/ping", h.testConnectionDB)
	h.httpServer.POST("/api/shorten/batch", h.batchHandler)
	h.httpServer.POST("/api/user/urls", h.getMyShortURL)
	h.httpServer.DELETE("/api/user/urls", h.deleteURL)

	if err := h.httpServer.Start(httpPort); err != nil && err != http.ErrServerClosed {
		h.httpServer.Logger.Info("Сервер остановлен: %v", err)
	}

	return nil
}

// остнавка http сервера
func (h *Handlers) StopHTTP(ctx context.Context) {
	h.httpServer.Shutdown(ctx)
}
func GzipMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Проверяем, что клиент поддерживает сжатие ответа
		acceptEncoding := c.Request().Header.Get(echo.HeaderAcceptEncoding)
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		// Проверяем, что клиент отправил сжатые данные
		contentEncoding := c.Request().Header.Get(echo.HeaderContentEncoding)
		sendsGzip := strings.Contains(contentEncoding, "gzip")

		// Если клиент отправил сжатые данные, оборачиваем тело запроса
		if sendsGzip {
			cr, err := newCompressReader(c.Request().Body)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			c.Request().Body = cr
			defer cr.Close()
		}

		// Если клиент поддерживает сжатие ответа, оборачиваем ResponseWriter
		if supportsGzip {
			res := c.Response()
			cw := newCompressWriter(res.Writer)
			res.Writer = cw
			defer cw.Close()
			res.Header().Set(echo.HeaderContentEncoding, "gzip")
		}

		// Передаём управление следующему обработчику
		return next(c)
	}
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
