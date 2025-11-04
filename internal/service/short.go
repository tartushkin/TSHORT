package service

import (
	"context"
	"encoding/json"
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
	return sh, nil
}

func (s *Short) NewFile() (*model.FileStorage, error) {
	file, err := os.OpenFile(s.PathStorage, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &model.FileStorage{
		File:    file,
		Encoder: json.NewEncoder(file),
	}, nil
}
