package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Robot-Project/internal/bot"
	"Robot-Project/internal/endpoint"

	"github.com/gin-gonic/gin"
)

const dataFile = "data/store.json"

func main() {
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("failed to create data dir: %v", err)
	}

	botStore, err := bot.Load(dataFile)
	if err != nil {
		log.Fatalf("failed to load store: %v", err)
	}

	router := gin.Default()

	// dependencies
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

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// run the server in a goroutine so main() is free to wait for a signal
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// block until Ctrl+C or a kill signal arrives
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutting down, saving store...")
	if err := botStore.Save(dataFile); err != nil {
		log.Printf("failed to save store on exit: %v", err)
	}

	// give in-flight requests up to 5s to finish before the process exits
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("done")
}
