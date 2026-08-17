package websocket

import (
	"encoding/json"
	"net/http"
	"os"

	"chat-service/internal/kafka"
	"chat-service/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

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
		clients: make(map[*Client]bool), Broadcast: make(chan repository.Message),
		register: make(chan *Client), unregister: make(chan *Client),
		kafkaSvc: k, repo: r,
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
	defer func() { c.hub.unregister <- c; c.conn.Close() }()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		msg := repository.Message{Username: c.username, Content: string(message)}
		c.hub.repo.SaveMessage(msg.Username, msg.Content)
		c.hub.kafkaSvc.ProduceMessage(msg)
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func ServeWS(hub *Hub, c *gin.Context) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(c.Query("token"), func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	conn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), username: token.Claims.(jwt.MapClaims)["username"].(string)}
	client.hub.register <- client
	go client.writePump()
	go client.readPump()
}
