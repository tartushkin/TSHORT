// Package model содержит DTO, запросы, ответы и структуры данных,
// используемые между слоями приложения (handler, service, storage).

package model

import (
	"encoding/json"
	"os"
)

// Константы типов хранилищ и форматов.

const (
	DATABASE    = "DB"
	FILE        = "FILE"
	List        = "list"
	CACHE       = "cache"
	One         = "one"
	Text        = "text"
	CONFLICT    = "данный URL"
	ERRCONFLICT = "ERRCONFLICT"
	JSON        = "json"
)

// PostURLHandlerRequest — входной запрос для сокращения URL.
type PostURLHandlerRequest struct {
	URL string `json:"url"`
}

// generate:reset
// PostURLHandlerResponse — ответ хендлера (универсальный).
type PostURLHandlerResponse struct {
	ErrMsg string
	Result string
}

// FileStorage — обёртка для файлового хранилища URL.
// generate:reset
type FileStorage struct {
	SURL    *os.File
	Encoder *json.Encoder
}

// BranchRequest — запрос для сокращения .
type BranchRequest struct {
	CorrID      string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

// BranchResponse — ответ для  сокращения.
type BranchResponse struct {
	CorrID   string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

// AliasFullCore — полная модель записи в хранилище.
type AliasFullCore struct {
	Alias       string `json:"alias"`
	OriginalURL string `json:"original_url"`
	CorrID      string `json:"correlation_id"`
	UserID      string `json:"user_id"`
	DeletedFlag bool   `json:"is_deleted"`
}

// UserURLResponse — ответ API для списка URL пользователя.
type UserURLResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Stats struct {
	Urls  int `json:"urls"`
	Users int `json:"users"`
}
