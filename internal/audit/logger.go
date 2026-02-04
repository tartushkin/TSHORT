// Package audit предоставляет логгеры для аудита: файловый и удалённый.
package audit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

// Logger — интерфейс для записи событий аудита.
type Logger interface {
	Log(event Event) error
}

// FileLogger - пишет события в файл.
type FileLogger struct {
	file *os.File
}

// RemoteLogger -  внешний аудит.
type RemoteLogger struct {
	url    string
	client *http.Client
}

// NewFileLogger создаёт новый файловый логгер.
func NewFileLogger(path string) (*FileLogger, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &FileLogger{file: file}, nil
}

// NewRemoteLogger - создаёт логгер для внешнего аудита.
func NewRemoteLogger(url string) *RemoteLogger {
	return &RemoteLogger{
		url:    url,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

// Log - запись в файл
func (l *FileLogger) Log(event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, err = l.file.Write(append(data, '\n'))
	return err
}

func (l *FileLogger) Close() error {
	return l.file.Close()
}

// Log - отправляем во внешний ресивер.
func (l *RemoteLogger) Log(event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	resp, err := l.client.Post(l.url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
