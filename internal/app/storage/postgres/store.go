package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"strings"
	"time"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
	"github.com/golang-migrate/migrate/v4"
	pgsql "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

const (
	insertLinkSQL   = "INSERT INTO links (short, link, user_id) VALUES ($1, $2, $3)"
	selectLinkSQL   = "SELECT link FROM links WHERE short = $1"
	selectShortSQL  = "SELECT short FROM links WHERE link = $1"
	hasLinkSQL      = "SELECT count(*) FROM links WHERE short = $1"
	selectUserLinks = "SELECT short, link FROM links WHERE user_id = $1"
)

type store struct {
	params        map[string]string
	connectString string
}

var ErrLinkUniqueConflict = errors.New("url is not unique")

func (s *store) Init() error {
	conf := config.GetParams()

	connectParams, err := parseConnectString(conf.GetDBConnectString())
	if err == nil {
		s.params = connectParams
		s.connectString = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			connectParams[DBUser],
			connectParams[DBPassword],
			connectParams[DBHost],
			connectParams[DBPort],
			connectParams[DBName],
			connectParams[DBSslMode])

		con, err := pgx.Connect(context.Background(), s.connectString)
		if err != nil {
			logger.Log.Infof("Connect string is incorrect %s", err)
			return err
		}

		err = con.Ping(context.Background())
		if err != nil {
			logger.Log.Infof("Unable to connect to database %s", err)
			return err
		}
	} else {
		logger.Log.Infof("Faled parse connect string: %s", conf.GetDBConnectString())
		s.connectString = conf.GetDBConnectString()
	}

	if err = s.migration(); err != nil {
		return err
	}

	return nil
}

func (s *store) PutLink(link string, short string, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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

func (s *store) PutBatchLinksArray(StoreBatchLicksArray map[string]string, userID int) error {
	// TODO use prepare statement
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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
			return err
		}

		if cTag.RowsAffected() != 1 {
			return errors.New("filed save data: " + cTag.String())
		}
	}

	transaction.Commit(ctx)

	return nil
}

func (s *store) GetLink(short string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, s.connectString)
	if err != nil {
		return "", errors.New("database connection error: " + err.Error())
	}
	defer conn.Close(ctx)

	var link string
	err = conn.QueryRow(ctx, selectLinkSQL, short).Scan(&link)
	if err != nil {
		return "", errors.New("database error: " + err.Error())
	}

	return link, nil
}

func (s *store) HasShort(short string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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

func (s *store) GetShort(link string) (string, error) {
	// todo use context with timeout from handlers
	ctx := context.Background()
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

func (s *store) GetUserLinks(userID int) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
		logger.Log.Errorf("could not init migrations: %s", err)
	} else {
		m.Up()
	}

	return nil
}

var Store = store{}

func parseConnectString(connectString string) (map[string]string, error) {
	params := make(map[string]string)
	if connectString == "" {
		return params, errors.New("connect params is empty")
	}

	p := strings.Split(connectString, " ")
	for _, v := range p {
		val := strings.Split(v, "=")
		switch string(val[0]) {
		case "host":
			params[DBHost] = strings.TrimSpace(val[1])
		case "user":
			params[DBUser] = strings.TrimSpace(val[1])
		case "password":
			params[DBPassword] = strings.TrimSpace(val[1])
		case "dbname":
			params[DBName] = strings.TrimSpace(val[1])
		case "sslmode":
			params[DBSslMode] = strings.TrimSpace(val[1])
		}
	}

	if len(params) == 0 {
		return params, errors.New("connect params not set")
	}

	params[DBPort] = "5432"

	return params, nil
}
