package compression

import (
	"compress/gzip"
	"io"
	"net/http"
)

type gzipWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// Header получение заголовков запроса
func (c *gzipWriter) Header() http.Header {
	return c.w.Header()
}

// Write запись данных
func (c *gzipWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader запись заголовков запроса
func (c *gzipWriter) WriteHeader(statusCode int) {
	c.w.Header().Set("Content-Encoding", "gzip")
	c.w.WriteHeader(statusCode)
}

// Close закрытие врайтера
func (c *gzipWriter) Close() error {
	return c.zw.Close()
}

type gzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// Read чтение данных
func (c gzipReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close закрытие ридера
func (c *gzipReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
