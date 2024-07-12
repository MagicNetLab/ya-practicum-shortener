package storage

import (
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage/local"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage/postgres"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

var storageList = map[string]LinkStore{
	"local":    &local.Store,
	"postgres": &postgres.Store,
}

func GetStore() (LinkStore, error) {
	var store LinkStore
	appConfig := config.GetParams()
	if appConfig.GetDBConnectString() != "" {
		store = storageList["postgres"]
	} else {
		store = storageList["local"]
	}

	err := store.Init()
	if err != nil {
		logger.Log.Infof("Storage init error: %s", err.Error())
		return nil, err
	}

	return store, nil

}
