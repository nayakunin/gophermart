package token

import (
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nayakunin/gophermart/internal/config"
)

type Service struct{}

func (s *Service) CreateToken(cfg config.Config, userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		string(cfg.AuthKey): strconv.FormatInt(userID, 10),
	})

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewService() *Service {
	return &Service{}
}
