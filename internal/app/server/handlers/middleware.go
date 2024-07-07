package handlers

import (
	"net/http"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/compression"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

type middlewareList []func(handlerFunc http.HandlerFunc) http.HandlerFunc

var middlewares = middlewareList{
	compression.GzipMiddleware,
	logger.Middleware,
}

func applyMiddlewares(h http.HandlerFunc) http.HandlerFunc {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}
