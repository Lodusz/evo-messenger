package service

import (
	"errors"
	"time"

	"auth-service/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      *repository.AuthRepository
	jwtSecret []byte
}

func NewAuthService(repo *repository.AuthRepository, secret string) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: []byte(secret),
	}
}

func (s *AuthService) RegisterUser(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.CreateUser(username, string(hash))
}

func (s *AuthService) LoginUser(username, password string) (string, error) {
	dbPasswordHash, err := s.repo.GetPasswordByUsername(username)
	if err != nil {
		return "", errors.New("неверный логин или пароль")
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("неверный логин или пароль")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString(s.jwtSecret)
}
