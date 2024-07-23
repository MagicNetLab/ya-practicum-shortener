package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage/postgres"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
	"github.com/go-chi/chi/v5"
)

func encodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		userID, err := parseCookie(r)
		if err != nil {
			logger.Log.Errorf("failed get user from token: %v", err)
			http.Error(w, "incorrect user token", http.StatusBadRequest)
			return
		}

		link, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Log.Errorf("Error reading body: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if string(link) == "" {
			http.Error(w, "Missing link", http.StatusBadRequest)
			return
		}

		c := config.GetParams()
		short, status := getShortLink(string(link), userID)
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(status)
		_, err = w.Write([]byte(fmt.Sprintf("http://%s/%s", c.GetShortHost(), short)))
		if err != nil {
			logger.Log.Errorf("Failed to write response %v", err)
		}
	}
}

func decodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		short := chi.URLParam(r, "short")
		store, err := storage.GetStore()
		if err != nil {
			logger.Log.Errorf("Failed to get store %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		hasShort, err := store.HasShort(short)
		if err != nil {
			logger.Log.Errorf("Failed to check short %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		if hasShort {
			link, err := store.GetLink(short)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			http.Redirect(w, r, link, http.StatusTemporaryRedirect)

			return
		}

		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func pingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isPostgresOK := postgres.Ping()
		logger.Log.Infoln(isPostgresOK)
		if !isPostgresOK {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}
}
