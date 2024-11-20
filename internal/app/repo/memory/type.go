package memory

import "errors"

// linkEntity структура данных для хранения ссылок.
type linkEntity struct {
	userID      int
	shortLink   string
	originalURL string
	isDeleted   bool
}

// StoreEntity структура данных для хранения в кэше
type StoreEntity struct {
	Short     string `json:"short"`
	Link      string `json:"link"`
	UserID    int    `json:"user_id"`
	IsDeleted bool   `json:"is_deleted"`
}

// ErrorLinkNotUnique ошибка в случае попытки сохранения уже существующей в базе ссылки
var ErrorLinkNotUnique = errors.New("link not unique")
