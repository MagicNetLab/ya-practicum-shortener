package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/shortgen"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage/local"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage/postgres"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/jwttoken"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

func getShortLink(ctx context.Context, url string, userID int) (short string, httpResponseStatus int) {
	store, err := storage.GetStore()
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("error initializing storage", args)
		return "", http.StatusInternalServerError
	}

	short = shortgen.GetShortLink(7)
	httpResponseStatus = http.StatusCreated
	err = store.PutLink(ctx, url, short, userID)
	if err != nil {
		httpResponseStatus = http.StatusInternalServerError
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("error storing short link", args)
		notUniqueError := errors.Is(err, postgres.ErrLinkUniqueConflict) || errors.Is(err, local.ErrorLinkNotUnique)
		if notUniqueError {
			short, err = store.GetShort(ctx, url)
			if err == nil {
				httpResponseStatus = http.StatusConflict
			}
		}
	}

	return short, httpResponseStatus
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

func batchDeleteLinks(ctx context.Context, shorts []string, userID int) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	inputCh := make(chan string, 5)
	go func() {
		for _, short := range shorts {
			select {
			case <-doneCh:
				return
			case inputCh <- short:
			}
		}
		close(inputCh)
	}()

	numWorkers := 2
	wg := sync.WaitGroup{}
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(dataChan <-chan string) {
			defer wg.Done()

			var batch []string
			for {
				select {
				case short, ok := <-dataChan:
					if !ok {
						deleteLinks(batch, userID)
						return
					}
					batch = append(batch, short)
					if len(batch) == 5 {
						deleteLinks(batch, userID)
						batch = nil
					}
				case <-doneCh:
					return
				default:
					time.Sleep(100 * time.Millisecond)
				}
			}
		}(inputCh)
	}

	wg.Wait()
}

func deleteLinks(shorts []string, userID int) {
	if len(shorts) == 0 {
		return
	}

	store, err := storage.GetStore()
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("error initializing storage", args)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err = store.DeleteBatchLinksArray(ctx, shorts, userID)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("error deleting short links", args)
	}
}
