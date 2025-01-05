package repo

import (
	"context"
	"errors"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
)

// ErrorLinkNotUnique ошибка в случае попытки сохранения уже существующей в базе ссылки
var ErrorLinkNotUnique = errors.New("link not unique")

// Repository интерфейс драйвера хранения данных
type Repository interface {
	// HasUserLogin проверка занятости логина пользователя
	HasUserLogin(ctx context.Context, username string) (bool, error)

	// AuthUser аутентификация пользователя
	AuthUser(ctx context.Context, username string, secret string) (int64, error)

	// CreateUser создание пользователя
	CreateUser(ctx context.Context, username string, secret string) (bool, error)

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

	// GetLinksCount возвращает количество сокращенных ссылок в системе
	GetLinksCount(ctx context.Context) (int, error)

	// GetUsersCount возвращает количество пользователей в системе
	GetUsersCount(ctx context.Context) (int, error)

	// Initialize инициализация хранилища
	Initialize(config *config.Configurator) error

	// Close Закрывает хранилище
	Close() error
}
