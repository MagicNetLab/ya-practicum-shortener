package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/repo"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/repo/postgres"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

func encodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		userID, err := parseCookie(r)
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed get user id from cookie", args)
			http.Error(w, "incorrect user token", http.StatusBadRequest)
			return
		}

		link, err := io.ReadAll(r.Body)
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed get link from request body", args)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if string(link) == "" {
			http.Error(w, "Missing link", http.StatusBadRequest)
			return
		}

		c := config.GetParams()
		short, status := getShortLink(r.Context(), string(link), userID)
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(status)
		w.Write([]byte(fmt.Sprintf("http://%s/%s", c.GetShortHost(), short)))
	}
}

func decodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		short := chi.URLParam(r, "short")

		link, isDeleted, err := repo.GetLink(r.Context(), short)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if isDeleted {
			http.Error(w, http.StatusText(http.StatusGone), http.StatusGone)
			return
		}

		http.Redirect(w, r, link, http.StatusTemporaryRedirect)
	}
}

func pingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isPostgresOK := postgres.Ping()
		if !isPostgresOK {
			logger.Error("failed to ping postgres", nil)
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}
}
