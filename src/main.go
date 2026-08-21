package main

import (
	"Robot-Project/internal/bot"
	"Robot-Project/internal/endpoint"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// dependencies
	botStore := bot.NewStore()
	botHandler := endpoint.NewHandler(botStore)

	// routes
	router.GET("/health", endpoint.Check)

	// bot routes for creating IDs and registering bots
	router.GET("/id", botHandler.CreateId)
	router.POST("/id/:id", botHandler.Register)

	// message routes for sending and retrieving messages
	router.POST("/id/:id/message", botHandler.SendMessage)
	router.GET("/id/:id/message", botHandler.GetMessage)

	// creating status routes for setting and getting bot status
	router.POST("/id/:id/status", botHandler.SetStatus)
	router.GET("/id/:id/status", botHandler.GetStatus)

	router.Run(":8080")
}
