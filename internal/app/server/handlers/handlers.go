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
		Handler: applyTokenMiddleware(applyDefaultMiddlewares(encodeHandler())),
	}
	handlers["apiDefault"] = RouteHandler{
		Method:  http.MethodPost,
		Path:    "/api/shorten",
		Handler: applyTokenMiddleware(applyDefaultMiddlewares(apiEncodeHandler())),
	}
	handlers["apiBatchDefault"] = RouteHandler{
		Method:  http.MethodPost,
		Path:    "/api/shorten/batch",
		Handler: applyTokenMiddleware(applyDefaultMiddlewares(apiBatchEncodeHandler())),
	}
	handlers["short"] = RouteHandler{
		Method:  http.MethodGet,
		Path:    "/{short}",
		Handler: applyDefaultMiddlewares(decodeHandler()),
	}
	handlers["dbPing"] = RouteHandler{
		Method:  http.MethodGet,
		Path:    "/ping",
		Handler: applyDefaultMiddlewares(pingHandler()),
	}
	handlers["apiUserLinks"] = RouteHandler{
		Method:  http.MethodGet,
		Path:    "/api/user/urls",
		Handler: applyDefaultMiddlewares(applyAuthMiddleware(apiListUserLinksHandler())),
	}

	return handlers
}
