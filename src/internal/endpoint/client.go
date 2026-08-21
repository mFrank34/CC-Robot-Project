package endpoint

import "github.com/gin-gonic/gin"

func serveClient(c *gin.Context) {
	c.File("src/client/client.lua")
}
