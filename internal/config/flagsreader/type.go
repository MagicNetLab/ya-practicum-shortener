package flagsreader

import "errors"

const (
	defaultHostKey  = "a"
	shortHostKey    = "b"
	fileStoragePath = "f"
	dbConnectKey    = "d"
	jwtSecret       = "j"
	pProfKey        = "p"
	enableHTTPSKey  = "s"
	configFileKey   = "c"
	trustedSubnet   = "t"
)

// CliConf параметры приложения собранные из флагов указанных при запуске
type CliConf struct {
	defaultHost     string
	defaultPort     string
	shortHost       string
	shortPort       string
	fileStoragePath string
	dbConnectString string
	jwtSecret       string
	pProfHost       string
	enableHTTPS     bool
	hasEnableHTTPS  bool
	configFilePath  string
	trustedSubnet   string
}

// GetDefaultHost возвращает базовый хост для запуска приложения
func (c CliConf) GetDefaultHost() (string, error) {
	if c.defaultHost == "" {
		return "", errors.New("default host not set")
	}
	return c.defaultHost, nil
}

// GetDefaultPort возвращает базовый порт для запуска приложения
func (c CliConf) GetDefaultPort() (string, error) {
	if c.defaultPort == "" {
		return "", errors.New("default port not set")
	}
	return c.defaultPort, nil
}

// GetShortHost возвращает хост для обработки переходов по коротким ссылкам
func (c CliConf) GetShortHost() (string, error) {
	if c.shortHost == "" {
		return "", errors.New("short host not set")
	}
	return c.shortHost, nil
}

// GetShortPort возвращает порт для обработки переходов по коротким ссылкам
func (c CliConf) GetShortPort() (string, error) {
	if c.shortPort == "" {
		return "", errors.New("short port not set")
	}
	return c.shortPort, nil
}

// GetFileStoragePath возвращает путь до файла локального хранилища ссылок
func (c CliConf) GetFileStoragePath() (string, error) {
	if c.fileStoragePath == "" {
		return "", errors.New("file storage path not set")
	}
	return c.fileStoragePath, nil
}

// GetDBConnectString возвращает строку с парамерами для подключения к БД
func (c CliConf) GetDBConnectString() (string, error) {
	if c.dbConnectString == "" {
		return "", errors.New("db connect string not set")
	}
	return c.dbConnectString, nil
}

// GetJWTSecret возвращает строку секрет для генерации  JWT токенов
func (c CliConf) GetJWTSecret() (string, error) {
	if c.jwtSecret == "" {
		return "", errors.New("jwt secret not set")
	}
	return c.jwtSecret, nil
}

// GetPProfHost возвращает хост для запуска профилировщика приложения
func (c CliConf) GetPProfHost() (string, error) {
	if c.pProfHost == "" {
		return "", errors.New("pprof host not set")
	}
	return c.pProfHost, nil
}

// GetIsEnableHTTPS возвращает флаг необходимости использования https для запуска сервера
func (c CliConf) GetIsEnableHTTPS() bool {
	return c.enableHTTPS
}

// HasEnableHTTPS возвращает был ли задействован флаг использования https или нет
func (c CliConf) HasEnableHTTPS() bool {
	return c.hasEnableHTTPS
}

// GetConfigFilePath возвращает путь до файла с json конфигурацией
func (c CliConf) GetConfigFilePath() (string, error) {
	if c.configFilePath == "" {
		return "", errors.New("config file path not set")
	}

	return c.configFilePath, nil
}

// GetTrustedSubnet Возвращает адрес доверенной сети для доступа к статистике сервера
func (c CliConf) GetTrustedSubnet() (string, error) {
	if c.trustedSubnet == "" {
		return "", errors.New("trusted subnet not set")
	}
	return c.trustedSubnet, nil
}
