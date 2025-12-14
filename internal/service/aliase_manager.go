package service

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/tartushkin/TSHORT.git/internal/model"
)

// SetCouple - добавление пары сокращенный/оригинальный url в кеш
func (s *Short) setCouple(coupe *model.AliasFullCore) {
	s.mu.RLock()
	s.CacheURL[coupe.Alias] = coupe
	s.mu.RUnlock()
	//URL := fmt.Sprintf("%s/%s", s.Address, coupe.Alias) //URL := fmt.Sprintf("%s/%s", s.Address, coupe.Alias)
}

func (s *Short) getFull(alias string) string {
	URL := fmt.Sprintf("%s/%s", s.Address, alias)
	return URL
}

// GetAliasName - получение оригинального url
func (s *Short) GetAliasName(aliasURL string) (string, error) {
	s.mu.RLock()
	coupe, ok := s.CacheURL[aliasURL]
	s.mu.RUnlock()
	if !ok {
		orig, err := s.Repo.GetOriginalURL(s.Ctx, aliasURL)
		if err != nil {
			return "", fmt.Errorf("не удалось найти оригинальный url по сокращенному: "+aliasURL+". ERR - %s", err.Error)
		}
		return orig, nil
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
func (s *Short) checkURL(outURL string) (string, bool) {
	for _, coupe := range s.CacheURL {
		if coupe.OriginalURL == outURL {
			s.Logger.Info(fmt.Sprintf("данный URL - %s уже есть в кеше приложения по ключу: %s", outURL, coupe.Alias))
			return coupe.Alias, true
		}
	}
	return "", false
}

func (s *Short) checkSourse() string {
	if s.conn != nil {
		return model.DATABASE
	}
	if s.File != nil {
		return model.FILE
	}
	return model.CACHE
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
					errMsg := fmt.Errorf("%s -%s/%s", model.ERRCONFLICT, s.Address, alias)
					s.Logger.Error(fmt.Errorf("%s - данный URL - %v уже есть в БД приложения по ключу: %v", model.ERRCONFLICT, couple.OriginalURL, alias))
					return errMsg
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
			s.setCouple(couple)
		}
	case model.CACHE:
		s.Logger.Info("insertURL - хранилище для данных: " + model.CACHE)
		for _, couple := range listURL {
			s.Logger.Info(fmt.Sprintf("insertURL - запись в кеш: %v  ", couple))
			s.setCouple(couple)
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
