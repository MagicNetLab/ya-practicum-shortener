package server

import (
	"log"
	"net/http"

	handle "github.com/MagicNetLab/ya-practicum-shortener/internal/app/server/handlers"
)

func Run(configurator configurator) {
	listen := getListeners()

	for h, l := range listen {
		h := h
		l := l
		go func() { log.Fatal(http.ListenAndServe(h, l)) }()
	}

	if pprofHost := configurator.GetPProfHost(); pprofHost != "" {
		go func() {
			log.Println(http.ListenAndServe(pprofHost, nil))
		}()
	}

	select {}
}

func getListeners() listeners {
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
