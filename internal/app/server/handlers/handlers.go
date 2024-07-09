package handlers

import (
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
		Handler: applyMiddlewares(encodeHandler()),
	}
	handlers["apiDefault"] = RouteHandler{
		Method:  http.MethodPost,
		Path:    "/api/shorten",
		Handler: applyMiddlewares(apiEncodeHandler()),
	}
	handlers["short"] = RouteHandler{
		Method:  http.MethodGet,
		Path:    "/{short}",
		Handler: applyMiddlewares(decodeHandler()),
	}
	handlers["dbPing"] = RouteHandler{
		Method:  http.MethodGet,
		Path:    "/ping",
		Handler: applyMiddlewares(pingHandler()),
	}

	return handlers
}
