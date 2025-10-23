package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	cfg "github.com/tartushkin/TSHORT.git/cmd/config"
	"github.com/tartushkin/TSHORT.git/internal/handler"
	sr "github.com/tartushkin/TSHORT.git/internal/service"
)

func main() {
	lg := logrus.New()

	ctx := context.Background()
	cfg := cfg.NewConfig()        // инициализация конфига
	sh := sr.Create(ctx, lg, cfg) // инициализация сервиса

	h := handler.NewHandlers(sh)
	go h.StartHTTP(ctx, cfg.Port) // запуск сервера
	lg.Info("Listner", fmt.Sprintf("Запущен http слушатель на порту %s", cfg.Port))

	go func() {
		<-ctx.Done()
		h.StopHTTP(ctx)
	}()
	select {}
}
