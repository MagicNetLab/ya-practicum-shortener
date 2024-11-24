package handlers

import (
	"net/http"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
)

// RouteHandler объект с информацией необходимой для инициализации роута
type RouteHandler struct {
	Host    string
	Method  string
	Path    string
	Handler http.HandlerFunc
}

// MapHandlers массив с данными роутов для запуска приложения
type MapHandlers map[string]RouteHandler

// GetHandlers возвращает массив с данными роутов для старта приложения
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
	handlers["serverStats"] = RouteHandler{
		Host:    c.GetDefaultHost(),
		Method:  http.MethodGet,
		Path:    "/api/internal/stats",
		Handler: applyTrustedMiddleware(applyDefaultMiddlewares(apiInternalStatsHandler())),
	}

	return handlers
}
