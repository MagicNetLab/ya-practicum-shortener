package logger

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

// Log TODO Задача со звездочкой: а что если ты используешь интерфейс для логирования, что позволит легко заменить
//			 реализацию логирования в будущем или сделать его более настраиваемым через конфигурацию

var Log *zap.SugaredLogger = zap.NewNop().Sugar()

func Initialize() error {
	zl, err := zap.NewProduction()
	if err != nil {
		return err
	}

	Log = zl.Sugar()
	return nil
}

func Middleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &ResponseData{
			Status: 0,
			Size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		Log.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.Status,
			"duration", duration,
			"size", responseData.Size,
		)
	}
}
