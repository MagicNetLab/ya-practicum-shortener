package config

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config/env"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config/flags"
)

var servParams configParams

// GetParams возвращает параметры для запуска приложения
func GetParams() AppConfig {
	if servParams.IsValid() {
		return &servParams
	}

	// костыль для тестов. без этих значений приложение в принципе не должно запускаться
	servParams.defaultHost = "localhost"
	servParams.defaultPort = "8080"
	servParams.shortHost = "localhost"
	servParams.shortPort = "8080"
	servParams.fileStoragePath = "/tmp/short-url-db.json"
	servParams.jwtSecret = getRandomSecret()

	envConf, err := env.Parse()
	if err == nil {
		defaultHost, err := envConf.GetDefaultHost()
		if err == nil {
			servParams.defaultHost = defaultHost
		}

		defaultPort, err := envConf.GetDefaultPort()
		if err == nil {
			servParams.defaultPort = defaultPort
		}

		shortHost, err := envConf.GetShortHost()
		if err == nil {
			servParams.shortHost = shortHost
		}

		shortPort, err := envConf.GetShortPort()
		if err == nil {
			servParams.shortPort = shortPort
		}

		storagePath, err := envConf.GetFileStoragePath()
		if err == nil {
			servParams.fileStoragePath = storagePath
		}

		dbConnectParams, err := envConf.GetDBConnectString()
		if err == nil {
			servParams.dbConnectString = dbConnectParams
		}

		jwtSecret, err := envConf.GetJWTSecret()
		if err == nil {
			servParams.jwtSecret = jwtSecret
		}

		pprofHost, err := envConf.GetPProfHost()
		if err == nil {
			servParams.pProfHost = pprofHost
		}

		servParams.enableHTTPS = envConf.GetIsEnableHTTPS()
	}

	cliConf := flags.Parse()
	defaultHost, err := cliConf.GetDefaultHost()
	if err == nil {
		servParams.defaultHost = defaultHost
	}

	defaultPort, err := cliConf.GetDefaultPort()
	if err == nil {
		servParams.defaultPort = defaultPort
	}

	shortHost, err := cliConf.GetShortHost()
	if err == nil {
		servParams.shortHost = shortHost
	}

	shortPort, err := cliConf.GetShortPort()
	if err == nil {
		servParams.shortPort = shortPort
	}

	storagePath, err := cliConf.GetFileStoragePath()
	if err == nil {
		servParams.fileStoragePath = storagePath
	}

	dbConnectParams, err := cliConf.GetDBConnectString()
	if err == nil {
		servParams.dbConnectString = dbConnectParams
	}

	pprofHost, err := cliConf.GetPProfHost()
	if err == nil {
		servParams.pProfHost = pprofHost
	}

	if cliConf.HasEnableHTTPS() {
		servParams.enableHTTPS = cliConf.GetIsEnableHTTPS()
	}

	return &servParams
}

func getRandomSecret() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}
