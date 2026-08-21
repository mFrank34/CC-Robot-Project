package bot

import "time"

type Bot struct {
	Id string `json:"id"`
}

type RegisterRequest struct {
	Id string `json:"id"`
}

type Message struct {
	From    string    `json:"from"`
	Payload string    `json:"payload"`
	SentAt  time.Time `json:"timestamp"`
}

type SendRequest struct {
	From    string `json:"from"`
	Payload string `json:"payload"`
}
