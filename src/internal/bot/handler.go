package bot

import (
	"net/http"
	"time"

	"Robot-Project/internal/util"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) CreateId(c *gin.Context) {
	id, err := util.GenerateCode(8)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate code"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.store.Exists(req.Id) {
		c.JSON(http.StatusConflict, gin.H{"error": "Bot with this ID already exists"})
		return
	}

	h.store.Add(Bot{Id: req.Id})
	c.JSON(http.StatusCreated, gin.H{"message": "Bot registered successfully"})
}

func (h *Handler) SendMessage(c *gin.Context) {
	targetId := c.Param("id")

	if !h.store.Exists(targetId) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target bot not found"})
		return
	}

	var req SendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.store.SetLatest(targetId, Message{
		From:    req.From,
		Payload: req.Payload,
		SentAt:  time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "command set"})
}

func (h *Handler) GetLatestMessage(c *gin.Context) {
	id := c.Param("id")

	if !h.store.Exists(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bot not found"})
		return
	}

	msg, ok := h.store.GetLatest(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": msg})
}
