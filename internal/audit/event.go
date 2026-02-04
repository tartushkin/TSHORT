package audit

import "time"

// Event - собитие аудита.
type Event struct {
	TS     int64  `json:"ts"`      // Unix timestamp
	Action string `json:"action"`  // Shorten/followw
	UserID string `json:"user_id"` // Mожет быть пустым
	URL    string `json:"url"`     // Оригинальный URL
}

// NewEvent - формировние нового события.
func NewEvent(action, userID, url string) Event {
	return Event{
		TS:     time.Now().Unix(),
		Action: action,
		UserID: userID,
		URL:    url,
	}
}
