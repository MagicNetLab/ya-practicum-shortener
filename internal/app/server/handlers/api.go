package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage/postgres"
	"net/http"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/shortgen"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

func apiEncodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusForbidden)
			return
		}

		userID, err := parseCookie(r)
		if err != nil {
			logger.Log.Errorf("failed get user_id from cookie: %v", err)
			http.Error(w, "incorrect user token", http.StatusBadRequest)
			return
		}

		var shortRequest APIRequest
		if err := json.NewDecoder(r.Body).Decode(&shortRequest); err != nil {
			http.Error(w, "Missing link", http.StatusBadRequest)
			return
		}

		c := config.GetParams()
		apiResult := APIResponse{Result: ""}
		short, status := getShortLink(shortRequest.URL, userID)
		apiResult.Result = "http://" + c.GetShortHost() + "/" + short
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
			http.Error(w, "Method not allowed", http.StatusForbidden)
			return
		}

		userID, err := parseCookie(r)
		if err != nil {
			logger.Log.Errorf("failed get user_id from cookie: %v", err)
			http.Error(w, "incorrect user token", http.StatusBadRequest)
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
			return
		}

		var response APIBatchResponse
		storeData := make(map[string]string)
		c := config.GetParams()
		for _, v := range batchRequest {
			short := shortgen.GetShortLink(7)
			row := APIBatchResponseEntity{
				CorrelationID: v.CorrelationID,
				ShortURL:      fmt.Sprintf("http://%s/%s", c.GetShortHost(), short),
			}
			storeData[short] = v.OriginalURL
			response = append(response, row)
		}

		err = store.PutBatchLinksArray(storeData, userID)
		if err != nil {
			if errors.Is(err, postgres.ErrLinkUniqueConflict) {
				http.Error(w, "Conflict: one or more links are not unique", http.StatusConflict)
				return
			}

			logger.Log.Errorf("Failed to put batch links: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func apiListUserLinksHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, err := parseCookie(r)
		if err != nil {
			logger.Log.Errorf("failed get user_id from cookie: %v", err)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		store, err := storage.GetStore()
		if err != nil {
			logger.Log.Errorf("Failed to get store: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		res, err := store.GetUserLinks(userID)
		if err != nil {
			logger.Log.Errorf("Failed to get user links: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		userLinksResponse := UserLinksResponse{}
		if len(res) == 0 {
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusNoContent)
			if err := json.NewEncoder(w).Encode(userLinksResponse); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		c := config.GetParams()
		for k, v := range res {
			row := UserLinkEntity{
				ShortURL:    fmt.Sprintf("http://%s/%s", c.GetShortHost(), k),
				OriginalURL: v,
			}
			userLinksResponse = append(userLinksResponse, row)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(userLinksResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func deleteUserLinksHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, err := parseCookie(r)
		if err != nil {
			logger.Log.Errorf("failed get user_id from cookie: %v", err)
			http.Error(w, "incorrect user token", http.StatusBadRequest)
			return
		}

		var deleteRequest APIDeleteRequest
		if err := json.NewDecoder(r.Body).Decode(&deleteRequest); err != nil {
			logger.Log.Errorf("Failed to decode delete request: %v", err)
			http.Error(w, "Incorrect request data", http.StatusBadRequest)
			return
		}

		go batchDeleteLinks(deleteRequest, userID)

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		if err := json.NewEncoder(w).Encode("ok"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
