package config

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config/envreader"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config/flagsreader"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config/jsonreader"
)

// AppConfig объект с конфигурацией приложения
var AppConfig Configurator

// Initialize инициализация конфигурации приложения
func Initialize() error {
	if AppConfig.IsValid() {
		return nil
	}

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

	if !AppConfig.IsValid() {
		return errors.New("failed initialize application params")
	}

	return nil
}

// GetParams возвращает параметры для запуска приложения
func GetParams() *Configurator {
	return &AppConfig
}

// setDefaultParams установка дефолтных параметров на случай если из вне ни чего не пришло
func setDefaultParams() {
	AppConfig.defaultHost = "localhost"
	AppConfig.defaultPort = "8080"
	AppConfig.shortHost = "localhost"
	AppConfig.shortPort = "8080"
	AppConfig.fileStoragePath = "/tmp/short-url-db.json"
	AppConfig.jwtSecret = getRandomSecret()
}

func appendParams(reader ParamsReader) {
	defaultHost, err := reader.GetDefaultHost()
	if err == nil {
		AppConfig.defaultHost = defaultHost
	}

	defaultPort, err := reader.GetDefaultPort()
	if err == nil {
		AppConfig.defaultPort = defaultPort
	}

	shortHost, err := reader.GetShortHost()
	if err == nil {
		AppConfig.shortHost = shortHost
	}

	shortPort, err := reader.GetShortPort()
	if err == nil {
		AppConfig.shortPort = shortPort
	}

	storagePath, err := reader.GetFileStoragePath()
	if err == nil {
		AppConfig.fileStoragePath = storagePath
	}

	dbConnectParams, err := reader.GetDBConnectString()
	if err == nil {
		AppConfig.dbConnectString = dbConnectParams
	}

	jwtSecret, err := reader.GetJWTSecret()
	if err == nil {
		AppConfig.jwtSecret = jwtSecret
	}

	pprofHost, err := reader.GetPProfHost()
	if err == nil {
		AppConfig.pProfHost = pprofHost
	}

	if reader.HasEnableHTTPS() {
		AppConfig.enableHTTPS = reader.GetIsEnableHTTPS()
	}

	trustedSubnet, err := reader.GetTrustedSubnet()
	if err == nil {
		AppConfig.trustedSubnet = trustedSubnet
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
