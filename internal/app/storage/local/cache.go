package local

var storeFile cacheStore

// GetCacheStore получение кэша с данными
func GetCacheStore() CacheStoreInterface {
	return &storeFile
}
