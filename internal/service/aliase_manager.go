package service

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/tartushkin/TSHORT.git/internal/model"
)

// SetCouple - добавление пары сокращенный/оригинальный url в кеш
func (s *Short) SetCouple(coupe *model.AliasFullCore) (string, error) {
	s.mu.RLock()
	s.CacheURL[coupe.Alias] = coupe
	s.mu.RUnlock()
	URL := fmt.Sprintf("%s/%s", s.Address, coupe.Alias) //URL := fmt.Sprintf("%s/%s", s.Address, coupe.Alias)
	return URL, nil
}

// GetAliasName - получение оригинального url
func (s *Short) GetAliasName(aliasURL string) (string, error) {
	s.mu.RLock()
	coupe, ok := s.CacheURL[aliasURL]
	s.mu.RUnlock()
	if !ok {
		return "", fmt.Errorf("не удалось найти оригинальный url по сокращенному: %s", aliasURL)
	}

	return coupe.OriginalURL, nil
}

func (s *Short) write(event *model.AliasFullCore) error {
	err := s.File.Encoder.Encode(event)
	if err != nil {
		return err
	}
	return nil
}

// проверка наличия url в кеше
func (s *Short) checkURL(outURL string) error {
	for _, coupe := range s.CacheURL {
		if coupe.OriginalURL == outURL {
			return fmt.Errorf("данный URL - %s уже есть в кеше приложения по ключу: %s", outURL, coupe.Alias)
		}
	}
	return nil
}

func (s *Short) checkSourse() string {
	if s.DNS != "" {
		return model.DATABASE
	}
	if s.File != nil {
		return model.FILE
	}
	return model.Cache
}

func (s *Short) insertURL(listURL []*model.AliasFullCore, req string) error {
	sourse := s.checkSourse()
	switch sourse {
	case model.DATABASE:
		s.Logger.Info("insertURL - хранилище для данных: " + model.DATABASE)
		if req == model.Text || req == model.One {
			couple := listURL[0]
			err := s.Repo.InsertURL(s.Ctx, couple)
			if err != nil {
				if err.Error() == model.ERRCONFLICT {
					alias, err := s.Repo.GetAlias(s.Ctx, couple.OriginalURL)
					if err != nil {
						return err
					}
					return fmt.Errorf("данный URL - %v уже есть в БД приложения по ключу: %v", couple.OriginalURL, alias)
				}
				return fmt.Errorf("возникла ошибка: %w при записи в БД новую пару URL", err)
			}
		} else {
			list, err := json.Marshal(listURL)
			if err != nil {
				return err
			}
			err = s.Repo.InsertURLJson(s.Ctx, list)
			if err != nil {
				return fmt.Errorf("возникла ошибка: %w при записи в БД новую пару URL", err)
			}
		}
	case model.FILE:
		s.Logger.Info("insertURL - хранилище для данных: " + model.FILE)
		for _, couple := range listURL {
			s.Logger.Info(fmt.Sprintf("insertURL - запись в файл: %v  ", couple))
			err := s.write(couple)
			if err != nil {
				return fmt.Errorf("возникла ошибка: %w при записи в файл новую пару URL", err)
			}
		}

	}
	return nil
}
func (s *Short) getUUID() string {
	uuidURL := uuid.New()
	aliasURL := uuidURL.String()

	aliasURL = aliasURL[:8]
	return aliasURL
}
