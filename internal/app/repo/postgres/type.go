package postgres

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
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
	insertLinkSQL   = "INSERT INTO links (short, link, user_id) VALUES ($1, $2, $3)"
	selectLinkSQL   = "SELECT link, is_deleted FROM links WHERE short = $1"
	selectShortSQL  = "SELECT short FROM links WHERE link = $1"
	hasLinkSQL      = "SELECT count(*) FROM links WHERE short = $1"
	selectUserLinks = "SELECT short, link FROM links WHERE user_id = $1"
)

// ErrLinkUniqueConflict ошибка попытки записать не уникальную ссылку
var ErrLinkUniqueConflict = errors.New("url is not unique")

// Store объект хранилища
type Store struct {
	conn    *pgx.Conn
	pool    *pgxpool.Pool
	connStr string
}

// PutLink сохранение ссылки пользователя в хранилище.
func (s *Store) PutLink(ctx context.Context, link string, short string, userID int) error {
	result, err := s.conn.Exec(ctx, insertLinkSQL, short, link, userID)
	if err != nil {
		if strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
			return ErrLinkUniqueConflict
		}
		return err
	}

	if result.RowsAffected() != 1 {
		return errors.New("filed save data")
	}

	return nil
}

// PutBatchLinksArray сохранение пакета ссылок пользователя в хранилище.
func (s *Store) PutBatchLinksArray(ctx context.Context, StoreBatchLinksArray map[string]string, userID int) error {
	transaction, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer transaction.Rollback(ctx)

	_, err = transaction.Prepare(ctx, "batch-insert", insertLinkSQL)
	if err != nil {
		return err
	}

	for key, value := range StoreBatchLinksArray {
		cTag, err := transaction.Exec(ctx, "batch-insert", key, value, userID)
		if err != nil {
			if strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
				return ErrLinkUniqueConflict
			}
			return err
		}

		if cTag.RowsAffected() != 1 {
			return errors.New("filed save data: " + cTag.String())
		}
	}

	transaction.Commit(ctx)

	return nil
}

// GetLink получение оригинальной ссылки по короткому хэшу.
func (s *Store) GetLink(ctx context.Context, short string) (string, bool, error) {
	var link string
	var isDeleted bool
	err := s.conn.QueryRow(ctx, selectLinkSQL, short).Scan(&link, &isDeleted)
	if err != nil {
		return "", isDeleted, errors.New("database error: " + err.Error())
	}

	return link, isDeleted, nil
}

// HasShort проверка наличия коротко ссылки в хранилище
func (s *Store) HasShort(ctx context.Context, short string) (bool, error) {
	var count int
	err := s.conn.QueryRow(ctx, hasLinkSQL, short).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

// GetShort получение короткой ссылки из хранилища для оригинальной ссылки
func (s *Store) GetShort(ctx context.Context, link string) (string, error) {
	var short string
	err := s.conn.QueryRow(ctx, selectShortSQL, link).Scan(&short)
	if err != nil {
		return "", errors.New("database error: " + err.Error())
	}

	return short, nil
}

// GetUserLinks получение всех ссылок пользователя из хранилища
func (s *Store) GetUserLinks(ctx context.Context, userID int) (map[string]string, error) {
	res := make(map[string]string)

	rows, err := s.conn.Query(ctx, selectUserLinks, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return res, nil
		}

		return nil, errors.New("database error: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var short, link string
		if err := rows.Scan(&short, &link); err != nil {
			rows.Close()
			return nil, errors.New("database error: " + err.Error())
		}
		res[short] = link
	}
	rows.Close()

	return res, nil
}

// DeleteBatchLinksArray пометка массива ссылок пользователя как удаленных
func (s *Store) DeleteBatchLinksArray(ctx context.Context, shorts []string, userID int) error {
	var paramrefs string
	var ids []interface{}
	ids = append(ids, userID)
	for i, v := range shorts {
		paramrefs += `$` + strconv.Itoa(i+2) + `,`
		ids = append(ids, v)
	}
	paramrefs = paramrefs[:len(paramrefs)-1]

	sqlQuery := `UPDATE links SET is_deleted = true WHERE user_id = $1 AND short IN (` + paramrefs + `)`
	exec, err := s.pool.Exec(ctx, sqlQuery, ids...)
	if err != nil {
		return err
	}

	if exec.RowsAffected() < 1 {
		return errors.New("filed delete data: " + exec.String())
	}

	return nil
}

// Initialize инициализация хранилища
func (s *Store) Initialize(config config.AppConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var connectStr string

	connectParams, err := parseConnectString(config.GetDBConnectString())
	if err == nil {
		connectStr = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			connectParams[dbUser],
			connectParams[dbPassword],
			connectParams[dbHost],
			connectParams[dbPort],
			connectParams[dbName],
			connectParams[dbDBSslMode])
	} else {
		args := map[string]interface{}{"connect": config.GetDBConnectString()}
		logger.Error("failed parse connect string", args)
		connectStr = config.GetDBConnectString()
	}

	conn, err := pgx.Connect(ctx, connectStr)
	if err != nil {
		return errors.New("database connection error: " + err.Error())
	}

	err = conn.Ping(ctx)
	if err != nil {
		return errors.New("database ping error: " + err.Error())
	}

	err = migration(connectStr)
	if err != nil {
		return errors.New("migration error: " + err.Error())
	}

	pool, err := pgxpool.New(ctx, connectStr)
	if err != nil {
		return errors.New("database connection error: " + err.Error())
	}

	s.conn = conn
	s.pool = pool

	return nil
}

// Close Закрывает хранилище
func (s *Store) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.conn.Close(ctx)
	if err != nil {
		return err
	}
	return nil
}
