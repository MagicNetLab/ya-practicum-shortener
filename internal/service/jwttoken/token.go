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
	_, err := parseToken(tokenString, jwtSecret)
	if err != nil {
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
	appConfig := config.GetParams()
	claims, err := parseToken(tokenString, appConfig.GetJWTSecret())
	if err != nil {
		return 0, errors.New("failed get user_id: invalid token")
	}

	return claims.UserID, nil
}

func parseToken(tokenString string, jwtSecret string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
	if err != nil {
		return claims, errors.New("failed check token: invalid token")
	}

	if !token.Valid {
		return claims, errors.New("failed check token: invalid token")
	}

	return claims, nil
}
