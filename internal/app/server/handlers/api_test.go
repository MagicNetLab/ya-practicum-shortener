package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
			request: "/api/shorten",
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
			request: "/api/shorten",
		},
		{
			name:   "Test success",
			method: http.MethodPost,
			body:   "{\"url\": \"https://practicum.yandex.ru\"}",
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
				body:        `{"result":"http://localhost:8080/`,
			},
			request: "/api/shorten",
		},
		// todo test incorrect link
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))
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

// TODO test apiBatchEncodeHandler
