package handler

import (
	"compress/gzip"
	"context"
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

	h.httpServer.POST("/", h.oldPostURLHandler)
	h.httpServer.GET("/:id", h.getRedirectHandler)
	h.httpServer.POST("/api/shorten", h.postURLHandler)
	h.httpServer.GET("/ping", h.testConnectionDB)

	h.httpServer.Logger.Fatal(h.httpServer.Start(httpPort))

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
