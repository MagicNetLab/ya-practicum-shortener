package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
	"github.com/jackc/pgx/v5"
)

// GetStore возвращает ссылку на объект хранилища
func GetStore() *Store {
	return &Store{}
}

// Ping проверка соединения с БД
func Ping() bool {
	var connectStr string
	conf := config.GetParams()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if conf.GetDBConnectString() == "" {
		logger.Error("Postgres connection params not configured", nil)
		return false
	}

	connectParams, err := parseConnectString(conf.GetDBConnectString())
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
		args := map[string]interface{}{"connect": conf.GetDBConnectString()}
		logger.Error("failed parse connect string", args)
		connectStr = conf.GetDBConnectString()
	}

	conn, err := pgx.Connect(ctx, connectStr)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("failed connect to database", args)
		return false
	}

	err = conn.Ping(ctx)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("database ping error", args)
		return false
	}

	return true
}
