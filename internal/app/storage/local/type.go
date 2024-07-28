package local

type StoreEntity struct {
	Short     string `json:"short"`
	Link      string `json:"link"`
	UserID    int    `json:"user_id"`
	IsDeleted bool   `json:"is_deleted"`
}

type CacheStoreInterface interface {
	Load() ([]StoreEntity, error)
	Save(link linkEntity) error
	IsInitialized() bool
	SetPath(path string)
	SetInitialized(isInitialized bool)
}
