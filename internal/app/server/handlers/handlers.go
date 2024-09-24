package handlers

import (
	"net/http"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
)

type RouteHandler struct {
	Host    string
	Method  string
	Path    string
	Handler http.HandlerFunc
}

type MapHandlers map[string]RouteHandler

func GetHandlers() MapHandlers {
	var handlers = MapHandlers{}
	c := config.GetParams()

	handlers["default"] = RouteHandler{
		Host:    c.GetDefaultHost(),
		Method:  http.MethodPost,
		Path:    "/",
		Handler: applyTokenMiddleware(applyDefaultMiddlewares(encodeHandler())),
	}
	handlers["apiDefault"] = RouteHandler{
		Host:    c.GetDefaultHost(),
		Method:  http.MethodPost,
		Path:    "/api/shorten",
		Handler: applyTokenMiddleware(applyDefaultMiddlewares(apiEncodeHandler())),
	}
	handlers["apiBatchDefault"] = RouteHandler{
		Host:    c.GetDefaultHost(),
		Method:  http.MethodPost,
		Path:    "/api/shorten/batch",
		Handler: applyTokenMiddleware(applyDefaultMiddlewares(apiBatchEncodeHandler())),
	}
	handlers["short"] = RouteHandler{
		Host:    c.GetShortHost(),
		Method:  http.MethodGet,
		Path:    "/{short}",
		Handler: applyDefaultMiddlewares(decodeHandler()),
	}
	handlers["dbPing"] = RouteHandler{
		Host:    c.GetDefaultHost(),
		Method:  http.MethodGet,
		Path:    "/ping",
		Handler: applyDefaultMiddlewares(pingHandler()),
	}
	handlers["apiUserLinks"] = RouteHandler{
		Host:    c.GetDefaultHost(),
		Method:  http.MethodGet,
		Path:    "/api/user/urls",
		Handler: applyAuthMiddleware(applyDefaultMiddlewares(apiListUserLinksHandler())),
	}
	handlers["apiDeleteLinks"] = RouteHandler{
		Host:    c.GetDefaultHost(),
		Method:  http.MethodDelete,
		Path:    "/api/user/urls",
		Handler: applyAuthMiddleware(applyDefaultMiddlewares(deleteUserLinksHandler())),
	}

	return handlers
}
