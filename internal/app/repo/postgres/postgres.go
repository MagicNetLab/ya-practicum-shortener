package postgres

import (
	"database/sql"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

// GetStore возвращает ссылку на объект хранилища
func GetStore() *Store {
	return &Store{}
}

// Ping проверка соединения с БД
func Ping() bool {
	parameterConfig := config.GetParams()
	dbConnect := parameterConfig.GetDBConnectString()

	if dbConnect == "" {
		logger.Error("Postgres connection params not configured", nil)
		return false
	}

	db, err := sql.Open("pgx", dbConnect)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("Postgres connection error: %v", args)
		return false
	}
	defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		args := map[string]interface{}{"error": pingErr.Error()}
		logger.Error("Postgres ping error: %v", args)
		return false
	}

	return true
}
