package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	model "github.com/Sunr1s/webclient"
	"github.com/Sunr1s/webclient/pkg/repository"
	"github.com/golang-jwt/jwt"
)

const (
	salt       = "NK5UJ0vBobKi1MOKuo3BoFk9xlABvd6jQNKhdJFv24IuXWi2lLlfJyCL7km0"
	signingKey = "pDxONfwbE4c9vBooWAmE6JO0hL5AQYu2wWJR0h9oh2E1h8PRdZZYVyYA1bFE"
	tokenTTL   = 12 * time.Hour
)

// Claims defines the structure of the JWT claims
type Claims struct {
	jwt.StandardClaims
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

// AuthService handles operations related to authorization
type AuthService struct {
	repo repository.Authorization
}

// NewAuthService initializes a new AuthService
func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

// RegisterUser registers a new user in the system
func (s *AuthService) RegisterUser(user model.User) (int, error) {
	user.Password = hashPassword(user.Password)
	return s.repo.RegisterUser(user)
}

// GenerateToken generates a new JWT token for a user
func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, hashPassword(password))
	if err != nil {
		return "", err
	}

	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID:   user.Id,
		Username: user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(signingKey))
}

// ParseToken parses a JWT token and returns the user details
func (s *AuthService) ParseToken(accessToken string) (int, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, "", errors.New("token claims are not of expected type")
	}

	return claims.UserID, claims.Username, nil
}

// hashPassword hashes a password with a salt
func hashPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
