package jwttoken

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const TokenLifeTime = time.Hour * 3

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

func validateToken(tokenString string, jwtSecret string) bool {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(jwtSecret), nil
		})
	if err != nil {
		return false
	}

	if !token.Valid {
		return false
	}

	return true
}

func GenerateToken(userID int, jwtSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenLifeTime)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
