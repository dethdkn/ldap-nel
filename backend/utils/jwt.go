package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func JWTGenerate(username string, admin bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"admin":    admin,
		"exp":      time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func JWTValidate(tokenString string) (string, bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return "", false, err
	}

	if !token.Valid {
		return "", false, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", false, errors.New("could not parse claims")
	}

	username, ok := claims["username"].(string)
	if !ok || username == "" {
		return "", false, errors.New("username not found in token")
	}

	admin, ok := claims["admin"].(bool)
	if !ok {
		return "", false, errors.New("admin status not found in token")
	}

	return username, admin, nil
}
