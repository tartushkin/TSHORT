package service

import (
	"encoding/hex"
	"fmt"
)

// SetAliasName - формирование сокращенного url
func (s *Short) SetAliasName(url string) string {
	aliasURL := hex.EncodeToString([]byte(url))
	aliasURL = aliasURL[:8]
	s.mu.RLock()
	s.CacheURL[aliasURL] = url
	s.mu.RUnlock()
	return aliasURL
}

// GetAliasName - получение оригинального url
func (s *Short) GetAliasName(aliasURL string) (string, error) {
	s.mu.RLock()
	value, ok := s.CacheURL[aliasURL]
	s.mu.RUnlock()
	if !ok {
		return "", fmt.Errorf("не удалось найти оригинальный url по сокращенному: %s", aliasURL)
	}
	return value, nil
}
