package model

import (
	"encoding/json"
	"os"
)

type PostURLHandlerRequest struct {
	URL string
}
type PostURLHandlerResponse struct {
	ErrMsg string
	Result string
}

type FileStorage struct {
	File    *os.File
	Encoder *json.Encoder
}

type StorageURL struct {
	Alias    string `json:"short_url"`
	Original string `json:"original_url"`
}
