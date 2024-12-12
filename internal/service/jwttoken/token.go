package jwttoken

import (
	"errors"
	"fmt"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// ValidateToken Валидация jwt токена пользователя
func ValidateToken(tokenString string, jwtSecret string) bool {
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

// GenerateToken генерация jwt токена для пользователя
func GenerateToken(userID int64, jwtSecret string) (string, error) {
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

// GetUserIDFromToken получение userID из токена
func GetUserIDFromToken(tokenString string) (int64, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			appConfig := config.GetParams()
			return []byte(appConfig.GetJWTSecret()), nil
		})
	if err != nil {
		return 0, errors.New("failed get user_id: invalid token")
	}

	if !token.Valid {
		fmt.Println("Token is not valid")
		return 0, errors.New("failed get user_id: invalid token")
	}

	return claims.UserID, nil
}
