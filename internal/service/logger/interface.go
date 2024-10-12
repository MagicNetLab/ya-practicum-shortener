package logger

// AppLogger логер приложения
type AppLogger interface {
	// Info отправка в логи информационного сообщения
	Info(msg string, args map[string]interface{})

	// Error отправка в логи сообщения об ошибке
	Error(msg string, args map[string]interface{})

	// Fatal отправка в логи сообщения об ошибке и завершение работы приложения
	Fatal(msg string, args map[string]interface{})

	// Sync сброс кеша сообщений в логи при завершении работы приложения
	Sync()
}
