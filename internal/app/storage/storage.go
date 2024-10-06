package storage

import (
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage/local"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage/postgres"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

var storageList = map[string]ILinkStore{
	"local":    &local.Store,
	"postgres": &postgres.Store,
}

func GetStore() (ILinkStore, error) {
	var store ILinkStore
	appConfig := config.GetParams()
	if appConfig.GetDBConnectString() != "" {
		store = storageList["postgres"]
	} else {
		store = storageList["local"]
	}

	err := store.Init()
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("Storage init error: %s", args)
		//return nil, err
	}

	return store, nil

}
