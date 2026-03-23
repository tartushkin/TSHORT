package service

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/tartushkin/TSHORT.git/internal/model"
	pb "github.com/tartushkin/TSHORT.git/pkg/shortenerservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const undefined = "undefined"

func (sh *Short) getUserId(ctx context.Context) (string, error) {
	// Извлекаем метаданные из контекста
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("метаданные отсутсвуют")
	}

	for key, values := range md {
		sh.Logger.Infof("Ключ: %s, Значения: %v", key, values)
	}
	// Получаем ID пользователя из метаданных
	userIDs := md.Get("userID")
	if len(userIDs) == 0 {
		return "", errors.New("user ID отсутствуют")
	}

	userID := userIDs[0]
	sh.Logger.Info("User ID: ", userID)
	return userID, nil
}
func (sh *Short) StartGRPC(port string) error {
	// Создаем TCP-слушатель на указанном порту
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf(" failed to listen: %v", err)
	}

	// Создаем новый gRPC-сервер
	s := grpc.NewServer()

	// Регистрируем сервис Greeter на сервере
	pb.RegisterShortenerServiceServer(s, sh)

	// Запускаем сервер
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	sh.Logger.Info("StartGRPC - запущен GRPC сервер")

	return nil
}

func (sh *Short) ShortenURL(ctx context.Context, r *pb.URLShortenRequest) (*pb.URLShortenResponse, error) {
	sh.Logger.Info("ShortenURL.start - регистрации нового URL")

	newList := make([]*model.AliasFullCore, 0, 1)
	id, err := sh.getUserId(ctx)
	if err != nil {
		return nil, err
	}
	couple := &model.AliasFullCore{
		Alias:       sh.getUUID(),
		OriginalURL: r.Url,
		CorrID:      undefined,
		UserID:      id,
	}
	newList = append(newList, couple)
	err = sh.insertURL(newList, model.One)
	if err != nil {
		return nil, err
	}
	res := pb.URLShortenResponse{Result: sh.getFull(couple.Alias)}
	sh.Logger.Info("ShortenURL.finish - успешно зарегестрировали новый URL")
	return &res, nil
}

func (sh *Short) ExpandURL(ctx context.Context, r *pb.URLExpandRequest) (*pb.URLExpandResponse, error) {
	sh.Logger.Info("ExpandURL.start - запрос оригинального URL по алиасу")
	res := &pb.URLExpandResponse{}
	originalURL, err := sh.GetAliasName(r.Id)
	if err != nil {
		sh.Logger.Error("getRedirectHandler.err - возникла ошбка при получениии оригинального URL: " + err.Error())
		return nil, err
	}
	res.Result = originalURL
	sh.Logger.Info("ExpandURL.finish - успешно получили оригинальный URL")
	return res, nil
}
func (sh *Short) ListUserURLs(ctx context.Context, r *pb.Empty) (*pb.UserURLsResponse, error) {
	sh.Logger.Info("ListUserURLs.start - получние url-ов пользователя")
	res := &pb.UserURLsResponse{}
	id, err := sh.getUserId(ctx)
	if err != nil {
		return nil, err
	}
	userIDList, err := sh.GetUserURL(id)
	if err != nil {
		sh.Logger.Error("Ошибка при работе с телом запроса: " + err.Error())
		return nil, err
	}
	if len(userIDList) == 0 {
		return nil, errors.New("у пользователя отсутсвуют URL")
	}
	for _, couple := range userIDList {
		res.Url = append(res.Url, &pb.URLData{
			ShortUrl:    couple.ShortURL,
			OriginalUrl: couple.OriginalURL,
		})
	}

	sh.Logger.Info("ListUserURLs.finish - возвращаем список URL: ", id)
	return res, nil
}
