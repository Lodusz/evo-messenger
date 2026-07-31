package websocket

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"chat-service/internal/kafka"
	"chat-service/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var jwtSecret = []byte("my_secret_key_12345")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients    map[*Client]bool
	Broadcast  chan repository.Message
	register   chan *Client
	unregister chan *Client
	kafkaSvc   *kafka.Service
	repo       *repository.ChatRepository
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	username string
}

func NewHub(k *kafka.Service, r *repository.ChatRepository) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		Broadcast:  make(chan repository.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		kafkaSvc:   k,
		repo:       r,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case msg := <-h.Broadcast:

			bytes, _ := json.Marshal(msg)
			for client := range h.clients {
				select {
				case client.send <- bytes:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		msg := repository.Message{
			Username: c.username,
			Content:  string(message),
		}

		c.hub.repo.SaveMessage(msg.Username, msg.Content)

		c.hub.kafkaSvc.ProduceMessage(msg)

		c.hub.Broadcast <- msg
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func parseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}
	return claims["username"].(string), nil
}

func ServeWS(hub *Hub, c *gin.Context) {

	tokenString := c.Query("token")
	username, err := parseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Недействительный токен"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Ошибка апгрейда соединения:", err)
		return
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		username: username,
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
