package storage

import "github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage/local"

type LinkStore interface {
	PutLink(link string, short string) error

	GetLink(short string) (string, error)

	HasShort(short string) bool
}

var storageList = map[string]LinkStore{
	"local": &local.Store,
}

func GetStore() LinkStore {
	// TODO тип хранилища получать из конфига из конфига
	return storageList["local"]
}
