package handler

import (
	"net/http"

	"chat-service/internal/repository"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	repo *repository.ChatRepository
}

func NewChatHandler(repo *repository.ChatRepository) *ChatHandler {
	return &ChatHandler{repo: repo}
}

func (h *ChatHandler) GetHistory(c *gin.Context) {
	messages, err := h.repo.GetHistory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения истории"})
		return
	}
	c.JSON(http.StatusOK, messages)
}
