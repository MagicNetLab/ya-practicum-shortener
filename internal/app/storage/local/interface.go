package local

// CacheStoreInterface интерфейс кэша локального хранилища ссылок в файле json
type CacheStoreInterface interface {
	// Load загрузка данных из кэша в память
	Load() ([]StoreEntity, error)

	// Save сохранение ссылки в кэше
	Save(link linkEntity) error

	// IsInitialized проверка инициализирован кэш или нет
	IsInitialized() bool

	// SetPath установка пути до локального файла с кэшем
	SetPath(path string)

	// SetInitialized установка параметра инициализации хранилища
	SetInitialized(isInitialized bool)
}
