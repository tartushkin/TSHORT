package model

import (
	"encoding/json"
	"os"
)

const (
	DATABASE    = "DB"
	FILE        = "FILE"
	List        = "list"
	One         = "one"
	Text        = "text"
	CONFLICT    = "данный URL"
	ERRCONFLICT = "ERRCONFLICT"
	JSON        = "json"
)

type PostURLHandlerRequest struct {
	URL string
}
type PostURLHandlerResponse struct {
	ErrMsg string
	Result string
}

type FileStorage struct {
	SURL    *os.File
	Encoder *json.Encoder
}

type BranchRequest struct {
	CorrID      string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}
type BranchResponse struct {
	CorrID   string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type AliasFullCore struct {
	Alias       string `json:"alias"`
	OriginalURL string `json:"original_url"`
	CorrID      string `json:"correlation_id"`
}
