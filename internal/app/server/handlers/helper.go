package handlers

import (
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/shortgen"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

func generateShortLink(url string) (string, error) {
	short := shortgen.GetShortLink(7)
	store, err := storage.GetStore()
	if err != nil {
		logger.Log.Errorf("Error init storage: %v", err)
		return "", err
	}

	err = store.PutLink(url, short)
	if err != nil {
		return "", err
	}

	return short, nil
}
