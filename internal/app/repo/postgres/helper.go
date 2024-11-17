package postgres

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
	"github.com/golang-migrate/migrate/v4"
	pgsql "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// parseConnectString парсинг строки с параметрами соединения с БД.
// ожидаемый формат разделенные пробелами key=value
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
			params[dbHost] = strings.TrimSpace(val[1])
		case "user":
			params[dbUser] = strings.TrimSpace(val[1])
		case "password":
			params[dbPassword] = strings.TrimSpace(val[1])
		case "dbname":
			params[dbName] = strings.TrimSpace(val[1])
		case "sslmode":
			params[dbDBSslMode] = strings.TrimSpace(val[1])
		}
	}

	if len(params) == 0 {
		return params, errors.New("connect params not set")
	}

	params[dbPort] = "5432"

	return params, nil
}

// migration применение миграций к БД
func migration(connectString string) error {
	db, err := sql.Open("postgres", connectString)
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
