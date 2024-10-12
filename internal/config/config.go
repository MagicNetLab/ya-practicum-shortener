package config

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config/env"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config/flags"
)

var servParams configParams

// GetParams возвращает параметры для запуска приложения
func GetParams() ParameterConfig {
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
		if envConf.HasBaseHost() {
			host, hostErr := envConf.GetBaseHost()
			port, portErr := envConf.GetBasePort()

			if hostErr == nil && portErr == nil {
				servParams.defaultHost = host
				servParams.defaultPort = port
			}
		}

		if envConf.HasShortHost() {
			host, hostErr := envConf.GetShortHost()
			port, portErr := envConf.GetShortPort()

			if hostErr == nil && portErr == nil {
				servParams.shortHost = host
				servParams.shortPort = port
			}
		}

		if envConf.HasFileStoragePath() {
			storagePath, storageErr := envConf.GetFileStoragePath()
			if storageErr == nil {
				servParams.fileStoragePath = storagePath
			}
		}

		if envConf.HasDBConnectString() {
			dbConnectParams, dbParamsErr := envConf.GetDBConnectString()
			if dbParamsErr == nil {
				servParams.dbConnectString = dbConnectParams
			}
		}

		if envConf.HasJWTSecret() {
			jwtSecret, jwtSecretErr := envConf.GetJWTSecret()
			if jwtSecretErr == nil {
				servParams.jwtSecret = jwtSecret
			}
		}

		servParams.pProfHost = envConf.GetPPROFHost()
	}

	cliConf := flags.Parse()

	if cliConf.HasDefaultHost() {
		host, hostErr := cliConf.GetDefaultHost()
		port, portErr := cliConf.GetDefaultPort()
		if hostErr == nil && portErr == nil {
			servParams.defaultHost = host
			servParams.defaultPort = port
		}
	}

	if cliConf.HasShortHost() {
		host, hostErr := cliConf.GetShortHost()
		port, portErr := cliConf.GetShortPort()
		if hostErr == nil && portErr == nil {
			servParams.shortHost = host
			servParams.shortPort = port
		}
	}

	if cliConf.HasFileStoragePath() {
		storagePath, storageErr := cliConf.GetFileStoragePath()
		if storageErr == nil {
			servParams.fileStoragePath = storagePath
		}
	}

	if cliConf.HasDBConnectString() {
		dbConnectParams, dbParamsErr := cliConf.GetDBConnectString()
		if dbParamsErr == nil {
			servParams.dbConnectString = dbConnectParams
		}
	}

	if cliConf.HasPProfHost() {
		servParams.pProfHost = cliConf.GetPProfHost()
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
