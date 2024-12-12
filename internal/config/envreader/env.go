package envreader

import (
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

// Parse парсинг env параметров. Возвращает IConfigurator и ошибку если что-то пошло не так при сборе параметров.
func Parse() Configurator {
	var envConf Configurator

	err := godotenv.Load(".env")
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error(".env file not found", args)
	}

	baseHost := os.Getenv("SERVER_ADDRESS")
	// ожидаемый формат host:port но в тестах на github данные могут приходить в формате  http://host или http://host:port
	// приводим значение к ожидаемому и если не получилось то будет вставлено дефолтное значение
	baseHost = strings.Replace(baseHost, "http://", "", 1)
	if baseHost != "" && strings.Contains(baseHost, ":") {
		envConf.baseHost = strings.Split(baseHost, ":")
	} else {
		args := map[string]interface{}{"server_address_value": os.Getenv("SERVER_ADDRESS")}
		logger.Error("invalid value SERVER_ADDRESS from env. default value set", args)
	}

	shortHost := os.Getenv("BASE_URL")
	// ожидаемый формат host:port но в тестах на github данные могут приходить в формате  http://host или http://host:port
	// приводим значение к ожидаемому и если не получилось то будет вставлено дефолтное значение
	shortHost = strings.Replace(shortHost, "http://", "", 1)
	if shortHost != "" && strings.Contains(shortHost, ":") {
		envConf.shortHost = strings.Split(shortHost, ":")
	} else {
		args := map[string]interface{}{"base_url_value": os.Getenv("BASE_URL")}
		logger.Error("invalid value BASE_URL from env. default value set", args)
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

	if configFilePath := os.Getenv("CONFIG"); configFilePath != "" {
		envConf.configFilePath = configFilePath
	}

	if trustedSubnet := os.Getenv("TRUSTED_SUBNET"); trustedSubnet != "" {
		envConf.trustedSubnet = trustedSubnet
	}

	return envConf
}
