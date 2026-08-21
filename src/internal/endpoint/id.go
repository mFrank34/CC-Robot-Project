package endpoint

import (
	"Robot-Project/internal/bot"
	"Robot-Project/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateId(c *gin.Context) {
	id, err := util.GenerateCode(8)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate code"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) Register(c *gin.Context) {
	var req bot.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.store.Exists(req.Id) {
		c.JSON(http.StatusConflict, gin.H{"error": "Bot with this ID already exists"})
		return
	}

	h.store.Add(bot.Bot{Id: req.Id})
	c.JSON(http.StatusCreated, gin.H{"message": "Bot registered successfully"})
}
