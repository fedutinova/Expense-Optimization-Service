package service

import (
	"crypto/sha1"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"os"
	"time"
	wallet "wallet-app/pkg/models"
	errs "wallet-app/pkg/models/errors"
	"wallet-app/pkg/repository"
)

const (
	tokenTTL = 12 * time.Hour
)

var (
	signingKey = os.Getenv("signingKey")
	salt       = os.Getenv("salt")
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo   repository.Authorization
	logger *logrus.Logger
}

func NewAuthService(repo repository.Authorization, logger *logrus.Logger) *AuthService {
	return &AuthService{
		repo:   repo,
		logger: logger,
	}
}

func (s *AuthService) CreateUser(user wallet.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)

	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}
func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.logger.Errorf("failed signing")
			return nil, errs.InternalServer
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		s.logger.Errorf("failed parse token %v", err)
		return 0, errs.InternalServer
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		s.logger.Error("token claims are not of type *tokenClaims")
		return 0, errs.InternalServer
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
