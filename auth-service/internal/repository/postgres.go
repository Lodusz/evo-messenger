package repository

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository() *AuthRepository {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/evo_messenger?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Ошибка драйвера БД: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("БД нsедоступна: %v", err)
	}

	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(username, password string) error {
	_, err := r.db.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, password)
	return err
}

func (r *AuthRepository) GetPasswordByUsername(username string) (string, error) {
	var hash string
	err := r.db.QueryRow("SELECT password_hash FROM users WHERE username = $1", username).Scan(&hash)
	return hash, err
}
