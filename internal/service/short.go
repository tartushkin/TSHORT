package service

import (
	"context"

	"github.com/sirupsen/logrus"
	cfg "github.com/tartushkin/TSHORT.git/cmd/config"
)

type Short struct {
	Logger   *logrus.Logger
	Ctx      context.Context
	HTTPPort string
	Address  string

	CacheURL map[string]string
}

// NewShort - заполнение структуры приложения
func Create(ctx context.Context, lg *logrus.Logger, cfg *cfg.Configure) *Short {
	cacheURL := map[string]string{}

	sh := &Short{
		Logger:   lg,
		Ctx:      ctx,
		CacheURL: cacheURL,
		HTTPPort: cfg.Port,
	}
	sh.Address = cfg.Address + "/"
	return sh
}
