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
	list, err := s.getBody(ctx, req) // получаем тело запроса
	if err != nil {

		return nil, err
	}
	for _, couple := range list {
		err := s.checkURL(couple.OriginalURL) //сначала проверяем кеш, потом проверяем наличие в базе
		if err != nil {
			return nil, err
		}
		couple.Alias = s.getUUID()
	}
	err = s.insertURL(list, req)
	if err != nil {
		return nil, err
	}
	responseList := []*model.BranchResponse{}
	for _, couple := range list {
		fullURL, err := s.SetCouple(couple)
		if err != nil {
			return nil, err
		}

		responseList = append(responseList, &model.BranchResponse{ShortURL: fullURL, CorrID: couple.CorrID})
		s.Logger.Info(fmt.Sprintf("ReaderBody.jsonList - запись пары в список для отправки в хранилище: %s/%s/%s ",
			couple.Alias, couple.OriginalURL, couple.CorrID))
	}
	return responseList, nil
}

func (s *Short) getBody(ctx echo.Context, req string) ([]*model.AliasFullCore, error) {

	s.Logger.Info("ReaderBody.start - чтение тела запроса")
	userID, err := s.GetUserID(ctx)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	coupleList := []*model.AliasFullCore{}
	var couple *model.PostURLHandlerRequest
	defer ctx.Request().Body.Close()

	switch req {
	case model.One:
		err = json.Unmarshal(body, &couple)
		if err != nil {
			errMsg := fmt.Errorf("возникла ошибка при чтении тела запроса: %w", err)
			return nil, errMsg
		}
		coupleList = append(coupleList, &model.AliasFullCore{
			OriginalURL: couple.URL,
			CorrID:      "",
			UserID:      userID,
		})
	case model.List:
		s.Logger.Info("ReaderBody.JSON - чтение jsonList запроса")

		err = json.Unmarshal(body, &coupleList)
		if err != nil {
			errMsg := fmt.Errorf("возникла ошибка при чтении тела запроса: %w", err)
			return nil, errMsg
		}
		for _, couple := range coupleList {
			couple.UserID = userID
		}

		return coupleList, nil
	case model.Text:
		s.Logger.Info("ReaderBody.text - чтение текстового URL")

		originalURL := string(body)
		coupleList = append(coupleList, &model.AliasFullCore{
			OriginalURL: originalURL,
			CorrID:      "",
			UserID:      userID,
		})
	}
	return coupleList, nil
}

func (s *Short) GetUserURL(ctx echo.Context) ([]*model.UserURLResponse, error) {
	coupleList := []*model.UserURLResponse{}

	userID, err := s.GetUserID(ctx)
	if err != nil {
		return nil, ctx.JSON(http.StatusUnauthorized, err.Error())
	}

	for _, couple := range s.CacheURL {
		if couple.UserID == userID {
			coupleList = append(coupleList, &model.UserURLResponse{
				OriginalURL: couple.OriginalURL,
				ShortURL:    couple.Alias,
			})
		}
	}
	if len(coupleList) == 0 {
		return nil, ctx.JSON(http.StatusNotFound, "Пользователь не имеет созданных коротких URL")
	}
	return coupleList, nil
}
