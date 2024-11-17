package repo

import (
	"context"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/repo/memory"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/repo/postgres"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
)

var driver Driver

// Initialize инициализация репозитория
func Initialize(conf *config.Configurator) error {
	if conf.GetDBConnectString() != "" {
		driver = postgres.GetStore()
	} else {
		driver = memory.GetStore()
	}

	err := driver.Initialize(conf)
	if err != nil {
		return err
	}

	return nil
}

// Close закрывает репозиторий
func Close() error {
	return driver.Close()
}

// PutLink сохраняет ссылку в хранилище
func PutLink(ctx context.Context, link string, short string, userID int) error {
	return driver.PutLink(ctx, link, short, userID)
}

// PutBatchLinksArray сохранение пакета ссылок пользователя в хранилище.
func PutBatchLinksArray(ctx context.Context, StoreBatchLinksArray map[string]string, userID int) error {
	return driver.PutBatchLinksArray(ctx, StoreBatchLinksArray, userID)
}

// GetLink получение оригинальной ссылки по короткому хэшу.
func GetLink(ctx context.Context, short string) (string, bool, error) {
	return driver.GetLink(ctx, short)
}

// HasShort проверка наличия коротко ссылки в хранилище
func HasShort(ctx context.Context, short string) (bool, error) {
	return driver.HasShort(ctx, short)
}

// GetShort получение короткой ссылки из хранилища для оригинальной ссылки
func GetShort(ctx context.Context, link string) (string, error) {
	return driver.GetShort(ctx, link)
}

// GetUserLinks получение всех ссылок пользователя из хранилища
func GetUserLinks(ctx context.Context, userID int) (map[string]string, error) {
	return driver.GetUserLinks(ctx, userID)
}

// DeleteBatchLinksArray пометка массива ссылок пользователя как удаленных
func DeleteBatchLinksArray(ctx context.Context, shorts []string, userID int) error {
	return driver.DeleteBatchLinksArray(ctx, shorts, userID)
}
