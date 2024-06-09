package handlers

import (
	"fmt"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/server/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/shortgen"
	storage "github.com/MagicNetLab/ya-practicum-shortener/internal/app/store"
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
		Handler: encodeLinkHeader,
	}

	handlers["short"] = RouteHandler{
		Method:  http.MethodGet,
		Path:    "/{short}",
		Handler: decodeLinkHeader,
	}

	return handlers
}

func encodeLinkHeader(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		response.Header().Set("content-type", "text/plain")
		response.WriteHeader(http.StatusForbidden)
		response.Write([]byte("Method not allowed"))
		return

	}

	link, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if string(link) == "" {
		response.Header().Set("content-type", "text/plain")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Missing link"))
		return
	}

	short := shortgen.GetShortLink(7)
	store := storage.GetStore()
	err = store.PutLink(string(link), short)
	if err != nil {
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	conf := config.GetParams()
	response.Header().Set("content-type", "text/plain")
	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(fmt.Sprintf("http://%s/%s", conf.GetShortHost(), short)))
}

func decodeLinkHeader(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		resp.Header().Set("content-type", "text/plain")
		resp.WriteHeader(http.StatusForbidden)
		resp.Write([]byte("Method not allowed"))
		return
	}

	short := chi.URLParam(req, "short")
	store := storage.GetStore()

	if store.HasShort(short) {
		link, err := store.GetLink(short)
		if err != nil {
			http.NotFound(resp, req)
		}

		http.Redirect(resp, req, link, http.StatusTemporaryRedirect)
	}

	http.NotFound(resp, req)
}
