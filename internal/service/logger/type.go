package logger

import (
	"net/http"

	"go.uber.org/zap"
)

// Logger объект логера для приложения
type Logger struct {
	log *zap.Logger
}

// Info отправка информационного сообщения в логи
func (l *Logger) Info(msg string, args map[string]interface{}) {
	l.log.Info(msg, l.prepareArgs(args)...)
}

// Error отправка сообщения об ошибке в логи
func (l *Logger) Error(msg string, args map[string]interface{}) {
	l.log.Error(msg, l.prepareArgs(args)...)
}

// Fatal отправка сообщения об ошибке в логи и завершения работы приложения
func (l *Logger) Fatal(msg string, args map[string]interface{}) {
	l.log.Fatal(msg, l.prepareArgs(args)...)
}

func (l *Logger) prepareArgs(args map[string]interface{}) []zap.Field {
	var r []zap.Field
	for k, v := range args {
		r = append(r, zap.Any(k, v))
	}

	return r
}

// Sync сброс закешированных сообщений в логи
func (l *Logger) Sync() {
	l.log.Sync()
}

type responseData struct {
	Status int
	Size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// Write запись данных о запросе
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.Size += size
	return size, err
}

// WriteHeader запись заголовков запроса
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.Status = statusCode
}
