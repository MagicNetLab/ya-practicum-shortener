package handlers

import (
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
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
	handlers["apiDefault"] = RouteHandler{
		Method:  http.MethodPost,
		Path:    "/api/shorten",
		Handler: logger.RequestLogger(apiEncodeHandler()),
	}
	handlers["short"] = RouteHandler{
		Method:  http.MethodGet,
		Path:    "/{short}",
		Handler: logger.RequestLogger(decodeHandler()),
	}

	return handlers
}
