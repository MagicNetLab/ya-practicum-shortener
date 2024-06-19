package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TODO добавить тестов
func Test_encodeLinkHeader(t *testing.T) {

	type want struct {
		contentType string
		statusCode  int
		body        string
	}

	tests := []struct {
		name    string
		method  string
		body    string
		want    want
		request string
	}{
		{
			name:   "Test wrong method",
			method: http.MethodGet,
			body:   "",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusForbidden,
				body:        "Method not allowed",
			},
			request: "/",
		},
		{
			name:   "Test empty body",
			method: http.MethodPost,
			body:   "",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
				body:        "Missing link",
			},
			request: "/",
		},
		{
			name:   "Test success",
			method: http.MethodPost,
			body:   "http://yandex.ru",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
				body:        "http://localhost:8080/",
			},
			request: "/",
		},
		// todo test incorrect link
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(encodeLinkHeader)
			h(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			bodyResult, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Contains(t, string(bodyResult), tt.want.body)
		})
	}
}

// TODO добавить тестов
func Test_decodeLinkHeader(t *testing.T) {
	type want struct {
		statusCode int
	}

	tests := []struct {
		name    string
		method  string
		want    want
		request string
	}{
		{
			name:   "Test wrong method",
			method: http.MethodPost,
			want: want{
				statusCode: http.StatusForbidden,
			},
			request: "/sl/jsdhkahs",
		},
		{
			name:   "Test incorrect short link",
			method: http.MethodGet,
			want: want{
				statusCode: http.StatusNotFound,
			},
			request: "/jsdhkahs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(decodeLinkHeader)
			h(w, request)

			result := w.Result()
			err := result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
