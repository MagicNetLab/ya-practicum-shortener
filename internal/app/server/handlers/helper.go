package handlers

import (
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/shortgen"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage"
)

func generateShortLink(url string) (string, error) {
	short := shortgen.GetShortLink(7)
	store := storage.GetStore()
	err := store.PutLink(url, short)
	if err != nil {
		return "", err
	}

	return short, nil
}
