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

	"auth-service/internal/handler"
	"auth-service/internal/repository"
	"auth-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Запуск Auth Service")

	repo := repository.NewAuthRepository()
	authService := service.NewAuthService(repo)
	authHandler := handler.NewAuthHandler(authService)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(corsMiddleware())

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	go func() {
		log.Println("Auth Service слушает порт 8081")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Ошибка при работе HTTP сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Получен сигнал завершения. Остановка Auth Service")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Принудительное завершение сервера Auth Service: %v\n", err)
	}

	log.Println("Auth Service остановлен!")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
