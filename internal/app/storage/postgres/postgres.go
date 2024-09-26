package postgres

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

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

	pingErr := db.Ping()
	if pingErr != nil {
		args := map[string]interface{}{"error": pingErr.Error()}
		logger.Error("Postgres ping error: %v", args)
		return false
	}

	defer db.Close()

	return true
}
