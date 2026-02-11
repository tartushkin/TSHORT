// Package service реализует бизнес-логику приложения: сокращение URL,
// работа с кешем, хранилищем.

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/tartushkin/TSHORT.git/internal/model"
)

// setCouple - добавление пары сокращенный/оригинальный url в кеш.
func (s *Short) setCouple(coupe *model.AliasFullCore) {
	s.mu.RLock()
	s.CacheURL[coupe.Alias] = coupe
	s.mu.RUnlock()

}

// GetOriginalURL возвращает оригинальный URL по сокращённому.
func (s *Short) GetOriginalURL(shortURL string) string {
	parts := strings.Split(shortURL, "/")
	alias := parts[len(parts)-1]

	if s.Repo != nil {
		org, err := s.Repo.GetOriginalURL(s.Ctx, alias)
		if err != nil {
			s.Logger.Info("GetOriginalURL.err - не удалось найтти оригинальный URL для - " + alias)
			return ""
		}
		return org.OriginalURL
	}
	return ""
}

func (s *Short) getFull(alias string) string {
	URL := fmt.Sprintf("%s/%s", s.Address, alias)
	return URL
}

// GetAliasName - возвращает оригинальный URL по alias.
// Возвращает ошибку, если URL помечен как удалённый.
func (s *Short) GetAliasName(aliasURL string) (string, error) {
	var couple *model.AliasFullCore
	s.mu.RLock()
	couple, ok := s.CacheURL[aliasURL]
	s.mu.RUnlock()

	if !ok {
		var err error
		couple, err = s.Repo.GetOriginalURL(s.Ctx, aliasURL)
		if err != nil {
			return "", fmt.Errorf("не удалось найти оригинальный url по сокращенному: "+aliasURL+". ERR - %s", err.Error())
		}
	}
	if couple.DeletedFlag {
		return "", fmt.Errorf("URL deleted: %s", aliasURL)
	}
	return couple.OriginalURL, nil
}

func (s *Short) write(event *model.AliasFullCore) error {
	err := s.FileStorage.Encoder.Encode(event)
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

// checkSourse - определяет тип хранилища: DATABASE, FILE или CACHE.
func (s *Short) checkSourse() string {
	if s.Conn != nil {
		return model.DATABASE
	}
	if s.FileStorage != nil {
		return model.FILE
	}
	return model.CACHE
}

// insertURL -  сохраняет список URL в выбранное хранилище.
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

// getUUID - генерирует короткий alias (8 символов).
func (s *Short) getUUID() string {
	uuidURL := uuid.New()
	aliasURL := uuidURL.String()

	aliasURL = aliasURL[:8]
	return aliasURL
}

// DeleteMarkedURLs - удаляет помеченные URL из БД и кеша.
func (s *Short) DeleteMarkedURLs(ctx context.Context) error {

	err := s.Repo.DeleteMarkedURLs(ctx)
	if err != nil {
		return err
	}
	s.Logger.Info("DeleteMarkedURLs.complete - успешное удаление URL из БД")
	for _, couple := range s.CacheURL {
		if couple.DeletedFlag {
			delete(s.CacheURL, couple.Alias)
		}
	}
	s.Logger.Info("DeleteMarkedURLs.complete - успешное удаление URL из кеша")
	return nil
}
