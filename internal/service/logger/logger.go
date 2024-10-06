package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log = Logger{log: zap.NewNop()}

// Initialize инициализация логера
func Initialize() error {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		zap.InfoLevel,
	)

	logger := zap.New(core)
	log = Logger{log: logger}

	return nil
}

// Info отправка информационного сообщения в логи
func Info(msg string, args map[string]interface{}) {
	log.Info(msg, args)
}

// Error отправка сообщения об ошибке в логи
func Error(msg string, args map[string]interface{}) {
	log.Error(msg, args)
}

// Fatal отправка сообщения об ошибке в логи и завершение работы приложения
func Fatal(msg string, args map[string]interface{}) {
	log.Fatal(msg, args)
}

// Sync сброс закэшированных данных в логи
func Sync() {
	log.log.Sync()
}
