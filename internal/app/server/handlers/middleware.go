package handlers

import (
	"net/http"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/compression"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/jwttoken"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

const trustedSubnetHeader = "X-Real-IP"

type middlewareList []func(handlerFunc http.HandlerFunc) http.HandlerFunc

var middlewares = middlewareList{
	compression.GzipMiddleware,
	logger.Middleware,
}

func applyDefaultMiddlewares(h http.HandlerFunc) http.HandlerFunc {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

func applyAuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return jwttoken.CheckAuthMiddleware(h)
}

func applyTokenMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return jwttoken.TokenMiddleware(h)
}

func applyTrustedMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return trustedAccessMiddleware(h)
}

func trustedAccessMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conf := config.GetParams()
		trustedSubnet := conf.GetTrustedSubnet()
		header := r.Header.Get(trustedSubnetHeader)
		if header == "" || header != trustedSubnet {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, r)
	}
}
