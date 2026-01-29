package audit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type Logger interface {
	Log(event Event) error
}

type FileLogger struct {
	file *os.File
}

// внешний аудит
type RemoteLogger struct {
	url    string
	client *http.Client
}

func NewFileLogger(path string) (*FileLogger, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &FileLogger{file: file}, nil
}

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

// отправляем во внешний
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
