package service

import (
	"encoding/json"
	"fmt"

	"github.com/tartushkin/TSHORT.git/internal/model"
)

func (s *Short) ReaderBody(body []byte, userID, req string) ([]*model.BranchResponse, error) {
	list, err := s.checkBody(body, userID, req) // получаем тело запроса
	if err != nil {
		return nil, err
	}
	responseList := make([]*model.BranchResponse, 0, len(list))
	newList := make([]*model.AliasFullCore, 0, len(list))

	for _, couple := range list {
		alias, ok := s.checkURL(couple.OriginalURL) //сначала проверяем кеш, потом проверяем наличие в базе
		if ok {
			responseList = append(responseList, &model.BranchResponse{
				ShortURL: alias,
				CorrID:   couple.CorrID,
			})
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
		fullURL := s.getFull(couple.Alias)
		responseList = append(responseList, &model.BranchResponse{ShortURL: fullURL, CorrID: couple.CorrID})
		s.Logger.Info(fmt.Sprintf("ReaderBody.jsonList - запись пары в список для отправки в хранилище: %s/%s/%s/%s ",
			couple.Alias, couple.OriginalURL, couple.CorrID, couple.UserID))
	}
	return responseList, nil
}

func (s *Short) checkBody(body []byte, userID, req string) ([]*model.AliasFullCore, error) {

	s.Logger.Info("ReaderBody.start - чтение тела запроса")

	coupleList := []*model.AliasFullCore{}
	var couple *model.PostURLHandlerRequest

	switch req {
	case model.One:
		err := json.Unmarshal(body, &couple)
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

		err := json.Unmarshal(body, &coupleList)
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

func (s *Short) GetUserURL(userID string) ([]*model.UserURLResponse, error) {
	coupleList := []*model.UserURLResponse{}
	for _, couple := range s.CacheURL {
		if couple.UserID == userID {
			coupleList = append(coupleList, &model.UserURLResponse{
				OriginalURL: couple.OriginalURL,
				ShortURL:    fmt.Sprintf("%s/%s", s.Address, couple.Alias),
			})
		}
	}
	if len(coupleList) == 0 {
		list, err := s.Repo.GetUserURL(s.Ctx, userID)
		if err != nil {
			return nil, err
		}
		coupleList = make([]*model.UserURLResponse, 0, len(coupleList))
		for _, couple := range list {
			coupleList = append(coupleList, &model.UserURLResponse{
				OriginalURL: couple.OriginalURL,
				ShortURL:    s.Address + "/" + couple.Alias, //fmt.Sprintf("%s/%s", s.Address, couple.Alias),
			})
		}
		if len(list) == 0 {
			return []*model.UserURLResponse{}, nil
		}
		return coupleList, nil
	}
	return coupleList, nil
}

func (s *Short) DeleteUserURL(userID string, deleteList []string) error {
	listDel := make([]string, 0, len(deleteList))

	for _, couple := range s.CacheURL {
		for _, del := range deleteList {
			if couple.Alias == del { // если нашли в кеше
				if couple.UserID == userID { // если пользователь совпадает
					if couple.DeletedFlag { // если тру идем на некст итерацию
						continue
					}
					couple.DeletedFlag = true
					listDel = append(listDel, couple.Alias)
				}

			}
		}
	}

	if len(listDel) == 0 { // в кеше пусто, смотримм в базе
		list, err := s.Repo.GetUserURL(s.Ctx, userID)
		if err != nil {
			return err
		}
		for _, couple := range list {
			for _, del := range deleteList {
				if couple.Alias == del { // если нашли в базе
					if couple.DeletedFlag { // если тру идем на некст итерацию
						continue
					}
					couple.DeletedFlag = true
					listDel = append(listDel, couple.Alias)

				}
			}
		}

	}
	if len(listDel) == 0 {
		return fmt.Errorf("ошибка: не найден ни один url")
	}
	err := s.Repo.DeleteURL(s.Ctx, listDel)
	if err != nil {
		s.Logger.Error("ошибка - не удалось пометить URL на удаление - ", err.Error())
		return err
	}
	return nil
}
