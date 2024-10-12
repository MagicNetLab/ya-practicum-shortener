package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type configurator interface {
	GetDefaultHost() string
	GetShortHost() string
	IsValid() bool
	GetPProfHost() string
}

type route struct {
	path    string
	method  string
	handler http.HandlerFunc
}

type listeners map[string]chi.Router

func (l listeners) append(host string, route route) {
	var r chi.Router
	if _, ok := l[host]; !ok {
		r = chi.NewRouter()
	} else {
		r = l[host]
	}

	switch route.method {
	case http.MethodPost:
		r.Post(route.path, route.handler)
	case http.MethodDelete:
		r.Delete(route.path, route.handler)
	default:
		r.Get(route.path, route.handler)
	}

	l[host] = r
}
