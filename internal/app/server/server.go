package server

import (
	handle "github.com/MagicNetLab/ya-practicum-shortener/internal/app/server/handlers"
	"log"
	"net/http"
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

	defaultRoute := route{
		path:    handlers["default"].Path,
		method:  handlers["default"].Method,
		handler: handlers["default"].Handler,
	}
	shortRoute := route{
		path:    handlers["short"].Path,
		method:  handlers["shot"].Method,
		handler: handlers["short"].Handler,
	}

	l := make(listeners)
	l.Append(configurator.GetDefaultHost(), defaultRoute)
	l.Append(configurator.GetShortHost(), shortRoute)

	return l
}
