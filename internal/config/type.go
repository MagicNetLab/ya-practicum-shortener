package config

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
	// GetTrustedSubnet возвращает вдрес доверенной сети для просмотра статистики сервера
	GetTrustedSubnet() (string, error)
	// GetGRPCPort возвращает номер порта для запуска grpc сервера
	GetGRPCPort() (string, error)
}

// Configurator хранилище параметров для запуска и работы приложения
type Configurator struct {
	defaultHost     string
	defaultPort     string
	shortHost       string
	shortPort       string
	fileStoragePath string
	dbConnectString string
	jwtSecret       string
	pProfHost       string
	enableHTTPS     bool
	trustedSubnet   string
	grpcPort        string
}

// GetDefaultHost возвращает хост для запуска приложения
func (c Configurator) GetDefaultHost() string {
	return c.defaultHost + ":" + c.defaultPort
}

// GetShortHost возвращает хост для переходов по коротким ссылкам
func (c Configurator) GetShortHost() string {
	return c.shortHost + ":" + c.shortPort
}

// GetFileStoragePath возвращает пусть до файла локального хранилища ссылок
func (c Configurator) GetFileStoragePath() string {
	return c.fileStoragePath
}

// GetDBConnectString Возвращает строку с настройками для подключения к БД
func (c Configurator) GetDBConnectString() string {
	return c.dbConnectString
}

// GetJWTSecret возвращает секрет для генерации JWT токенов
func (c Configurator) GetJWTSecret() string {
	return c.jwtSecret
}

// GetPProfHost возвращает хост для запуска профилировщика
func (c Configurator) GetPProfHost() string {
	return c.pProfHost
}

// IsEnableHTTPS возвращает нужно ли использовать https при запуске сервера
func (c Configurator) IsEnableHTTPS() bool {
	return c.enableHTTPS
}

// IsValid проверяет корректность настроек для работы проложения
func (c Configurator) IsValid() bool {
	return c.defaultHost != "" && c.defaultPort != "" && c.shortHost != "" && c.shortPort != "" && (c.fileStoragePath != "" || c.dbConnectString != "") && c.jwtSecret != ""
}

// GetTrustedSubnet возвращает адрес доверенной сети для доступа к статистике сервера
func (c Configurator) GetTrustedSubnet() string {
	return c.trustedSubnet
}

// GetGRPCPort возвращает номер порта для запуска grpc сервера
func (c Configurator) GetGRPCPort() string {
	return c.grpcPort
}
