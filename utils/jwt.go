package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(email string, id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"id":    id,
		"exp":   time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("secretkey")))
}

func VerifyToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected Signing Method")
		}
		return []byte(os.Getenv("secretkey")), nil
	})
	if err != nil {
		return 0, errors.New("could not parse Token")
	}

	isValid := parsedToken.Valid

	if !isValid {
		return 0, errors.New("not a valid token")
	}

	claims, _ := parsedToken.Claims.(jwt.MapClaims)

	id := int64(claims["id"].(float64))

	return id, nil
}
