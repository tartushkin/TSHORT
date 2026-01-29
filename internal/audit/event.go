package audit

import "time"

type Event struct {
	TS     int64  `json:"ts"`      // Unix timestamp
	Action string `json:"action"`  // shorten/follow
	UserID string `json:"user_id"` // может быть пустым
	URL    string `json:"url"`     // оригинальный URL
}

func NewEvent(action, userID, url string) Event {
	return Event{
		TS:     time.Now().Unix(),
		Action: action,
		UserID: userID,
		URL:    url,
	}
}
