package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chat-service/internal/handler"
	"chat-service/internal/kafka"
	"chat-service/internal/repository"
	"chat-service/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Запуск Срфе Ыукмшсу")

	repo := repository.NewChatRepository()
	kafkaSvc := kafka.NewKafkaService()
	hub := websocket.NewHub(kafkaSvc, repo)
	chatHandler := handler.NewChatHandler(repo)

	go hub.Run()
	go kafkaSvc.ConsumeMessages(hub.Broadcast)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(corsMiddleware())

	r.GET("/history", chatHandler.GetHistory)
	r.GET("/ws", func(c *gin.Context) {
		websocket.ServeWS(hub, c)
	})

	srv := &http.Server{
		Addr:    ":8082",
		Handler: r,
	}

	go func() {
		log.Println("Chat Service слушает порт 8082")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Ошибка при работе HTTP сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Получен сигнал завершения. Остановка Chat Service...")

	// Даем 5 секунд на завершение текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Принудительное завершение сервера: ", err)
	}

	log.Println("Chat Service ыещзз")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
