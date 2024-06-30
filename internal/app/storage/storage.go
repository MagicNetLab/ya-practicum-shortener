package storage

import "github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage/local"

type LinkStore interface {
	PutLink(link string, short string) error

	GetLink(short string) (string, error)

	HasShort(short string) bool

	Init()
}

var storageList = map[string]LinkStore{
	"local": &local.Store,
}

func GetStore() LinkStore {
	store := storageList["local"]
	store.Init()
	return store
}
