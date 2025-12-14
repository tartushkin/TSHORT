package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	cfg "github.com/tartushkin/TSHORT.git/internal/config/app"
	"github.com/tartushkin/TSHORT.git/internal/config/db"
	"github.com/tartushkin/TSHORT.git/internal/model"
	"github.com/tartushkin/TSHORT.git/internal/repository"
)

const defaultParamDelete = 2 // дефолтный параметр на запуска процесса уадаления

type Short struct {
	Logger      *logrus.Logger
	Ctx         context.Context
	HTTPPort    string
	Address     string
	PathStorage string
	File        *model.FileStorage
	DNS         string
	CacheURL    map[string]*model.AliasFullCore
	conn        *sql.DB
	Repo        *repository.Repo
	mu          sync.RWMutex

	paramDelete time.Duration
}

// NewShort - заполнение структуры приложения
func Create(ctx context.Context, lg *logrus.Logger, cfg *cfg.Config) (*Short, error) {
	cacheURL := map[string]*model.AliasFullCore{}

	sh := &Short{
		Logger:   lg,
		Ctx:      ctx,
		CacheURL: cacheURL,
		HTTPPort: cfg.Port,
	}
	if cfg.FileStoragePath != "" {
		sh.PathStorage = cfg.FileStoragePath
		file, err := sh.NewFile()
		if err != nil {
			return nil, err
		}
		sh.File = file
	}
	if cfg.DNS != "" {
		conn, err := db.NewConnection(cfg.DNS)
		if err != nil {
			lg.Info("db: не удалось подключилиться к DB, используем другое хранилище")
		} else {
			sh.DNS = cfg.DNS
			lg.Info("db: успешно подключились к DB")
			sh.conn = conn
			sh.Repo = repository.NewRepository(sh.conn)
		}

	}
	sh.paramDelete = defaultParamDelete * time.Minute
	if cfg.ParamDelete != 0 {
		sh.paramDelete = time.Duration(cfg.ParamDelete) * time.Minute
	}
	sh.Address = cfg.Address
	go sh.StartCleanup(ctx, sh.paramDelete)
	//err := sh.LoadStorageURL() //подгрузка кеша
	//if err != nil {
	//	return nil, err
	//}
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
	sourse := s.checkSourse()

	switch sourse {
	case model.DATABASE:
		list, err := s.Repo.LoadCache(s.Ctx)
		if err != nil {
			return err
		}
		for _, line := range list {
			s.CacheURL[line.Alias] = line
			s.Logger.Info(fmt.Sprintf("Прочитано и подгружено из БД в кеш пара: key:%v, value:%v", line.Alias, line.OriginalURL))
		}
	case model.FILE:
		_, err := s.File.SURL.Seek(0, 0)
		if err != nil {
			s.Logger.Error("Ошибка перемещения указателя файла: ", err)
			return err
		}
		decoder := json.NewDecoder(s.File.SURL)
		var line model.AliasFullCore
		for decoder.More() {
			err := decoder.Decode(&line)
			if err != nil {
				s.Logger.Error("Ошибка декодирования JSON: ", err)
				return err
			}
			s.CacheURL[line.Alias] = &line
			s.Logger.Info(fmt.Sprintf("Прочитано и подгружено в кеш пара из файла: aliasKey:%v, originalUrl:%v, corrID:%v", line.Alias, line.OriginalURL, line.CorrID))
		}
	}

	return nil

}

func (s *Short) StartCleanup(ctx context.Context, param time.Duration) {
	ticker := time.NewTicker(time.Second)
	s.Logger.Info("StartCleanup.start - старт процесса очистки помеченных на удаления URL")
	for {
		s.Logger.Info("StartCleanup.wait - ожидание новой итерации очистки")
		select {
		case <-ctx.Done():
			s.Logger.Info("StartCleanup.cancel - контекст процесса был завршен")
			ticker.Stop()
			return
		case <-ticker.C:
			ticker.Reset(param)
			err := s.DeleteMarkedURLs(ctx)
			if err != nil {
				s.Logger.Error("Ошибка при удалении помеченных URL: " + err.Error())
				continue
			}
			s.Logger.Info("StartCleanup.complete - успешная очистка помеченных на удаления URL")
		}
	}
}
