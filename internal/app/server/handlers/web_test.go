package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/jwttoken"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusMethodNotAllowed,
				body:        "Method Not Allowed",
			},
			request: "/",
		},
		{
			name:   "Test empty body",
			method: http.MethodPost,
			body:   "",
			want: want{
				contentType: "text/plain; charset=utf-8",
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
		// todo test unique link
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))
			c := config.GetParams()
			token, _ := jwttoken.GenerateToken(3, c.GetJWTSecret())
			newCookie := http.Cookie{Name: "token", Value: token, Path: "/", Expires: time.Now().Add(5 * time.Minute)}
			request.AddCookie(&newCookie)

			w := httptest.NewRecorder()
			h := encodeHandler()
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
				statusCode: http.StatusMethodNotAllowed,
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
			c := config.GetParams()
			token, _ := jwttoken.GenerateToken(3, c.GetJWTSecret())
			newCookie := http.Cookie{Name: "token", Value: token, Path: "/", Expires: time.Now().Add(5 * time.Minute)}
			request.AddCookie(&newCookie)
			w := httptest.NewRecorder()
			h := decodeHandler()
			h(w, request)

			result := w.Result()
			err := result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
