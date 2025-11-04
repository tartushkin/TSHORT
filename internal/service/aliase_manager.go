package service

import (
	"encoding/hex"
	"fmt"

	"github.com/tartushkin/TSHORT.git/internal/model"
)

// SetAliasName - формирование сокращенного url
func (s *Short) SetAliasName(url string) (string, error) {
	aliasURL := hex.EncodeToString([]byte(url))
	aliasURL = aliasURL[:8]
	s.mu.RLock()
	s.CacheURL[aliasURL] = url
	s.mu.RUnlock()
	newURL := model.StorageURL{
		Alias:    aliasURL,
		Original: url,
	}
	err := s.write(&newURL)
	if err != nil {
		return "", fmt.Errorf("возникла ошибка: %w при записи в файл новую пару URL", err)
	}
	return aliasURL, nil
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

func (s *Short) write(event *model.StorageURL) error {
	err := s.File.Encoder.Encode(event)
	if err != nil {
		return err
	}
	return nil
}
