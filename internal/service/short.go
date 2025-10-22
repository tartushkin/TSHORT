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

// SetAliaseName - формирование сокращенного url
func (s *Short) SetAliaseName(url string) string {
	aliaseURL := hex.EncodeToString([]byte(url))
	aliaseURL = aliaseURL[:8]
	s.CacheURL[aliaseURL] = url
	return aliaseURL
}

// GetAliaseName - получение оригинального url
func (s *Short) GetAliasName(aliaseURL string) (string, error) {
	value, ok := s.CacheURL[aliaseURL]
	if !ok {
		return "", fmt.Errorf("не удалось найти оригинальный url по сокращенному: %s", aliaseURL)
	}
	return value, nil
}
