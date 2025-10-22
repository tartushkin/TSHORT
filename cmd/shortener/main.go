package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/tartushkin/TSHORT.git/internal/handler"
	sr "github.com/tartushkin/TSHORT.git/internal/service"
)

const (
	httpPort = 8080
)

func main() {
	lg := logrus.New()

	ctx := context.Background()

	sh := sr.Create(ctx, lg)
	h := handler.NewHandlers(sh)

	go h.StartHTTP(ctx, httpPort)
	lg.Info("Listner", "Запущен http слушатель на порту ", httpPort)

	go func() {
		<-ctx.Done()
		h.StopHTTP(ctx)
	}()
	select {}
}
