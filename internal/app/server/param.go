package server

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type configurator interface {
	SetDefaultHost(host string, port string) error
	SetShortHost(host string, port string) error
	GetDefaultHost() string
	GetShortHost() string
	IsValid() bool
}

type route struct {
	path    string
	method  string
	handler http.HandlerFunc
}

type listeners map[string]chi.Router

func (l listeners) Append(host string, route route) {
	var r chi.Router
	if _, ok := l[host]; !ok {
		r = chi.NewRouter()
	} else {
		r = l[host]
	}

	switch route.method {
	case http.MethodPost:
		r.Post(route.path, route.handler)
	default:
		r.Get(route.path, route.handler)
	}

	l[host] = r
}
