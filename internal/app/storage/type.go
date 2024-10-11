package storage

import "context"

// Storer интерфейс объекта для работы с хранилищем ссылок
type Storer interface {
	// PutLink сохранение ссылки пользователя в хранилище.
	PutLink(ctx context.Context, link string, short string, userID int) error

	// PutBatchLinksArray сохранение пакета ссылок пользователя в хранилище.
	PutBatchLinksArray(ctx context.Context, StoreBatchLinksArray map[string]string, userID int) error

	// GetLink получение оригинальной ссылки по короткому хэшу.
	GetLink(ctx context.Context, short string) (string, bool, error)

	// HasShort проверка наличия коротко ссылки в хранилище
	HasShort(ctx context.Context, short string) (bool, error)

	// GetShort получение короткой ссылки из хранилища для оригинальной ссылки
	GetShort(ctx context.Context, link string) (string, error)

	// GetUserLinks получение всех ссылок пользователя из хранилища
	GetUserLinks(ctx context.Context, userID int) (map[string]string, error)

	// DeleteBatchLinksArray пометка массива ссылок пользователя как удаленных
	DeleteBatchLinksArray(ctx context.Context, shorts []string, userID int) error

	// Init инициализация хранилища
	Init() error
}
