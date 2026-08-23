package endpoint

import (
	"Robot-Project/internal/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SendMessage(c *gin.Context) {
	targetId := c.Param("id")

	if !h.store.Exists(targetId) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target bot not found"})
		return
	}

	var req model.SendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.store.SetLatest(targetId, model.Message{
		From:    req.From,
		Payload: req.Payload,
		SentAt:  time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "command set"})
}

func (h *Handler) GetMessage(c *gin.Context) {
	id := c.Param("id")

	if !h.store.Exists(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bot not found"})
		return
	}

	msg, ok := h.store.GetLatest(id)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"message": model.Message{
			Payload: "stop",
			SentAt:  time.Now(),
		}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": msg})
}
