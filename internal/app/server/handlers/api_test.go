package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/jwttoken"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_apiEncodeHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		body        string
	}

	tests := []struct {
		name    string
		method  string
		cookie  bool
		body    string
		want    want
		request string
	}{
		{
			name:   "Test wrong method",
			method: http.MethodGet,
			cookie: true,
			body:   "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusForbidden,
				body:        "Method not allowed",
			},
			request: "/api/shorten",
		},
		{
			name:   "Test empty body",
			method: http.MethodPost,
			cookie: true,
			body:   "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				body:        "Missing link",
			},
			request: "/api/shorten",
		},
		{
			name:   "Test success",
			method: http.MethodPost,
			cookie: true,
			body:   "{\"url\": \"https://practicum.yandex.ru\"}",
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
				body:        `{"result":"http://localhost:8080/`,
			},
			request: "/api/shorten",
		},
		{
			name:   "Test missing link",
			method: http.MethodPost,
			cookie: true,
			body:   "{\"url\": \"\"}",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				body:        `Missing link`,
			},
			request: "/api/shorten",
		},
		{
			name:   "Test without cookie header",
			method: http.MethodPost,
			cookie: false,
			body:   "{\"url\": \"https://practicum.yandex.ru\"}",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				body:        `incorrect user token`,
			},
			request: "/api/shorten",
		},
	}

	c := config.GetParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))
			if tt.cookie {
				token, _ := jwttoken.GenerateToken(3, c.GetJWTSecret())
				newCookie := http.Cookie{Name: "token", Value: token, Path: "/", Expires: time.Now().Add(5 * time.Minute)}
				request.AddCookie(&newCookie)
			}

			w := httptest.NewRecorder()
			h := apiEncodeHandler()
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

func Test_apiEncodeLinkByUnique(t *testing.T) {
	t.Run("test send not unique link", func(t *testing.T) {
		store, err := storage.GetStore()
		assert.NoError(t, err)
		link := "https://mail.ru"
		userID := 3

		err = store.PutLink(link, "uweyiu", userID)
		assert.NoError(t, err)

		c := config.GetParams()
		t.Run("test not unique link", func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader("{\"url\": \"https://mail.ru\"}"))
			token, _ := jwttoken.GenerateToken(3, c.GetJWTSecret())
			newCookie := http.Cookie{Name: "token", Value: token, Path: "/", Expires: time.Now().Add(5 * time.Minute)}
			request.AddCookie(&newCookie)

			w := httptest.NewRecorder()
			h := apiEncodeHandler()
			h(w, request)

			result := w.Result()

			assert.Equal(t, http.StatusConflict, result.StatusCode)
			assert.Equal(t, "application/json", result.Header.Get("Content-Type"))

			bodyResult, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Contains(t, string(bodyResult), `{"result":"http://localhost:8080/`)
		})
	})
}

// TODO test apiBatchEncodeHandler
func Test_apiBatchEncodeHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		body        string
	}

	tests := []struct {
		name    string
		method  string
		cookie  bool
		body    string
		want    want
		request string
	}{
		{
			name:   "Test wrong method",
			method: http.MethodGet,
			cookie: true,
			body:   "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusForbidden,
				body:        "Method not allowed",
			},
			request: "/api/shorten/batch",
		},
		{
			name:   "Test empty body",
			method: http.MethodPost,
			cookie: true,
			body:   "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				body:        "Missing link",
			},
			request: "/api/shorten/batch",
		},
		{
			name:   "Test success",
			method: http.MethodPost,
			cookie: true,
			body:   "[{\"correlation_id\":\"364asds\",\"original_url\":\"https://okko.ru\"}]",
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
				body:        `"correlation_id":"364asds"`,
			},
			request: "/api/shorten",
		},
		{
			name:   "Test without cookie header",
			method: http.MethodPost,
			cookie: false,
			body:   "{\"url\": \"https://practicum.yandex.ru\"}",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				body:        `incorrect user token`,
			},
			request: "/api/shorten/batch",
		},
		// test missing link
	}

	c := config.GetParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))
			if tt.cookie {
				token, _ := jwttoken.GenerateToken(3, c.GetJWTSecret())
				newCookie := http.Cookie{Name: "token", Value: token, Path: "/", Expires: time.Now().Add(5 * time.Minute)}
				request.AddCookie(&newCookie)
			}

			w := httptest.NewRecorder()
			h := apiBatchEncodeHandler()
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
