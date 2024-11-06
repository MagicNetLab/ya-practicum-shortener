package env

import (
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

// Parse парсинг env параметров. Возвращает IConfigurator и ошибку если что-то пошло не так при сборе параметров.
func Parse() (Configurator, error) {
	var envConf Configurator

	err := godotenv.Load(".env")
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error(".env file not found", args)
	}

	baseHost := os.Getenv("SERVER_ADDRESS")
	// костыль для тестов на github. в env может приходить всякая ересь
	if strings.Contains(baseHost, "http://") {
		baseHost = strings.Replace(baseHost, "http://", "", 1)
	}
	if baseHost != "" && strings.Contains(baseHost, ":") {
		envConf.baseHost = strings.Split(baseHost, ":")
	}

	shortHost := os.Getenv("BASE_URL")
	// костыль для тестов на github. в env может приходить всякая ересь
	if strings.Contains(shortHost, "http://") {
		shortHost = strings.Replace(shortHost, "http://", "", 1)
	}
	if shortHost != "" && strings.Contains(shortHost, ":") {
		envConf.shortHost = strings.Split(shortHost, ":")
	}

	if fileStorage := os.Getenv("FILE_STORAGE_PATH"); fileStorage != "" {
		envConf.fileStoragePath = fileStorage
	}

	if dbParams := os.Getenv("DATABASE_DSN"); dbParams != "" {
		envConf.dbConnectString = dbParams
	}

	if JWTSecret := os.Getenv("JWT_SECRET"); JWTSecret != "" {
		envConf.jwtSecret = JWTSecret
	}

	if pProfHost := os.Getenv("PPROF_HOST"); pProfHost != "" {
		envConf.pProfHost = pProfHost
	}

	if enableHTTPS := os.Getenv("ENABLE_HTTPS"); enableHTTPS != "" {
		envConf.enableHTTPS = enableHTTPS == "true"
	}

	return envConf, nil
}
