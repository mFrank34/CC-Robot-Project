package main

import (
	"Robot-Project/internal/bot"
	"Robot-Project/internal/health"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// dependencies
	botStore := bot.NewStore()
	botHandler := bot.NewHandler(botStore)

	// routes
	router.GET("/health", health.Check)

	// bot routes for creating IDs and registering bots
	router.GET("/id", botHandler.CreateId)
	router.POST("/id/:id", botHandler.Register)

	// message routes for sending and retrieving messages
	router.POST("/id/:id/message", botHandler.SendMessage)
	router.GET("/id/:id/message", botHandler.GetLatestMessage)

	router.Run(":8080")
}
