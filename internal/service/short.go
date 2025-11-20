package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	cfg "github.com/tartushkin/TSHORT.git/internal/config/app"
	db "github.com/tartushkin/TSHORT.git/internal/config/db"
	"github.com/tartushkin/TSHORT.git/internal/model"
	"github.com/tartushkin/TSHORT.git/internal/repository"
)

type Short struct {
	Logger      *logrus.Logger
	Ctx         context.Context
	HTTPPort    string
	Address     string
	PathStorage string
	File        *model.FileStorage
	DNS         string
	CacheURL    map[string]string
	conn        *sql.DB
	Repo        *repository.Repo
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
		DNS:         cfg.DNS,
	}

	sh.Address = cfg.Address + "/"
	file, err := sh.NewFile()
	if err != nil {
		return nil, err
	}
	sh.File = file

	conn, err := db.NewConnection(cfg.DNS)
	if err != nil {
		panic(err)
	}
	lg.Info("db: успешно подключились к DB")
	sh.conn = conn
	sh.Repo = repository.NewRepository(sh.conn)

	sh.LoadStorageURL() //подгрузка кеша
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
	if s.conn != nil {
		s.conn.Close()
		s.Logger.Info("main: соединение с БД закрыто")
	}
	s.Logger.Info("main: file - для закрытия отсутствует")
}

// LoadStorageURL - подгрузка в кеш
func (s *Short) LoadStorageURL() error {
	_, err := s.File.SURL.Seek(0, 0)
	if err != nil {
		s.Logger.Error("Ошибка перемещения указателя файла: ", err)
		return err
	}
	sourse := s.checkSourse()

	switch sourse {
	case model.DATABASE:
		list, err := s.Repo.GetURLList(s.Ctx)
		if err != nil {
			return err
		}
		for _, line := range list {
			s.CacheURL[line.Alias] = line.Original
			s.Logger.Info(fmt.Sprintf("Прочитано и подгружено из БД в кеш пара: key:%v, value:%v", line.Alias, line.Original))
		}
	case model.FILE:
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
	}

	return nil

}
