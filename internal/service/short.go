package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	cfg "github.com/tartushkin/TSHORT.git/internal/config/app"
	"github.com/tartushkin/TSHORT.git/internal/model"
)

type Short struct {
	Logger      *logrus.Logger
	Ctx         context.Context
	HTTPPort    string
	Address     string
	PathStorage string
	File        *model.FileStorage
	CacheURL    map[string]string
	mu          sync.RWMutex
}

// NewShort - заполнение структуры приложения
func Create(ctx context.Context, lg *logrus.Logger, cfg *cfg.Config) (*Short, error) {
	cacheURL := map[string]string{}

	sh := &Short{
		Logger:      lg,
		Ctx:         ctx,
		CacheURL:    cacheURL,
		HTTPPort:    cfg.Port,
		PathStorage: cfg.StorageURL,
	}

	sh.Address = cfg.Address + "/"
	file, err := sh.NewFile()
	if err != nil {
		return nil, err
	}
	sh.File = file
	sh.LoadStorageURL()
	return sh, nil
}

// NewFile - создание файла для хранения пар URL
func (s *Short) NewFile() (*model.FileStorage, error) {
	file, err := os.OpenFile(s.PathStorage, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &model.FileStorage{
		SURL:    file,
		Encoder: json.NewEncoder(file),
	}, nil
}

// Close - закрытие файла
func (s *Short) Close() {
	if s.File != nil {
		s.File.SURL.Close()
		s.Logger.Info("main: ", fmt.Sprintf("file - %s, успешно закрыт", s.PathStorage))
	}
	s.Logger.Info("main: file - для закрытия отсутствует")
}

// LoadStorageURL - подгрузка в кеш из файла
func (s *Short) LoadStorageURL() error {
	_, err := s.File.SURL.Seek(0, 0)
	if err != nil {
		s.Logger.Error("Ошибка перемещения указателя файла: ", err)
		return err
	}

	decoder := json.NewDecoder(s.File.SURL)
	var line model.StorageURL
	for decoder.More() {
		err := decoder.Decode(&line)
		if err != nil {
			s.Logger.Error("Ошибка декодирования JSON: ", err)
			return err
		}
		s.CacheURL[line.Alias] = line.Original
		s.Logger.Info(fmt.Sprintf("Прочитано и подгружено в кеш пара из файла: key:%v, value:%v", line.Alias, line.Original))
	}

	return nil

}
