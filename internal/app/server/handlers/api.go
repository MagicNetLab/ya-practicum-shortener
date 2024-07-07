package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
)

type APIRequest struct {
	URL string `json:"url"`
}

type APIResponse struct {
	Result string `json:"result"`
}

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

		short, err := generateShortLink(shortRequest.URL)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		conf := config.GetParams()
		redirectLink := "http://" + conf.GetShortHost() + "/" + short
		apiResult := APIResponse{
			Result: redirectLink,
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(apiResult); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
