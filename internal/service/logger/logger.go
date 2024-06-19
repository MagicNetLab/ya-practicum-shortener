package logger

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

var Log *zap.SugaredLogger = zap.NewNop().Sugar()

func Initialize() error {
	zl, err := zap.NewProduction()
	if err != nil {
		return err
	}

	Log = zl.Sugar()
	return nil
}

func RequestLogger(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
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
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	}
}
