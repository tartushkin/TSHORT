package audit

import (
	"net/http"
	"time"
)

const (
	get     = "follow"
	post    = "shorten"
	unknown = "unknown"
)

// Event - собитие аудита.
type Event struct {
	TS     int64  `json:"ts"`      // Unix timestamp
	Action string `json:"action"`  // Shorten/followw
	UserID string `json:"user_id"` // Mожет быть пустым
	URL    string `json:"url"`     // Оригинальный URL
}

// NewEvent - формировние нового события.
func NewEvent(method, userID, url string) Event {
	var action string
	switch method {
	case http.MethodGet:
		action = get
	case http.MethodPost:
		action = post
	default:
		action = unknown
	}

	return Event{
		TS:     time.Now().Unix(),
		Action: action,
		UserID: userID,
		URL:    url,
	}
}
