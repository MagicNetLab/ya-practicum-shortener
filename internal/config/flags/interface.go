package flags

// CliConfigurator хранилище параметров установленных в cli при запуске приложения.
type CliConfigurator interface {
	// HasDefaultHost возвращает bool установлен ли базовый хост для работы приложения или нет.
	HasDefaultHost() bool

	// GetDefaultHost возвращает базовый хост установленный для работы приложения или пустую строку с ошибкой.
	GetDefaultHost() (string, error)

	// GetDefaultPort возвращает базовый порт установленный для запуска приложения или пустую строку с ошибкой
	GetDefaultPort() (string, error)

	// HasShortHost возвращает bool - установлен ли хост для обработки коротких ссылок или нет.
	HasShortHost() bool

	// GetShortHost возвращает хост установленный для обработки коротких ссылок или пустую строку с ошибкой.
	GetShortHost() (string, error)

	// GetShortPort возвращает порт установленный для обработки коротких ссылок.
	GetShortPort() (string, error)

	// HasFileStoragePath возвращает bool - установлен ли путь до файла для локального хранения ссылок.
	HasFileStoragePath() bool

	// GetFileStoragePath возвращает установленный пусть до файла для локального хранения коротких ссылок.
	GetFileStoragePath() (string, error)

	// HasDBConnectString возвращает bool - установлена ли строка конфигурации подключения к БД или нет.
	HasDBConnectString() bool

	// GetDBConnectString возвращает установленную троку конфигурации для подключения к БД или пустую строку с ошибкой.
	GetDBConnectString() (string, error)

	// HasJWTSecret возвращает bool - установлен ли секрет для генерации JWT токена.
	HasJWTSecret() bool

	// GetJWTSecret возвращает секрет для генерации JWT токена или пустую строку с ошибкой.
	GetJWTSecret() (string, error)

	// HasPProfHost возвращает bool установлен ли хост для запуска профилировщика или нет.
	HasPProfHost() bool

	// GetPProfHost возвращает хост установленный для запуска профилировщика или пустую строку с ошибкой.
	GetPProfHost() string
}
