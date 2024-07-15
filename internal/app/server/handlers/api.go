package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/shortgen"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage/postgres"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

func apiEncodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shortRequest APIRequest

		if r.Method != http.MethodPost {
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Method not allowed"))
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&shortRequest); err != nil {
			http.Error(w, "Missing link", http.StatusBadRequest)
			return
		}

		conf := config.GetParams()
		status := http.StatusCreated
		short, err := generateShortLink(shortRequest.URL)
		if err != nil {
			if errors.Is(err, postgres.ErrLinkUniqueConflict) {
				short, err = getShortLink(shortRequest.URL)
				if err != nil {
					http.Error(w, "Unique conflict", http.StatusInternalServerError)
					return
				}
				status = http.StatusConflict
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		redirectLink := "http://" + conf.GetShortHost() + "/" + short
		apiResult := APIResponse{
			Result: redirectLink,
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(status)
		if err := json.NewEncoder(w).Encode(apiResult); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func apiBatchEncodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Method not allowed"))
			return
		}

		var batchRequest APIBatchRequest
		if err := json.NewDecoder(r.Body).Decode(&batchRequest); err != nil {
			logger.Log.Errorf("Failed to decode batch request: %v", err)
			http.Error(w, "Missing link", http.StatusBadRequest)
			return
		}

		store, err := storage.GetStore()
		if err != nil {
			logger.Log.Errorf("Failed to get store: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		var response APIBatchResponse
		storeData := make(map[string]string)
		conf := config.GetParams()

		for _, v := range batchRequest {
			short := shortgen.GetShortLink(7)
			row := APIBatchResponseEntity{
				CorrelationID: v.CorrelationID,
				ShortURL:      "http://" + conf.GetShortHost() + "/" + short,
			}
			storeData[short] = v.OriginalURL
			response = append(response, row)
		}

		err = store.PutBatchLinksArray(storeData)
		if err != nil {
			logger.Log.Errorf("Failed to put batch links: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
