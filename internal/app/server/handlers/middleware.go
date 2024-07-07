package handlers

import (
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/compression"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
	"net/http"
)

func requestMiddlewares(h http.HandlerFunc) http.HandlerFunc {
	return logger.Middleware(compression.Middleware(h))
}
