package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
	cfg "github.com/tartushkin/TSHORT.git/internal/config/app"
	"github.com/tartushkin/TSHORT.git/internal/handler"
	sr "github.com/tartushkin/TSHORT.git/internal/service"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

/*
main запускает HTTP-сервер для сокращения URL.

Программа:
 1. Инициализирует логгер, конфигурацию и сервис.
 2. Запускает HTTP-сервер в отдельной горутине.
 3. Ожидает сигнала завершения (Ctrl+C, SIGTERM).
 4. При получении сигнала — корректно останавливается.
 5. Если включено профилирование (RunProfile), останавливает CPU- и memory-профили.

Зависимости:
  - Логирование: github.com/sirupsen/logrus
  - Конфигурация: internal/config/app
  - Сервис: internal/service
  - Хендлеры: internal/handler
*/
func main() {
	buildInfo()
	lg := logrus.New()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()
	cfg := cfg.NewConfig(lg)                // инициализация конфига
	sh, dis, err := sr.Create(ctx, lg, cfg) // инициализация сервиса
	if err != nil {
		lg.Error("ошибка инициализация сервиса", "error", err)
		return
	}
	defer sh.Close()

	h := handler.NewHandlers(dis, sh, cfg.SubNet)
	go func() {
		if err := h.StartHTTP(ctx, cfg); err != nil && err != http.ErrServerClosed {
			lg.Error("ошибка HTTP-сервера", "error", err)
			cancel()
		}
	}()
	lg.Info("HTTP-сервер запущен", "port", cfg.Port)

	// Ждём сигнала остановки
	<-ctx.Done()
	lg.Info("получен сигнал остановки, завершаем работу...")
	if cfg.RunProfile {
		pprof.StopCPUProfile()
		sh.Fcpu.Close()

		err := sh.MemProfile()
		if err != nil {
			lg.Error("ошибка профилирования памяти", "error", err)
			//cancel()
		}
		sh.Fmem.Close()
	}
}

func buildInfo() {
	// Выводим информацию о сборке
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
