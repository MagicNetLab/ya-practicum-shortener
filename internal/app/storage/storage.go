package storage

import (
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage/local"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage/postgres"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
)

var storageList = map[string]LinkStore{
	"local":    &local.Store,
	"postgres": &postgres.Store,
}

func GetStore() (LinkStore, error) {
	var store LinkStore
	appConfig := config.GetParams()
	if appConfig.HasDBConnectParams() {
		store = storageList["postgres"]
	} else {
		store = storageList["local"]
	}

	err := store.Init()
	if err != nil {
		return nil, err
	}

	return store, nil

}
