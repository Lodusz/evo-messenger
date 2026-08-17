package handler

import (
	"log"
	"net/http"

	"auth-service/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат"})
		return
	}

	if err := h.service.RegisterUser(req.Username, req.Password); err != nil {
		log.Printf("ОШИБКА ПРИ РЕГИСТРАЦИИ: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при регистрации пользователя. Возможно, имя уже занято."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Успешно"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат"})
		return
	}

	token, err := h.service.LoginUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
