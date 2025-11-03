package service

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Short struct {
	Logger   *logrus.Logger
	Ctx      context.Context
	CacheURL map[string]string
}

// NewShort - заполнение структуры приложения
func Create(ctx context.Context, lg *logrus.Logger) *Short {
	cacheURL := map[string]string{}

	sh := &Short{
		Logger:   lg,
		Ctx:      ctx,
		CacheURL: cacheURL,
	}
	return sh
}
