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

type handlerParams struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

type MapHandlers map[string]handlerParams

var handlers = MapHandlers{}

func GetHandlers() MapHandlers {

	handlers["default"] = handlerParams{
		Method:  http.MethodPost,
		Path:    "/",
		Handler: encodeLinkHeader,
	}

	handlers["short"] = handlerParams{
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
		body := "Method not allowed"
		_, err := response.Write([]byte(body))
		if err != nil {
			panic(err)
		}
		return

	}

	link, err := io.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	if string(link) == "" {
		response.Header().Set("content-type", "text/plain")
		response.WriteHeader(http.StatusBadRequest)
		body := "Missing link"
		_, err := response.Write([]byte(body))
		if err != nil {
			panic(err)
		}
		return
	}

	short := shortgen.GetShort(7)
	store := storage.GetStore()
	err = store.PutLink(string(link), short)
	if err != nil {
		panic(err)
	}

	conf := config.GetParams()
	response.Header().Set("content-type", "text/plain")
	response.WriteHeader(http.StatusCreated)
	body := fmt.Sprintf("http://%s/%s", conf.GetShortHost(), short)
	_, err = response.Write([]byte(body))
	if err != nil {
		panic(err)
	}
}

func decodeLinkHeader(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		resp.Header().Set("content-type", "text/plain")
		resp.WriteHeader(http.StatusForbidden)
		body := "Method not allowed"
		_, err := resp.Write([]byte(body))
		if err != nil {
			panic(err)
		}
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
