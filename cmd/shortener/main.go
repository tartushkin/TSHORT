package main

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
	cfg "github.com/tartushkin/TSHORT.git/internal/config/app"
	"github.com/tartushkin/TSHORT.git/internal/handler"
	sr "github.com/tartushkin/TSHORT.git/internal/service"
)

func main() {
	lg := logrus.New()
	//ALTER TABLE t_short.t_list ADD CONSTRAINT IF NOT EXISTS unique_s_full UNIQUE (s_full);
	ctx := context.Background()
	cfg := cfg.NewConfig()             // инициализация конфига
	sh, err := sr.Create(ctx, lg, cfg) // инициализация сервиса
	if err != nil {
		panic(err)
	}
	defer sh.Close()

	h := handler.NewHandlers(sh)
	go h.StartHTTP(ctx, cfg.Port) // запуск сервера
	lg.Info("Listner: ", fmt.Sprintf("Запущен http слушатель на порту %s", cfg.Port))

	go func() {
		<-ctx.Done()
		h.StopHTTP(ctx)

	}()
	select {}
}
