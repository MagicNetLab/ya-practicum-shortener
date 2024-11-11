package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	pgsql "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"

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

type store struct {
	params        map[string]string
	connectString string
}

// Init инициализация БД
func (s *store) Init() error {
	// todo отрефакторить и избавиться от постоянных подключений к БД (держать канал).
	conf := config.GetParams()

	connectParams, err := parseConnectString(conf.GetDBConnectString())
	if err == nil {
		s.params = connectParams
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		s.connectString = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			connectParams[dbUser],
			connectParams[dbPassword],
			connectParams[dbHost],
			connectParams[dbPort],
			connectParams[dbName],
			connectParams[dbDBSslMode])

		con, err := pgx.Connect(ctx, s.connectString)
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed to connect to database", args)
			return err
		}

		err = con.Ping(ctx)
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed to ping database", args)
			return err
		}
	} else {
		args := map[string]interface{}{"connect": conf.GetDBConnectString()}
		logger.Error("failed parse connect string", args)
		s.connectString = conf.GetDBConnectString()
	}

	if err = s.migration(); err != nil {
		return err
	}

	return nil
}

// PutLink сохранение ссылки пользователя в БД
func (s *store) PutLink(ctx context.Context, link string, short string, userID int) error {
	conn, err := pgx.Connect(ctx, s.connectString)
	if err != nil {
		return errors.New("database connection error: " + err.Error())
	}
	defer conn.Close(ctx)

	commandTag, err := conn.Exec(ctx, insertLinkSQL, short, link, userID)
	if err != nil {
		if strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
			return ErrLinkUniqueConflict
		}
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("filed save data: " + commandTag.String())
	}

	return nil
}

// PutBatchLinksArray пакетное сохранение ссылок пользователя в БД
func (s *store) PutBatchLinksArray(ctx context.Context, StoreBatchLicksArray map[string]string, userID int) error {
	conn, err := pgx.Connect(ctx, s.connectString)
	if err != nil {
		return errors.New("database connection error: " + err.Error())
	}
	defer conn.Close(ctx)

	transaction, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer transaction.Rollback(ctx)

	_, err = transaction.Prepare(ctx, "batch-insert", insertLinkSQL)
	if err != nil {
		return err
	}

	for key, value := range StoreBatchLicksArray {
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

// GetLink получение оригинальной ссылки из БД по сокращенному хэшу
func (s *store) GetLink(ctx context.Context, short string) (string, bool, error) {
	conn, err := pgx.Connect(ctx, s.connectString)
	if err != nil {
		return "", false, errors.New("database connection error: " + err.Error())
	}
	defer conn.Close(ctx)

	var link string
	var isDeleted bool
	err = conn.QueryRow(ctx, selectLinkSQL, short).Scan(&link, &isDeleted)
	if err != nil {
		return "", isDeleted, errors.New("database error: " + err.Error())
	}

	return link, isDeleted, nil
}

// HasShort проверка наличия короткой ссылки в БД
func (s *store) HasShort(ctx context.Context, short string) (bool, error) {
	conn, err := pgx.Connect(ctx, s.connectString)
	if err != nil {
		return false, errors.New("database connection error: " + err.Error())
	}
	defer conn.Close(ctx)

	var count int
	err = conn.QueryRow(ctx, hasLinkSQL, short).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

// GetShort получение сокращенного хэша для ссылки из БД
func (s *store) GetShort(ctx context.Context, link string) (string, error) {
	conn, err := pgx.Connect(ctx, s.connectString)
	if err != nil {
		return "", errors.New("database connection error: " + err.Error())
	}
	defer conn.Close(ctx)

	var short string
	err = conn.QueryRow(ctx, selectShortSQL, link).Scan(&short)
	if err != nil {
		return "", errors.New("database error: " + err.Error())
	}

	return short, nil

}

// GetUserLinks получение всех ссылок пользователя из БД
func (s *store) GetUserLinks(ctx context.Context, userID int) (map[string]string, error) {
	conn, err := pgx.Connect(ctx, s.connectString)
	if err != nil {
		return nil, errors.New("database connection error: " + err.Error())
	}
	defer conn.Close(ctx)

	res := make(map[string]string)

	rows, err := conn.Query(ctx, selectUserLinks, userID)
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

// DeleteBatchLinksArray пакетное удаление ссылок пользователя
func (s *store) DeleteBatchLinksArray(ctx context.Context, shorts []string, userID int) error {
	conn, err := pgx.Connect(ctx, s.connectString)
	if err != nil {
		return errors.New("database connection error: " + err.Error())
	}
	defer conn.Close(ctx)

	var paramrefs string
	var ids []interface{}
	ids = append(ids, userID)
	for i, v := range shorts {
		paramrefs += `$` + strconv.Itoa(i+2) + `,`
		ids = append(ids, v)
	}
	paramrefs = paramrefs[:len(paramrefs)-1]
	sqlQuery := `UPDATE links SET is_deleted = true WHERE user_id = $1 AND short IN (` + paramrefs + `)`
	exec, err := conn.Exec(ctx, sqlQuery, ids...)
	if err != nil {
		return err
	}

	if exec.RowsAffected() < 1 {
		return errors.New("filed delete data: " + exec.String())
	}

	return nil
}

func (s *store) Close() error {
	// todo сделать нормально закрытие коннектов к базе.
	return nil
}

func (s *store) migration() error {
	db, err := sql.Open("postgres", s.connectString)
	if err != nil {
		return errors.New("could not connect to postgres")
	}

	driver, err := pgsql.WithInstance(db, &pgsql.Config{})
	if err != nil {
		return errors.New("could not connect to postgres")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)

	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("could not init migrations", args)
	} else {
		m.Up()
	}

	return nil
}
