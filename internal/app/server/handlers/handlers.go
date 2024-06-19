package handlers

import (
	"fmt"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/shortgen"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type RouteHandler struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

type MapHandlers map[string]RouteHandler

func GetHandlers() MapHandlers {
	var handlers = MapHandlers{}

	handlers["default"] = RouteHandler{
		Method:  http.MethodPost,
		Path:    "/",
		Handler: logger.RequestLogger(encodeHandler()),
	}

	handlers["short"] = RouteHandler{
		Method:  http.MethodGet,
		Path:    "/{short}",
		Handler: logger.RequestLogger(decodeHandler()),
	}

	return handlers
}

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

		short := shortgen.GetShortLink(7)
		store := storage.GetStore()
		err = store.PutLink(string(link), short)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)

			return
		}

		conf := config.GetParams()
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(fmt.Sprintf("http://%s/%s", conf.GetShortHost(), short)))
		if err != nil {
			logger.Log.Errorf("Failed to write response %v", err)
		}
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
		store := storage.GetStore()

		if store.HasShort(short) {
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
