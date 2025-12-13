package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tartushkin/TSHORT.git/internal/model"
)

func (s *Short) ReaderBody(ctx echo.Context, req string) ([]*model.BranchResponse, error) {
	list, err := s.getBody(ctx, req) // получаем тело запроса
	if err != nil {
		return nil, err
	}
	responseList := []*model.BranchResponse{}
	newList := []*model.AliasFullCore{}
	for _, couple := range list {
		alias, ok := s.checkURL(couple.OriginalURL) //сначала проверяем кеш, потом проверяем наличие в базе
		if ok {
			responseList = append(responseList, &model.BranchResponse{
				ShortURL: alias,
				CorrID:   couple.CorrID,
			})
			//return nil, err
			continue
		}
		newList = append(newList, &model.AliasFullCore{
			Alias:       s.getUUID(),
			OriginalURL: couple.OriginalURL,
			CorrID:      couple.CorrID,
			UserID:      couple.UserID,
		})
	}
	if len(newList) > 0 {
		err = s.insertURL(newList, req)
		if err != nil {
			return nil, err
		}
	}

	for _, couple := range newList {
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
			CorrID:      "undefined",
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
			CorrID:      "undefined",
			UserID:      userID,
		})
	}
	return coupleList, nil
}

func (s *Short) GetUserURL(ctx echo.Context) ([]*model.UserURLResponse, error) {
	coupleList := []*model.UserURLResponse{}

	userID, err := s.GetUserID(ctx)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, err.Error())
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
		return nil, echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("не нашли у пользователя URL"))
	}
	return coupleList, nil
}

func (s *Short) DeleteUserURL(ctx echo.Context, deleteList []string) {
	userID, err := s.GetUserID(ctx)
	if err != nil {
		s.Logger.Error("ошибка - не удалось найти пользователя - ", err.Error())
		//return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	list := []string{}

	for _, couple := range s.CacheURL {
		for _, del := range deleteList {
			if couple.Alias == del { // если нашли в кеше
				if couple.UserID == userID { // если пользователь совпадает
					if couple.DeletedFlag { // если тру идем на некст итерацию
						continue
					}
					couple.DeletedFlag = true
					list = append(list, couple.Alias)
				}

			}
		}
	}

	delStr := "'" + strings.Join(list, "','") + "'"
	err = s.Repo.DeleteURL(s.Ctx, delStr)
	if err != nil {
		s.Logger.Error("ошибка - не удалось удалить URL - ", err.Error())
	}
}
