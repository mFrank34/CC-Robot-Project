package endpoint

import (
	"Robot-Project/internal/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SetStatus(c *gin.Context) {
	id := c.Param("id")

	if !h.store.Exists(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bot not found"})
		return
	}

	var st model.Status
	if err := c.ShouldBindJSON(&st); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	st.UpdatedAt = time.Now()

	h.store.SetStatus(id, st)
	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}

func (h *Handler) GetStatus(c *gin.Context) {
	id := c.Param("id")

	if !h.store.Exists(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bot not found"})
		return
	}

	st, ok := h.store.GetStatus(id)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"status": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": st})
}
