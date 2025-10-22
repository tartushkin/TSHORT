package service

import (
	"encoding/hex"
	"fmt"
)

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
