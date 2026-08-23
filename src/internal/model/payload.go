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
	Args    string    `json:"arg,omitempty"`
	SentAt  time.Time `json:"timestamp"`
}

type InventoryItem struct {
	Slot  int    `json:"slot"`
	Item  string `json:"item"`
	Count int    `json:"count"`
}

type Status struct {
	Fuel        int             `json:"fuel"`
	Inventory   []InventoryItem `json:"inventory"`
	State       string          `json:"state"`
	LastCommand Command         `json:"last_command"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
