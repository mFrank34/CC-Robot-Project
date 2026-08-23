package model

import "time"

type Bot struct {
	Id string `json:"id"`
}

type RegisterRequest struct {
	Id string `json:"id"`
}

type Ids struct {
	Ids []string `json:"ids"`
}

type Message struct {
	From    string    `json:"from"`
	Payload Command   `json:"payload"`
	Arg     string    `json:"arg,omitempty"`
	SentAt  time.Time `json:"timestamp"`
}

type SendRequest struct {
	From    string `json:"from"`
	Payload string `json:"payload"`
	Args    string `json:"args"`
}

type InventoryItem struct {
	Slot  int    `json:"slot"`
	Item  string `json:"item"`
	Count int    `json:"count"`
}

type Status struct {
	Fuel      int             `json:"fuel"`
	Inventory []InventoryItem `json:"inventory"`
	UpdatedAt time.Time       `json:"updated_at"`
}
