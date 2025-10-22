package service

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/sirupsen/logrus"
)
//
type Short struct {
	Logger   *logrus.Logger
	Ctx      context.Context
	CacheUrl map[string]string
}

// NewShort - заполнение структуры приложения
func Create(ctx context.Context, lg *logrus.Logger) *Short {
	cacheUrl := map[string]string{}

	sh := &Short{
		Logger:   lg,
		Ctx:      ctx,
		CacheUrl: cacheUrl,
	}
	return sh
}

// SetAliaseName - формирование сокращенного url
func (s *Short) SetAliaseName(url string) string {
	aliaseUrl := hex.EncodeToString([]byte(url))
	aliaseUrl = aliaseUrl[:8]
	s.CacheUrl[aliaseUrl] = url
	return aliaseUrl
}

// GetAliaseName - получение оригинального url
func (s *Short) GetAliasName(aliaseUrl string) (string, error) {
	value, ok := s.CacheUrl[aliaseUrl]
	if !ok {
		return "", fmt.Errorf("не удалось найти оригинальный url по сокращенному: %s", aliaseUrl)
	}
	return value, nil
}
