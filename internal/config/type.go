package config

// AppConfig хранилище параметров для запуска и работы приложения
type AppConfig interface {
	// GetDefaultHost возвращает хост для запуска приложения
	GetDefaultHost() string
	// GetShortHost возвращает хост для переходов по коротким ссылкам
	GetShortHost() string
	// GetFileStoragePath возвращает пусть до файла локального хранилища ссылок
	GetFileStoragePath() string
	// GetDBConnectString Возвращает строку с настройками для подключения к БД
	GetDBConnectString() string
	// GetJWTSecret возвращает секрет для генерации JWT токенов
	GetJWTSecret() string
	// GetPProfHost возвращает хост для запуска профилировщика
	GetPProfHost() string
	// IsEnableHTTPS возвращает нужно ли использовать https при запуске сервера
	IsEnableHTTPS() bool
	// IsValid проверяет корректность настроек для работы проложения
	IsValid() bool
}

// ParamsReader интерфейс для плагинов чтения конфигураций
type ParamsReader interface {
	// GetDefaultHost возвращает базовый хост для запуска приложения
	GetDefaultHost() (string, error)
	// GetDefaultPort возвращает базовый порт для запуска приложения
	GetDefaultPort() (string, error)
	// GetShortHost возвращает хост для обработки переходов по коротким ссылкам
	GetShortHost() (string, error)
	// GetShortPort возвращает порт для обработки переходов по коротким ссылкам
	GetShortPort() (string, error)
	// GetFileStoragePath возвращает путь до файла локального хранилища ссылок
	GetFileStoragePath() (string, error)
	// GetDBConnectString возвращает строку с парамерами для подключения к БД
	GetDBConnectString() (string, error)
	// GetJWTSecret возвращает строку секрет для генерации  JWT токенов
	GetJWTSecret() (string, error)
	// GetPProfHost возвращает хост для запуска профилировщика приложения
	GetPProfHost() (string, error)
	// GetIsEnableHTTPS возвращает флаг необходимости использования https для запуска сервера
	GetIsEnableHTTPS() bool
	// HasEnableHTTPS возвращает был ли установлен параметр enableHTTPS
	HasEnableHTTPS() bool
	// GetConfigFilePath возвращает имя файла с json конфигурацией
	GetConfigFilePath() (string, error)
}

type configParams struct {
	defaultHost     string
	defaultPort     string
	shortHost       string
	shortPort       string
	fileStoragePath string
	dbConnectString string
	jwtSecret       string
	pProfHost       string
	enableHTTPS     bool
}

// GetDefaultHost возвращает хост для запуска приложения
func (c configParams) GetDefaultHost() string {
	return c.defaultHost + ":" + c.defaultPort
}

// GetShortHost возвращает хост для переходов по коротким ссылкам
func (c configParams) GetShortHost() string {
	return c.shortHost + ":" + c.shortPort
}

// GetFileStoragePath возвращает пусть до файла локального хранилища ссылок
func (c configParams) GetFileStoragePath() string {
	return c.fileStoragePath
}

// GetDBConnectString Возвращает строку с настройками для подключения к БД
func (c configParams) GetDBConnectString() string {
	return c.dbConnectString
}

// GetJWTSecret возвращает секрет для генерации JWT токенов
func (c configParams) GetJWTSecret() string {
	return c.jwtSecret
}

// GetPProfHost возвращает хост для запуска профилировщика
func (c configParams) GetPProfHost() string {
	return c.pProfHost
}

// IsEnableHTTPS возвращает нужно ли использовать https при запуске сервера
func (c configParams) IsEnableHTTPS() bool {
	return c.enableHTTPS
}

// IsValid проверяет корректность настроек для работы проложения
func (c configParams) IsValid() bool {
	return c.defaultHost != "" && c.defaultPort != "" && c.shortHost != "" && c.shortPort != "" && (c.fileStoragePath != "" || c.dbConnectString != "") && c.jwtSecret != ""
}
