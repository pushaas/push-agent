package models

type (
	Message struct {
		Channels []string `json:"channels"`
		Content string `json:"content"`
	}
)
