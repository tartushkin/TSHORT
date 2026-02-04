package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/tartushkin/TSHORT.git/internal/audit"
	cfg "github.com/tartushkin/TSHORT.git/internal/config/app"
	"github.com/tartushkin/TSHORT.git/internal/config/db"
	"github.com/tartushkin/TSHORT.git/internal/model"
	"github.com/tartushkin/TSHORT.git/internal/repository"
)

const defaultParamDelete = 20 // дефолтный параметр на запуска процесса уадаления
// Short — основной сервис для работы с URL.
type Short struct {
	mu sync.RWMutex

	Logger *logrus.Logger
	Ctx    context.Context

	CacheURL map[string]*model.AliasFullCore
	Conn     *sql.DB
	Repo     *repository.Repo
	//client   *Client

	HTTPPort    string
	Address     string
	PathStorage string

	DNS         string
	paramDelete time.Duration

	auditLocal string
	auditURL   string

	FileStorage *model.FileStorage
	Dis         *audit.Dispatcher
	Fcpu        *os.File
	Fmem        *os.File

	RunProfile bool
}

// NewShort - заполнение структуры приложения.
func Create(ctx context.Context, lg *logrus.Logger, cfg *cfg.Config) (*Short, error) {
	cacheURL := map[string]*model.AliasFullCore{}

	sh := &Short{
		Logger:   lg,
		Ctx:      ctx,
		CacheURL: cacheURL,
		HTTPPort: cfg.Port,
		Address:  cfg.Address,
		Dis:      audit.NewDispatcher(),
	}

	if cfg.FileStoragePath != "" {
		sh.PathStorage = cfg.FileStoragePath
		file, err := sh.NewFile(sh.PathStorage)
		if err != nil {
			return nil, err
		}
		sh.FileStorage = file
	}

	if cfg.DNS != "" {
		conn, err := db.NewConnection(cfg.DNS)
		if err != nil {
			lg.Info("db: не удалось подключилиться к DB, используем другое хранилище")
		} else {
			sh.DNS = cfg.DNS
			lg.Info("db: успешно подключились к DB")
			sh.Conn = conn
			sh.Repo = repository.NewRepository(sh.Conn)
		}

	}

	sh.paramDelete = defaultParamDelete * time.Second
	if cfg.ParamDelete != 0 {
		sh.paramDelete = time.Duration(cfg.ParamDelete) * time.Second
	}

	if cfg.LocalAuditPath != "" {
		sh.auditLocal = cfg.LocalAuditPath
		file, err := audit.NewFileLogger(sh.auditLocal)
		if err != nil {
			lg.Error("create.audit - ошибка создания файла для аудита", err)
		}
		sh.Dis.AddLogger(file)
	}
	if cfg.AuditPath != "" {
		sh.auditURL = cfg.AuditPath
		au := audit.NewRemoteLogger(sh.auditURL)
		sh.Dis.AddLogger(au)
	}

	if cfg.RunProfile {
		err := sh.CPUProfile()
		if err != nil {
			return nil, err
		}
	}
	go sh.StartCleanup(ctx, sh.paramDelete)
	return sh, nil
}

// NewFile - создание файла для хранения пар URL
func (s *Short) NewFile(path string) (*model.FileStorage, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
	if s.FileStorage != nil {
		s.FileStorage.SURL.Close()
		s.Logger.Info("main: ", fmt.Sprintf("file - %s, успешно закрыт", s.PathStorage))
	}
	if s.Conn != nil {
		s.Conn.Close()
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
		_, err := s.FileStorage.SURL.Seek(0, 0)
		if err != nil {
			s.Logger.Error("Ошибка перемещения указателя файла: ", err)
			return err
		}
		decoder := json.NewDecoder(s.FileStorage.SURL)
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

// StartCleanup - процесс очистки помеченных на удаления URL
func (s *Short) StartCleanup(ctx context.Context, param time.Duration) {
	ticker := time.NewTicker(param)
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

// CPUProfile -  профилирования CPU
func (s *Short) CPUProfile() error {
	fcpu, err := os.OpenFile("./profiles/result.pprof", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	s.Fcpu = fcpu

	if err := pprof.StartCPUProfile(fcpu); err != nil {
		return err
	}
	s.Logger.Info("CPU профилирование запущено")
	return nil
}

// MemProfile -  профилирования памяти
func (s *Short) MemProfile() error {
	// создаём файл журнала профилирования памяти
	fmem, err := os.OpenFile("./profiles/result.pprof", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	s.Fmem = fmem
	runtime.GC() // получаем статистику по использованию памяти
	if err := pprof.WriteHeapProfile(fmem); err != nil {
		return err
	}
	s.Logger.Info("Профили сохранены: cpu.pprof, base.pprof")
	return nil
}
