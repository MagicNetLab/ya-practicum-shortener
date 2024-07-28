package server

import (
	"log"
	"net/http"

	handle "github.com/MagicNetLab/ya-practicum-shortener/internal/app/server/handlers"
)

func Run(configurator configurator) {
	listen := getListeners(configurator)

	for h, l := range listen {
		h := h
		l := l
		go func() { log.Fatal(http.ListenAndServe(h, l)) }()
	}

	select {}
}

func getListeners(configurator configurator) listeners {
	handlers := handle.GetHandlers()
	l := make(listeners)
	for _, v := range handlers {
		l.Append(v.Host, route{
			path:    v.Path,
			method:  v.Method,
			handler: v.Handler,
		})
	}
	return l
}
