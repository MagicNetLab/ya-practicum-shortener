package handlers

import (
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/compression"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
	"net/http"
)

type middlewares []func(h http.HandlerFunc) http.HandlerFunc

var middlewareList = middlewares{
	compression.GzipMiddleware,
	logger.Middleware,
}

func requestMiddlewares(h http.HandlerFunc) http.HandlerFunc {
	for _, m := range middlewareList {
		h = m(h)
	}
	return h
}
