package auth

import (
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(userID int64, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": strconv.FormatInt(userID, 10),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
