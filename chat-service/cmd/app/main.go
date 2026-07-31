package main

import (
	"log"

	"chat-service/internal/handler"
	"chat-service/internal/kafka"
	"chat-service/internal/repository"
	"chat-service/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Запуск Chat Service")

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

	log.Println("Chat Service порт 8082")
	if err := r.Run(":8082"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
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
