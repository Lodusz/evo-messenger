package service

import (
	"errors"
	"time"

	"auth-service/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("my_secret_key_12345")

type AuthService struct {
	repo *repository.AuthRepository
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) RegisterUser(username, password string) error {
	return s.repo.CreateUser(username, password)
}

func (s *AuthService) LoginUser(username, password string) (string, error) {
	dbPassword, err := s.repo.GetPasswordByUsername(username)
	if err != nil {
		return "", errors.New("пользователь не найден")
	}

	if dbPassword != password {
		return "", errors.New("неверный пароль")
	}

	// Генерируем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString(jwtSecret)
}
