package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
	cfg "github.com/tartushkin/TSHORT.git/internal/config/app"
	"github.com/tartushkin/TSHORT.git/internal/handler"
	sr "github.com/tartushkin/TSHORT.git/internal/service"
)

func main() {
	lg := logrus.New()
	//ctx := context.Background()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	cfg := cfg.NewConfig()             // инициализация конфига
	sh, err := sr.Create(ctx, lg, cfg) // инициализация сервиса
	if err != nil {
		panic(err)
	}
	defer sh.Close()

	h := handler.NewHandlers(sh)
	//serverDone := make(chan struct{})
	go func() {
		if err := h.StartHTTP(ctx, cfg.Port, cfg.SecretKey); err != nil && err != http.ErrServerClosed {
			lg.Error("ошибка HTTP-сервера", "error", err)
			stop()
		}
	}()
	lg.Info("HTTP-сервер запущен", "port", cfg.Port)

	// Ждём сигнала остановки
	<-ctx.Done()
	lg.Info("получен сигнал остановки, завершаем работу...")
	//time.Sleep(100 * time.Millisecond)
	if cfg.RunProfile {
		fmt.Println("зашли")
		pprof.StopCPUProfile()
		sh.Fcpu.Close()

		err := sh.MemProfile()
		if err != nil {
			panic(err)
		}
		sh.Fmem.Close()
	}
	time.Sleep(time.Second)
}
