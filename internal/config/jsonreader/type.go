package jsonreader

import (
	"errors"
	"strings"
)

type Configurator struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DataBaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
}

// GetDefaultHost возвращает базовый хост для запуска приложения
func (c Configurator) GetDefaultHost() (string, error) {
	if !strings.HasSuffix(c.ServerAddress, ":") {
		return "", errors.New("default host not set")
	}

	if defaultHost := strings.Split(c.ServerAddress, ":")[0]; defaultHost != "" {
		return defaultHost, nil
	}

	return "", errors.New("default host not set")
}

// GetDefaultPort возвращает базовый порт для запуска приложения
func (c Configurator) GetDefaultPort() (string, error) {
	if !strings.HasSuffix(c.ServerAddress, ":") {
		return "", errors.New("default port not set")
	}

	if defaultPort := strings.Split(c.ServerAddress, ":")[1]; defaultPort != "" {
		return defaultPort, nil
	}

	return "", errors.New("default port not set")
}

// GetShortHost возвращает хост для обработки переходов по коротким ссылкам
func (c Configurator) GetShortHost() (string, error) {
	if !strings.HasSuffix(c.BaseURL, ":") {
		return "", errors.New("short host not set")
	}

	if shortHost := strings.Split(c.BaseURL, ":")[0]; shortHost != "" {
		return shortHost, nil
	}

	return "", errors.New("short host not set")
}

// GetShortPort возвращает порт для обработки переходов по коротким ссылкам
func (c Configurator) GetShortPort() (string, error) {
	if !strings.HasSuffix(c.BaseURL, ":") {
		return "", errors.New("short port not set")
	}

	if shortPort := strings.Split(c.BaseURL, ":")[1]; shortPort != "" {
		return shortPort, nil
	}

	return "", errors.New("short port not set")
}

// GetFileStoragePath возвращает путь до файла локального хранилища ссылок
func (c Configurator) GetFileStoragePath() (string, error) {
	if c.FileStoragePath == "" {
		return "", errors.New("file storage path not set")
	}
	return c.FileStoragePath, nil
}

// GetDBConnectString возвращает строку с парамерами для подключения к БД
func (c Configurator) GetDBConnectString() (string, error) {
	if c.DataBaseDSN == "" {
		return "", errors.New("no db connect string specified")
	}
	return c.DataBaseDSN, nil
}

// GetJWTSecret возвращает строку секрет для генерации  JWT токенов
func (c Configurator) GetJWTSecret() (string, error) {
	return "", errors.New("no jwt secret specified")
}

// GetPProfHost возвращает хост для запуска профилировщика приложения
func (c Configurator) GetPProfHost() (string, error) {
	return "", errors.New("no pprof host specified")
}

// GetIsEnableHTTPS возвращает флаг необходимости использования https для запуска сервера
func (c Configurator) GetIsEnableHTTPS() bool {
	return c.EnableHTTPS
}

// HasEnableHTTPS возвращает был ли установлен параметр enableHTTPS
func (c Configurator) HasEnableHTTPS() bool {
	return true
}

func (c Configurator) GetConfigFilePath() (string, error) {
	return "", errors.New("no config file path specified")
}
