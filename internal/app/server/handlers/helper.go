package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/shortgen"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage/postgres"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/jwttoken"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
	"github.com/golang-jwt/jwt/v4"
)

func getShortLink(url string, userID int) (string, int) {
	store, err := storage.GetStore()
	if err != nil {
		logger.Log.Errorf("Error init storage: %v", err)
		return "", http.StatusInternalServerError
	}

	short := shortgen.GetShortLink(7)
	status := http.StatusCreated
	err = store.PutLink(url, short, userID)
	if err != nil {
		status = http.StatusInternalServerError
		logger.Log.Errorf("Error putting short link: %v", err)
		if errors.Is(err, postgres.ErrLinkUniqueConflict) {
			short, err = store.GetShort(url)
			if err == nil {
				status = http.StatusConflict
			}
		}
	}

	return short, status
}

func getUserID(tokenString string) (int, error) {
	claims := &jwttoken.Claims{}
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

func parseCookie(r *http.Request) (int, error) {
	cookie, err := r.Cookie("token")
	if err != nil || cookie.Value == "" {
		return 0, fmt.Errorf("failed parse cookie %v", err)
	}

	userID, err := getUserID(cookie.Value)
	if err != nil {
		return 0, fmt.Errorf("failed get userID from cookie %v", err)
	}

	return userID, nil
}
