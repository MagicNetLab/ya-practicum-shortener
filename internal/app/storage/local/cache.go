package local

var storeFile cacheStore

// getCacheStore получение кэша с данными
func getCacheStore() *cacheStore {
	return &storeFile
}
