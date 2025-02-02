package jwttoken

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/user"
)

// CheckAuthMiddleware миддлвара для проверки авторизации пользователя.
func CheckAuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appConfig := config.GetParams()
		claims := &Claims{}
		cookie, err := r.Cookie("token")
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed get token from cookie", args)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, claims,
			func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}

				return []byte(appConfig.GetJWTSecret()), nil
			})

		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Info("failed to parse token", args)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !token.Valid || claims.UserID == 0 {
			logger.Error("invalid token", nil)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	}
}

// TokenMiddleware миддлвара для установки токена пользователя.
// Тестовый вариант пока нет методов для регистрации и авторизации.
func TokenMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appConfig := config.GetParams()
		cookie, err := r.Cookie("token")
		if err != nil || !ValidateToken(cookie.Value, appConfig.GetJWTSecret()) {
			userID, err := user.Create("", "")
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			token, err := GenerateToken(userID, appConfig.GetJWTSecret())
			if err != nil {
				args := map[string]interface{}{"error": err.Error()}
				logger.Error("failed to generate token", args)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			newCookie := http.Cookie{Name: "token", Value: token, Path: "/", Expires: time.Now().Add(TokenLifeTime)}
			r.AddCookie(&newCookie)
			http.SetCookie(w, &newCookie)
		}

		h.ServeHTTP(w, r)
	}
}
