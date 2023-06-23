package dto

type Notification struct {
	Channel    string `json:"channel"`
	Message    string `json:"message"`
	RetryCount int    `json:"retry_count"`
}
