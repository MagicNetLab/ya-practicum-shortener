package envreader

import "errors"

// Configurator параметры приложения собранные из переменных окружения
type Configurator struct {
	baseHost        []string `env:"SERVER_ADDRESS" envSeparator:":"`
	shortHost       []string `env:"BASE_URL" envSeparator:":"`
	fileStoragePath string   `env:"FILE_STORAGE_PATH"`
	dbConnectString string   `env:"DATABASE_DSN"`
	jwtSecret       string   `env:"JWT_SECRET"`
	pProfHost       string   `env:"PPROF_HOST" envDefault:"localhost:5000"`
	enableHTTPS     bool     `env:"ENABLE_HTTPS" envDefault:"false"`
	configFilePath  string   `env:"CONFIG"`
	trustedSubnet   string   `env:"TRUSTED_SUBNET"`
}

// GetDefaultHost возвращает базовый хост для запуска приложения
func (c Configurator) GetDefaultHost() (string, error) {
	if len(c.baseHost) == 0 || c.baseHost[0] == "" {
		return "", errors.New("base host not set")
	}
	return c.baseHost[0], nil
}

// GetDefaultPort возвращает базовый порт для запуска приложения
func (c Configurator) GetDefaultPort() (string, error) {
	if len(c.baseHost) == 0 || c.baseHost[1] == "" {
		return "", errors.New("base port not set")
	}
	return c.baseHost[1], nil
}

// GetShortHost возвращает хост для обработки переходов по коротким ссылкам
func (c Configurator) GetShortHost() (string, error) {
	if len(c.shortHost) == 0 || c.shortHost[0] == "" {
		return "", errors.New("short host not set")
	}
	return c.shortHost[0], nil
}

// GetShortPort возвращает порт для обработки переходов по коротким ссылкам
func (c Configurator) GetShortPort() (string, error) {
	if len(c.shortHost) == 0 || c.shortHost[1] == "" {
		return "", errors.New("short port not set")
	}
	return c.shortHost[1], nil
}

// GetFileStoragePath возвращает путь до файла локального хранилища ссылок
func (c Configurator) GetFileStoragePath() (string, error) {
	if c.fileStoragePath == "" {
		return "", errors.New("file storage path not set")
	}
	return c.fileStoragePath, nil
}

// GetDBConnectString возвращает строку с парамерами для подключения к БД
func (c Configurator) GetDBConnectString() (string, error) {
	if c.dbConnectString == "" {
		return "", errors.New("database connect string not set")
	}
	return c.dbConnectString, nil
}

// GetJWTSecret возвращает строку секрет для генерации  JWT токенов
func (c Configurator) GetJWTSecret() (string, error) {
	if c.jwtSecret == "" {
		return "", errors.New("jwt secret not set")
	}
	return c.jwtSecret, nil
}

// GetPProfHost возвращает хост для запуска профилировщика приложения
func (c Configurator) GetPProfHost() (string, error) {
	if c.pProfHost == "" {
		return "", errors.New("pprof host not set")
	}
	return c.pProfHost, nil
}

// GetIsEnableHTTPS возвращает флаг необходимости использования https для запуска сервера
func (c Configurator) GetIsEnableHTTPS() bool {
	return c.enableHTTPS
}

// HasEnableHTTPS возвращает был ли установлен параметр enableHTTPS
func (c Configurator) HasEnableHTTPS() bool {
	return true
}

// GetConfigFilePath возвращает путь до файла с json конфигурацией
func (c Configurator) GetConfigFilePath() (string, error) {
	if c.configFilePath == "" {
		return "", errors.New("config file path not set")
	}
	return c.configFilePath, nil
}

// GetTrustedSubnet возвращает адрес доверенной сети для доступа к статистике сервера
func (c Configurator) GetTrustedSubnet() (string, error) {
	if c.trustedSubnet == "" {
		return "", errors.New("trusted subnet not set")
	}
	return c.trustedSubnet, nil
}
