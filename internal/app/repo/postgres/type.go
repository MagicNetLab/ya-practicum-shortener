package postgres

import (
	"errors"
)

const (
	dbUser      = "DB_USER"
	dbPassword  = "DB_PASSWORD"
	dbHost      = "DB_HOST"
	dbPort      = "DB_PORT"
	dbName      = "DB_NAME"
	dbDBSslMode = "DB_SSL_MODE"
)

const (
	insertLinkSQL             = "INSERT INTO links (short, link, user_id) VALUES ($1, $2, $3)"
	selectLinkSQL             = "SELECT link, is_deleted FROM links WHERE short = $1"
	selectShortSQL            = "SELECT short FROM links WHERE link = $1"
	hasLinkSQL                = "SELECT count(*) FROM links WHERE short = $1"
	selectUserLinks           = "SELECT short, link FROM links WHERE user_id = $1"
	selectLinksCount          = "SELECT count(*) as linksCount FROM links WHERE is_deleted = false"
	selectUsersCount          = "SELECT count(distinct(links.user_id)) FROM links"
	existsUserLogin           = "SELECT count(id) FROM users WHERE login = $1"
	getUserIDByLoginAndSecret = "SELECT id FROM users WHERE login = $1 AND secret = $2"
	createUser                = "INSERT INTO users (login, secret) VALUES ($1, $2)"
)

// ErrLinkUniqueConflict ошибка попытки записать не уникальную ссылку
var ErrLinkUniqueConflict = errors.New("url is not unique")
