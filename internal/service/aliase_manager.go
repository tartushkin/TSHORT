package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tartushkin/TSHORT.git/internal/model"
)

// SetAliasName - формирование сокращенного url
func (s *Short) SetAliasName(url string) (string, error) {
	err := s.checkURL(url)
	if err != nil {
		return "", err
	}
	uuidURL := uuid.New()
	aliasURL := uuidURL.String()

	aliasURL = aliasURL[:8]
	s.mu.RLock()
	s.CacheURL[aliasURL] = url
	s.mu.RUnlock()
	newURL := model.StorageURL{
		Alias:    aliasURL,
		Original: url,
	}
	sourse := s.checkSourse()
	switch sourse {
	case model.DATABASE:
		err := s.Repo.InsertURL(s.Ctx, url, aliasURL)
		if err != nil {
			return "", fmt.Errorf("возникла ошибка: %w при записи в БД новую пару URL", err)
		}
	case model.FILE:
		err = s.write(&newURL)
		if err != nil {
			return "", fmt.Errorf("возникла ошибка: %w при записи в файл новую пару URL", err)
		}
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

// проверка наличия url в кеше
func (s *Short) checkURL(outURL string) error {
	for k, v := range s.CacheURL {
		if v == outURL {
			return fmt.Errorf("данный URL - %s уже есть в кеше приложения по ключу: %s", outURL, k)
		}
	}
	return nil
}

func (s *Short) checkSourse() string {
	if s.DNS != "" {
		return model.DATABASE
	}
	return model.FILE
}
