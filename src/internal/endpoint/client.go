package endpoint

import (
	"github.com/gin-gonic/gin"
)

func ServeClient(c *gin.Context) {
	c.File("client/client.lua")
}
