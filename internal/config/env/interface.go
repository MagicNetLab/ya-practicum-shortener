package env

// IConfigurator интерфейс хранилища параметров для работы приложения установленных в env переменных
type IConfigurator interface {
	// HasBaseHost проверяет установлен базовый хост или нет.
	HasBaseHost() bool

	// GetBaseHost возвращает установленный базовый хост для запуска приложения.
	GetBaseHost() (string, error)

	// HasShortHost проверяет установлен хост для обработки коротких ссылок или нет.
	HasShortHost() bool

	// GetShortHost возвращает хост установленный для обработки коротких ссылок.
	GetShortHost() (string, error)

	// GetBasePort проверяет установленный базовый порт для запуска приложения.
	GetBasePort() (string, error)

	// GetShortHostString возвращает адрес (хост + порт) установленный для обработки коротких ссылок.
	GetShortHostString() (string, error)

	// GetShortPort возвращает порт установленный для обработки коротких ссылок.
	GetShortPort() (string, error)

	// HasFileStoragePath проверяет установлен пусть до файла для локального хранения ссылок или нет.
	HasFileStoragePath() bool

	// GetFileStoragePath возвращает установленный пусть до файла для локального хранения ссылок.
	GetFileStoragePath() (string, error)

	// HasDBConnectString проверяет установлена ли строка конфигурации подключения к БД или нет.
	HasDBConnectString() bool

	// GetDBConnectString возвращает установленную строку с конфигурацией для подключения к БД.
	GetDBConnectString() (string, error)

	// HasJWTSecret проверяет установлен ли секрет для генерации JWT токена или нет.
	HasJWTSecret() bool

	// GetJWTSecret возвращает установленный секрет для генерации JWT токена.
	GetJWTSecret() (string, error)

	// GetPPROFHost возвращает хост для запуска web интерфейса профилировщика.
	GetPPROFHost() string
}
