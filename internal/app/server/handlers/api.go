package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/repo"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/repo/postgres"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/shortgen"
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
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed get user_id from cookie", args)
			http.Error(w, "incorrect user token", http.StatusBadRequest)
			return
		}

		var shortRequest APIRequest
		if err := json.NewDecoder(r.Body).Decode(&shortRequest); err != nil {
			http.Error(w, "Missing link", http.StatusBadRequest)
			return
		}

		if shortRequest.URL == "" {
			http.Error(w, "Missing link", http.StatusBadRequest)
			return
		}

		c := config.GetParams()
		apiResult := APIResponse{Result: ""}
		short, status := getShortLink(r.Context(), shortRequest.URL, userID)
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
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed get user_id from cookie", args)
			http.Error(w, "incorrect user token", http.StatusBadRequest)
			return
		}

		var batchRequest APIBatchRequest
		if err := json.NewDecoder(r.Body).Decode(&batchRequest); err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed to decode batch request", args)
			http.Error(w, "Missing link", http.StatusBadRequest)
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

		err = repo.PutBatchLinksArray(r.Context(), storeData, userID)
		if err != nil {
			if errors.Is(err, postgres.ErrLinkUniqueConflict) {
				http.Error(w, "Conflict: one or more links are not unique", http.StatusConflict)
				return
			}

			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed put batch links", args)
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
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed get user_id from cookie", args)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		res, err := repo.GetUserLinks(r.Context(), userID)
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed get user links", args)
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
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed get user_id from cookie", args)
			http.Error(w, "incorrect user token", http.StatusBadRequest)
			return
		}

		var deleteRequest APIDeleteRequest
		if err := json.NewDecoder(r.Body).Decode(&deleteRequest); err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed to decode delete request", args)
			http.Error(w, "Incorrect request data", http.StatusBadRequest)
			return
		}

		go batchDeleteLinks(r.Context(), deleteRequest, userID)

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		if err := json.NewEncoder(w).Encode("ok"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func apiInternalStatsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
		defer cancel()

		urls, err := repo.GetLinksCount(ctx)
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed get links count", args)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		users, err := repo.GetUsersCount(ctx)
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed get users count", args)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		response := APIInternalStatsResponse{Urls: urls, Users: users}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
