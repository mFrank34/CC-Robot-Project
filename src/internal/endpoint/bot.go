package endpoint

import "Robot-Project/internal/model"

type Handler struct {
	store *model.Store
}

func NewHandler(store *model.Store) *Handler {
	return &Handler{store: store}
}
