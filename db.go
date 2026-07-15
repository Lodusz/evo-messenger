package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

type MessageResponse struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func (db *DB) GetOrCreateUser(username string) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&id)
	if err == sql.ErrNoRows {
		err = db.QueryRow("INSERT INTO users (username) VALUES ($1) RETURNING id", username).Scan(&id)
	}
	return id, err
}

func (db *DB) SaveMessage(userID int, content string) error {
	_, err := db.Exec("INSERT INTO messages (user_id, content) VALUES ($1, $2)", userID, content)
	return err
}

func (db *DB) GetHistory() ([]MessageResponse, error) {
	rows, err := db.Query("SELECT u.username, m.content FROM messages m JOIN users u ON m.user_id = u.id ORDER BY m.id ASC LIMIT 50")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []MessageResponse
	for rows.Next() {
		var msg MessageResponse
		if err := rows.Scan(&msg.Username, &msg.Content); err == nil {
			msgs = append(msgs, msg)
		} else {
			log.Println("Ошибка парсинга ", err)
		}
	}
	return msgs, nil
}
