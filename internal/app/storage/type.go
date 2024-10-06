package storage

// ILinkStore интерфейс объекта для работы с хранилищем ссылок
type ILinkStore interface {
	// PutLink сохранение ссылки пользователя в хранилище.
	PutLink(link string, short string, userID int) error

	// PutBatchLinksArray сохранение пакета ссылок пользователя в хранилище.
	PutBatchLinksArray(StoreBatchLinksArray map[string]string, userID int) error

	// GetLink получение оригинальной ссылки по короткому хэшу.
	GetLink(short string) (string, bool, error)

	// HasShort проверка наличия коротко ссылки в хранилище
	HasShort(short string) (bool, error)

	// GetShort получение короткой ссылки из хранилища для оригинальной ссылки
	GetShort(link string) (string, error)

	// GetUserLinks получение всех ссылок пользователя из хранилища
	GetUserLinks(userID int) (map[string]string, error)

	// DeleteBatchLinksArray пометка массива ссылок пользователя как удаленных
	DeleteBatchLinksArray(shorts []string, userID int) error

	// Init инициализация хранилища
	Init() error
}
