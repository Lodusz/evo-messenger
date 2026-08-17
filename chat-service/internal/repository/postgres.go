package repository

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Message struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository() *ChatRepository {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("db не задан в окружении")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Ошибка драйвера БД: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("БД недоступна: %v", err)
	}

	return &ChatRepository{db: db}
}

func (r *ChatRepository) SaveMessage(username, content string) error {
	_, err := r.db.Exec("INSERT INTO messages (username, content) VALUES ($1, $2)", username, content)
	return err
}

func (r *ChatRepository) GetHistory() ([]Message, error) {
	rows, err := r.db.Query("SELECT id, username, content, created_at FROM messages ORDER BY created_at ASC LIMIT 100")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.Username, &msg.Content, &msg.CreatedAt); err != nil {
			log.Printf("Ошибка парсинга строки: %v", err)
			continue
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
