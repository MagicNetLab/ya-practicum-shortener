package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
	"github.com/golang-migrate/migrate/v4"
	pgsql "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

type store struct {
	params        map[string]string
	connectString string
}

func (s *store) Init() error {
	conf := config.GetParams()

	connectParams, err := parseConnectString(conf.GetDBConnectString())
	if err != nil {
		logger.Log.Infof("Faled parse connect string: %s", conf.GetDBConnectString())
		return err
	}

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
	}

	if err = s.migration(); err != nil {
		return err
	}

	return nil
}

func (s *store) PutLink(link string, short string) error {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, s.connectString)
	if err != nil {
		return errors.New("database connection error: " + err.Error())
	}
	defer conn.Close(ctx)

	commandTag, err := conn.Exec(ctx, "INSERT INTO links (short, link) VALUES ($1, $2)", short, link)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("filed save data: " + commandTag.String())
	}

	return nil
}

func (s *store) GetLink(short string) (string, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, s.connectString)
	if err != nil {
		return "", errors.New("database connection error: " + err.Error())
	}
	defer conn.Close(ctx)

	var link string
	err = conn.QueryRow(ctx, "SELECT link FROM links WHERE short = $1", short).Scan(&link)
	if err != nil {
		return "", errors.New("database error: " + err.Error())
	}

	return link, nil
}

func (s *store) HasShort(short string) (bool, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, s.connectString)
	if err != nil {
		return false, errors.New("database connection error: " + err.Error())
	}
	defer conn.Close(ctx)

	var count int
	err = conn.QueryRow(ctx, "SELECT count(*) FROM links WHERE short = $1", short).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
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
		"file://../../migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return errors.New("could not init migrations")
	}

	m.Up()

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
