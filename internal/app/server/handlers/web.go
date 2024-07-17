package handlers

import (
	"errors"
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
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Method not allowed"))

			return
		}

		link, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if string(link) == "" {
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("Missing link"))
			if err != nil {
				logger.Log.Errorf("Failed to write response %v", err)
			}

			return
		}

		c := config.GetParams()
		// todo избавиться от дублирования с api.go
		short, err := generateShortLink(string(link))
		if err == nil {
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			_, err = w.Write([]byte(fmt.Sprintf("http://%s/%s", c.GetShortHost(), short)))
			if err != nil {
				logger.Log.Errorf("Failed to write response %v", err)
			}
			return
		}

		if errors.Is(err, postgres.ErrLinkUniqueConflict) {
			short, err = getShortLink(string(link))
			if err != nil {
				logger.Log.Errorf("Failed to get short link %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusConflict)
			_, err = w.Write([]byte(fmt.Sprintf("http://%s/%s", c.GetShortHost(), short)))
			if err != nil {
				logger.Log.Errorf("Failed to write response %v", err)
			}
			return
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func decodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusForbidden)
			_, err := w.Write([]byte("Method not allowed"))
			if err != nil {
				logger.Log.Errorf("Failed to write response %v", err)
			}

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
				http.NotFound(w, r)
				return
			}

			http.Redirect(w, r, link, http.StatusTemporaryRedirect)

			return
		}

		http.NotFound(w, r)
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
