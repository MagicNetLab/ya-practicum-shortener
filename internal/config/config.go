package config

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config/envreader"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config/flagsreader"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config/jsonreader"
)

var servParams configParams

func Initialize() error {
	var confFile string
	setDefaultParams()
	envConf := envreader.Parse()
	cliConf := flagsreader.Parse()

	if filePath, err := envConf.GetConfigFilePath(); err == nil && filePath != "" {
		confFile = filePath
	}
	if filePath, err := cliConf.GetConfigFilePath(); err == nil && filePath != "" {
		confFile = filePath
	}
	if confFile != "" {
		jsonConf := jsonreader.Parse(confFile)
		appendParams(jsonConf)
	}

	appendParams(envConf)
	appendParams(cliConf)

	if !servParams.IsValid() {
		return errors.New("failed initialize application params")
	}

	return nil
}

// GetParams возвращает параметры для запуска приложения
func GetParams() AppConfig {
	return &servParams
}

func setDefaultParams() {
	// костыль для тестов. без этих значений приложение в принципе не должно запускаться
	servParams.defaultHost = "localhost"
	servParams.defaultPort = "8080"
	servParams.shortHost = "localhost"
	servParams.shortPort = "8080"
	servParams.fileStoragePath = "/tmp/short-url-db.json"
	servParams.jwtSecret = getRandomSecret()
}

func appendParams(reader ParamsReader) {
	defaultHost, err := reader.GetDefaultHost()
	if err == nil {
		servParams.defaultHost = defaultHost
	}

	defaultPort, err := reader.GetDefaultPort()
	if err == nil {
		servParams.defaultPort = defaultPort
	}

	shortHost, err := reader.GetShortHost()
	if err == nil {
		servParams.shortHost = shortHost
	}

	shortPort, err := reader.GetShortPort()
	if err == nil {
		servParams.shortPort = shortPort
	}

	storagePath, err := reader.GetFileStoragePath()
	if err == nil {
		servParams.fileStoragePath = storagePath
	}

	dbConnectParams, err := reader.GetDBConnectString()
	if err == nil {
		servParams.dbConnectString = dbConnectParams
	}

	jwtSecret, err := reader.GetJWTSecret()
	if err == nil {
		servParams.jwtSecret = jwtSecret
	}

	pprofHost, err := reader.GetPProfHost()
	if err == nil {
		servParams.pProfHost = pprofHost
	}

	if reader.HasEnableHTTPS() {
		servParams.enableHTTPS = reader.GetIsEnableHTTPS()
	}
}

func getRandomSecret() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}
