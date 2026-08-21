package endpoint

import "Robot-Project/internal/bot"

type Handler struct {
	store *bot.Store
}

func NewHandler(store *bot.Store) *Handler {
	return &Handler{store: store}
}
