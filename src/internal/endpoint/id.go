package endpoint

import (
	"Robot-Project/internal/model"
	"Robot-Project/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateId(c *gin.Context) {
	id, err := util.GenerateCode(4)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate code"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) Register(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bot ID is required"})
		return
	}

	if h.store.Exists(id) {
		c.JSON(http.StatusConflict, gin.H{"error": "Bot with this ID already exists"})
		return
	}

	h.store.Add(model.Bot{Id: id})
	c.JSON(http.StatusCreated, gin.H{"message": "Bot registered successfully"})
}

func (h *Handler) AllIDs(c *gin.Context) {
	rawids := h.store.GetAllIDs()

	c.JSON(http.StatusOK, gin.H{
		"ids": rawids,
	})
}
