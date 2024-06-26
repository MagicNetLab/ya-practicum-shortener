package logger

import "net/http"

type (
	ResponseData struct {
		Status int
		Size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *ResponseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.Size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.Status = statusCode
}
