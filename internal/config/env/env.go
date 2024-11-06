package env

import (
	"log"
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
	log.Println(baseHost)
	if baseHost != "" && strings.Contains(baseHost, ":") {
		envConf.baseHost = strings.Split(baseHost, ":")
	}

	shortHost := os.Getenv("BASE_URL")
	log.Println(shortHost)
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
