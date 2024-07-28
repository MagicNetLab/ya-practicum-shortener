package jwttoken

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/user"
	"github.com/golang-jwt/jwt/v4"
)

func CheckAuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appConfig := config.GetParams()
		claims := &Claims{}
		cookie, err := r.Cookie("token")
		if err != nil {
			logger.Log.Errorf("Cookie err: %v", err)
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
			logger.Log.Infof("failed to parse token: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !token.Valid || claims.UserID == 0 {
			logger.Log.Infof("invalid token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	}
}

func TokenMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appConfig := config.GetParams()
		cookie, err := r.Cookie("token")
		if err != nil || !validateToken(cookie.Value, appConfig.GetJWTSecret()) {
			u := user.Create()
			token, err := GenerateToken(u.ID, appConfig.GetJWTSecret())
			if err != nil {
				logger.Log.Errorf("failed to generate token: %v", err)
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
