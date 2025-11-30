package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tartushkin/TSHORT.git/internal/model"
)

func (s *Short) ReaderBody(ctx echo.Context, req string) ([]*model.BranchResponse, error) {
	s.Logger.Info("ReaderBody.start - чтение тела запроса")
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return nil, ctx.String(http.StatusBadRequest, "Возникал ошибка при чтении тела запроса: "+err.Error())
	}

	defer ctx.Request().Body.Close()
	list := []*model.AliasFullCore{}
	responseList := []*model.BranchResponse{}
	var originalURL string

	switch req {
	case model.One:
		couple := model.PostURLHandlerRequest{}
		err = json.Unmarshal(body, &couple)
		if err != nil {
			errMsg := fmt.Errorf("возникла ошибка при чтении тела запроса: %w", err)
			return nil, errMsg
		}
		if couple.URL == "" {
			return nil, fmt.Errorf("необходимо заполнить поле: URL")
		}
		s.Logger.Info("ReaderBody.json - запись одного URL: " + couple.URL)

		aliasURL, fullURL, err := s.setAliasName(couple.URL, "")
		if err != nil {
			return nil, err
		}
		list = append(list, &model.AliasFullCore{
			Alias:       aliasURL,
			OriginalURL: couple.URL,
		})
		err = s.insertURL(list)
		if err != nil {
			return nil, err
		}

		responseList = append(responseList, &model.BranchResponse{ShortUrl: fullURL})
		s.Logger.Info("ReaderBody.json - успешно отправили пару в хранилище, алиас: " + aliasURL)
	case model.List:
		coupleList := []*model.BranchRequest{}
		s.Logger.Info("ReaderBody.jsonList - запись списка URL")
		err = json.Unmarshal(body, &coupleList)
		if err != nil {
			errMsg := fmt.Errorf("возникла ошибка при чтении тела запроса: %w", err)
			return nil, errMsg
		}
		for _, couple := range coupleList {
			if couple.CorrID == "" {
				return nil, fmt.Errorf("необходимо заполнить поле: CorrID")
			}
			if couple.OriginalURL == "" {
				return nil, fmt.Errorf("необходимо заполнить поле: OriginalURL")
			}
			aliasURL, fullURL, err := s.setAliasName(couple.OriginalURL, couple.CorrID)
			if err != nil {
				return nil, err
			}
			fix := model.AliasFullCore{
				Alias:       aliasURL,
				OriginalURL: couple.OriginalURL,
				CorrID:      couple.CorrID,
			}
			list = append(list, &fix)
			responseList = append(responseList, &model.BranchResponse{ShortUrl: fullURL, CorrID: couple.CorrID})
			s.Logger.Info(fmt.Sprintf("ReaderBody.jsonList - запись пары в список для отправки в хранилище: %s/%s/%s ",
				aliasURL, couple.OriginalURL, couple.CorrID))
		}
		err = s.insertURL(list)
		if err != nil {
			errMsg := fmt.Errorf("возникла ошибка при записи спика URL в хранилище: %w", err)
			return nil, errMsg
		}
		s.Logger.Info("ReaderBody.jsonList - успешно отправили список в хранилище.")
	default:
		s.Logger.Info("ReaderBody.text - запись URL текстом")
		originalURL = string(body)
		if originalURL == "" {
			return nil, fmt.Errorf("необходимо заполнить тело запроса")
		}

		aliasURL, fullURL, err := s.setAliasName(originalURL, "")
		if err != nil {
			return nil, err
		}
		list = append(list, &model.AliasFullCore{
			Alias:       aliasURL,
			OriginalURL: originalURL,
		})

		err = s.insertURL(list)
		if err != nil {
			errMsg := fmt.Errorf("возникла ошибка при записи URL в хранилище: %w", err)
			return nil, errMsg
		}
		responseList = append(responseList, &model.BranchResponse{ShortUrl: fullURL})
		s.Logger.Info("ReaderBody.text - успешно отправили URL в хранилище.")
	}

	return responseList, nil
}

func (s *Short) setAliasName(originalURL, corrID string) (string, string, error) {
	coupe := &model.AliasFullCore{
		OriginalURL: originalURL,
		CorrID:      corrID,
	}
	aliasURL, err := s.SetAliasName(coupe)
	if err != nil {
		return "", "", err
	}
	URL := fmt.Sprintf("%s/%s", s.Address, aliasURL)
	return aliasURL, URL, nil
}
