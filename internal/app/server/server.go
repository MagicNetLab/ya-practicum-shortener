package server

import (
	"log"
	"net/http"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"

	handle "github.com/MagicNetLab/ya-practicum-shortener/internal/app/server/handlers"
)

func Run(configurator configurator) {
	listen := getListeners()

	for h, l := range listen {
		h := h
		l := l
		go func() { log.Fatal(http.ListenAndServe(h, l)) }()

		args := map[string]interface{}{"url": h}
		logger.Info("listener starting", args)
	}

	if pprofHost := configurator.GetPProfHost(); pprofHost != "" {
		go func() {
			log.Println(http.ListenAndServe(pprofHost, nil))
		}()

		args := map[string]interface{}{"url": pprofHost}
		logger.Info("pprof starting", args)
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
