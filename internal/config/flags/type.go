package flags

import "errors"

const (
	defaultHostKey  = "a"
	shortHostKey    = "b"
	fileStoragePath = "f"
	dbConnectKey    = "d"
	jwtSecret       = "j"
	pProfKey        = "p"
)

type cliConf struct {
	defaultHost     string
	defaultPort     string
	shortHost       string
	shortPort       string
	fileStoragePath string
	dbConnectString string
	jwtSecret       string
	pProfHost       string
}

// HasDefaultHost проверяет установлен базовый хост для запуска приложения или нет
func (cc *cliConf) HasDefaultHost() bool {
	return cc.defaultHost != "" && cc.defaultPort != ""
}

// GetDefaultHost возвращает базовый хост для запуска приложения
func (cc *cliConf) GetDefaultHost() (string, error) {
	if cc.defaultHost == "" {
		return "", errors.New("default host not set")
	}

	return cc.defaultHost, nil
}

// GetDefaultPort возвращает базовый порт для запуска приложения
func (cc *cliConf) GetDefaultPort() (string, error) {
	if cc.defaultPort == "" {
		return "", errors.New("default port not set")
	}

	return cc.defaultPort, nil
}

// HasShortHost проверяет установлен ли хост для перенаправлений при переходе по коротким ссылкам
func (cc *cliConf) HasShortHost() bool {
	return cc.shortHost != "" && cc.shortPort != ""
}

// GetShortHost возвращает хост для обработки переходов по коротким ссылкам
func (cc *cliConf) GetShortHost() (string, error) {
	if cc.shortHost == "" {
		return "", errors.New("short host not set")
	}

	return cc.shortHost, nil
}

// GetShortPort возвращает порт для обработки переходов по коротким ссылкам
func (cc *cliConf) GetShortPort() (string, error) {
	if cc.shortPort == "" {
		return "", errors.New("short port not set")
	}

	return cc.shortPort, nil
}

// HasFileStoragePath проверяет установлен ли путь до файла для локального хранения кэша
func (cc *cliConf) HasFileStoragePath() bool {
	return cc.fileStoragePath != ""
}

// GetFileStoragePath возвращает пусть до файла для локального хранения кэша
func (cc *cliConf) GetFileStoragePath() (string, error) {
	if !cc.HasFileStoragePath() {
		return "", errors.New("file storage path not set")
	}

	return cc.fileStoragePath, nil
}

// HasDBConnectString проверяет установлена ли строка с настройками для подключения к БД
func (cc *cliConf) HasDBConnectString() bool {
	return cc.dbConnectString != ""
}

// GetDBConnectString возвращает строку с настройками для подключения в БД
func (cc *cliConf) GetDBConnectString() (string, error) {
	if !cc.HasDBConnectString() {
		return "", errors.New("database connect param not set")
	}

	return cc.dbConnectString, nil
}

// HasJWTSecret проверяет установлен ли секрет для генерации JWT токенов
func (cc *cliConf) HasJWTSecret() bool {
	return cc.jwtSecret != ""
}

// GetJWTSecret возвращает секрет для генерации JWT токенов
func (cc *cliConf) GetJWTSecret() (string, error) {
	if cc.jwtSecret == "" {
		return "", errors.New("jwttoken secret not set")
	}
	return cc.jwtSecret, nil
}

// HasPProfHost проверяет установлен ли хост для запуска профилировщика
func (cc *cliConf) HasPProfHost() bool {
	return cc.pProfHost != ""
}

// GetPProfHost возвращает хост для запуска профилировщика
func (cc *cliConf) GetPProfHost() string {
	return cc.pProfHost
}
