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

	defaultRoute := route{
		path:    handlers["default"].Path,
		method:  handlers["default"].Method,
		handler: handlers["default"].Handler,
	}
	apiDefaultRoute := route{
		path:    handlers["apiDefault"].Path,
		method:  handlers["apiDefault"].Method,
		handler: handlers["apiDefault"].Handler,
	}
	apiBatchDefaultRoute := route{
		path:    handlers["apiBatchDefault"].Path,
		method:  handlers["apiBatchDefault"].Method,
		handler: handlers["apiBatchDefault"].Handler,
	}
	shortRoute := route{
		path:    handlers["short"].Path,
		method:  handlers["shot"].Method,
		handler: handlers["short"].Handler,
	}
	pingRoute := route{
		path:    handlers["dbPing"].Path,
		method:  handlers["dbPing"].Method,
		handler: handlers["dbPing"].Handler,
	}
	userLinksRoute := route{
		path:    handlers["apiUserLinks"].Path,
		method:  handlers["apiUserLinks"].Method,
		handler: handlers["apiUserLinks"].Handler,
	}

	l := make(listeners)
	l.Append(configurator.GetDefaultHost(), defaultRoute)
	l.Append(configurator.GetDefaultHost(), apiDefaultRoute)
	l.Append(configurator.GetShortHost(), shortRoute)
	l.Append(configurator.GetDefaultHost(), pingRoute)
	l.Append(configurator.GetDefaultHost(), apiBatchDefaultRoute)
	l.Append(configurator.GetDefaultHost(), userLinksRoute)

	return l
}
