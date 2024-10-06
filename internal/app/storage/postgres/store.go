package postgres

import (
	"errors"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// Store объект хранилища для работы с postgresql
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
