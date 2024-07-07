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
		Handler: requestMiddlewares(encodeHandler()),
	}
	handlers["apiDefault"] = RouteHandler{
		Method:  http.MethodPost,
		Path:    "/api/shorten",
		Handler: requestMiddlewares(apiEncodeHandler()),
	}
	handlers["short"] = RouteHandler{
		Method:  http.MethodGet,
		Path:    "/{short}",
		Handler: requestMiddlewares(decodeHandler()),
	}

	return handlers
}
