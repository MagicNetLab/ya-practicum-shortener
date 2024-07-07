package local

type StoreEntity struct {
	Short string `json:"short"`
	Link  string `json:"link"`
}

type CacheStoreInterface interface {
	Load() ([]StoreEntity, error)
	Save(short string, link string) error
	IsInitialized() bool
	SetPath(path string)
	SetInitialized(isInitialized bool)
}
