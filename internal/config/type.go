package config

import "strings"

type configParams struct {
	defaultHost     string
	defaultPort     string
	shortHost       string
	shortPort       string
	fileStoragePath string
	dbConnectString string
	jwtSecret       string
	pProfOn         bool
	pProfHost       string
}

// SetFileStoragePath установка пути до файла для локального хранения кэша
func (sp *configParams) SetFileStoragePath(path string) error {
	sp.fileStoragePath = path
	return nil
}

// GetFileStoragePath получение пути до файла для локального хранения кэша
func (sp *configParams) GetFileStoragePath() string {
	return sp.fileStoragePath
}

// IsValid проверка корректности настроек для запуска приложения
func (sp *configParams) IsValid() bool {
	return sp.defaultHost != "" &&
		sp.defaultPort != "" &&
		sp.shortHost != "" &&
		sp.shortPort != "" &&
		sp.fileStoragePath != ""

}

// SetDefaultHost установка хоста для работы с данными пользователя
func (sp *configParams) SetDefaultHost(host string, port string) error {
	sp.defaultHost = host
	sp.defaultPort = port

	return nil
}

// SetShortHost установка хоста для обработки переходов по  коротким ссылкам
func (sp *configParams) SetShortHost(host string, port string) error {
	sp.shortHost = host
	sp.shortPort = port

	return nil
}

// GetDefaultHost получение хоста для работы с данными пользователя
func (sp *configParams) GetDefaultHost() string {
	p := []string{sp.defaultHost, sp.defaultPort}
	return strings.Join(p, ":")
}

// GetShortHost получение хоста для обработки переходов по коротким ссылкам
func (sp *configParams) GetShortHost() string {
	p := []string{sp.shortHost, sp.shortPort}
	return strings.Join(p, ":")
}

// SetDBConnectString установка строки с параметрами для подключение к БД
func (sp *configParams) SetDBConnectString(params string) error {
	sp.dbConnectString = params

	return nil
}

// GetDBConnectString получение строки с параметрами для подключения к БД
func (sp *configParams) GetDBConnectString() string {
	return sp.dbConnectString
}

// SetJWTSecret установка секрета для генерации jwt токенов пользователей
func (sp *configParams) SetJWTSecret(secret string) error {
	sp.jwtSecret = secret
	return nil
}

// GetJWTSecret получение секрета для генерации jwt токенов пользователей
func (sp *configParams) GetJWTSecret() string {
	return sp.jwtSecret
}

// SetPProfOn установка параметра активности профилировщика
func (sp *configParams) SetPProfOn(isOn bool) error {
	sp.pProfOn = isOn

	return nil
}

// IsPProfOn проверка активен профилировщик или нет
func (sp *configParams) IsPProfOn() bool {
	return sp.pProfOn
}

// SetPProfHost установка хоста для запуска профилировщика
func (sp *configParams) SetPProfHost(host string) error {
	sp.pProfHost = host

	return nil
}

// GetPProfHost получение хоста для запуска профилировщика
func (sp *configParams) GetPProfHost() string {
	return sp.pProfHost
}
