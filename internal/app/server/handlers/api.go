package handlers

import (
	"encoding/json"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/shortgen"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"net/http"
)

type ApiRequest struct {
	Url string `json:"url"`
}

type ApiResponse struct {
	Result string `json:"result"`
}

var shortRequest ApiRequest

func apiEncodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		short := shortgen.GetShortLink(7)
		store := storage.GetStore()
		err := store.PutLink(shortRequest.Url, short)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		conf := config.GetParams()
		redirectLink := "http://" + conf.GetShortHost() + "/" + short
		apiResult := ApiResponse{
			Result: redirectLink,
		}

		resp, err := json.Marshal(apiResult)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(resp)
	}
}
