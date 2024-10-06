package config

// ParameterConfig хранилище параметров для запуска и работы приложения
type ParameterConfig interface {
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
	// IsValid проверяет корректность настроек для работы проложения
	IsValid() bool
}
