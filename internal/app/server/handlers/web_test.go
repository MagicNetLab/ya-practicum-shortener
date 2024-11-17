package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/repo"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/jwttoken"
)

type contextKey struct {
	name string
}

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
		userID  int
		cookies bool
		want    want
		request string
	}{
		{
			name:    "Test wrong method",
			method:  http.MethodGet,
			body:    "",
			userID:  28,
			cookies: true,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusMethodNotAllowed,
				body:        "Method Not Allowed",
			},
			request: "/",
		},
		{
			name:    "Test empty body",
			method:  http.MethodPost,
			body:    "",
			userID:  33,
			cookies: true,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				body:        "Missing link",
			},
			request: "/",
		},
		{
			name:    "Test success",
			method:  http.MethodPost,
			body:    "http://yandex.ru",
			userID:  23,
			cookies: true,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
				body:        "http://localhost:8080/",
			},
			request: "/",
		},
		{
			name:    "Test empty body",
			method:  http.MethodPost,
			body:    "",
			userID:  48,
			cookies: false,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				body:        "incorrect user token",
			},
			request: "/",
		},
	}

	c := config.GetParams()
	if !c.IsValid() {
		errConf := config.Initialize()
		assert.NoError(t, errConf)
		errConf = repo.Initialize(config.GetParams())
		assert.NoError(t, errConf)
		c = config.GetParams()
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))
			cookieName := "token"
			if tt.cookies == false {
				cookieName = "tokien"
			}
			token, _ := jwttoken.GenerateToken(tt.userID, c.GetJWTSecret())
			newCookie := http.Cookie{Name: cookieName, Value: token, Path: "/", Expires: time.Now().Add(5 * time.Minute)}
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

func Test_encodeLinkByUnique(t *testing.T) {
	t.Run("test send not unique link", func(t *testing.T) {
		c := config.GetParams()
		if !c.IsValid() {
			errConf := config.Initialize()
			assert.NoError(t, errConf)
			errConf = repo.Initialize(config.GetParams())
			assert.NoError(t, errConf)
			c = config.GetParams()
		}

		link := "https://cloud.ru"
		userID := 3

		err := repo.PutLink(context.Background(), link, "uweyiu", userID)
		assert.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(link))
		token, _ := jwttoken.GenerateToken(userID, c.GetJWTSecret())
		newCookie := http.Cookie{Name: "token", Value: token, Path: "/", Expires: time.Now().Add(5 * time.Minute)}
		request.AddCookie(&newCookie)

		w := httptest.NewRecorder()
		h := encodeHandler()
		h(w, request)

		result := w.Result()

		assert.Equal(t, http.StatusConflict, result.StatusCode)
		assert.Equal(t, "text/plain", result.Header.Get("Content-Type"))

		bodyResult, err := io.ReadAll(result.Body)
		require.NoError(t, err)
		err = result.Body.Close()
		require.NoError(t, err)

		assert.Contains(t, string(bodyResult), "http://localhost:8080/")
	})
}

func Test_decodeLinkHeader(t *testing.T) {
	type want struct {
		statusCode int
	}

	tests := []struct {
		name    string
		method  string
		deleted bool
		want    want
		request string
	}{
		{
			name:    "Test wrong method",
			method:  http.MethodPost,
			deleted: false,
			want: want{
				statusCode: http.StatusMethodNotAllowed,
			},
			request: "/sl/jsdhkahs",
		},
		{
			name:    "Test incorrect short link",
			method:  http.MethodGet,
			deleted: false,
			want: want{
				statusCode: http.StatusNotFound,
			},
			request: "/jsdhkahs",
		},

		// todo short link is deleted
		//{
		//	name:    "Test short link is deleted",
		//	method:  http.MethodGet,
		//	deleted: true,
		//	want: want{
		//		statusCode: http.StatusGone,
		//	},
		//	request: "/jsdhkahs",
		//},
		// todo success get link
	}

	c := config.GetParams()
	if !c.IsValid() {
		errConf := config.Initialize()
		assert.NoError(t, errConf)
		errConf = repo.Initialize(config.GetParams())
		assert.NoError(t, errConf)
		c = config.GetParams()
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.deleted {
				err := repo.PutLink(context.Background(), "http://test.link", "jsdhkahs", 3)
				assert.NoError(t, err)

				err = repo.DeleteBatchLinksArray(context.Background(), []string{"jsdhkahs"}, 3)
				assert.NoError(t, err)
			}

			request := httptest.NewRequest(tt.method, tt.request, nil)
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
