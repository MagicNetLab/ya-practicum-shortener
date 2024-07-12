package postgres

import (
	"database/sql"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Ping() bool {
	parameterConfig := config.GetParams()
	dbConnect := parameterConfig.GetDBConnectString()

	if dbConnect == "" {
		logger.Log.Errorf("Postgres connection params not configured")
		return false
	}

	db, err := sql.Open("pgx", dbConnect)
	if err != nil {
		logger.Log.Infof("Postgres connection error: %v", err)
		return false
	}

	pingErr := db.Ping()
	if pingErr != nil {
		logger.Log.Infof("Postgres ping error: %v", pingErr)
		return false
	}

	defer db.Close()

	return true
}
