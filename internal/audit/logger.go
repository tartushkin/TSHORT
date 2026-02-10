// Package audit предоставляет логгеры для аудита: файловый и удалённый.
package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
func (d *Dispatcher) NewFileLogger(path string) error { //(*FileLogger,
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	fl := &FileLogger{file: file}
	d.AddLogger(fl)

	return nil
}

// NewRemoteLogger - создаёт логгер для внешнего аудита.
func (d *Dispatcher) NewRemoteLogger(url string) { //*RemoteLogger
	rl := &RemoteLogger{
		url:    url,
		client: &http.Client{Timeout: 5 * time.Second},
	}
	d.AddLogger(rl)
	//return rl
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

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusInternalServerError:
		return fmt.Errorf("log.err - при отправке аудита влзникла ошибка на стороне сервера.")
	default:
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("log.err - неожиданный статус ответа: %d, body: %s", resp.StatusCode, string(body))
	}
}
