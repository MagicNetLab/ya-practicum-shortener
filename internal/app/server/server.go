package server

import (
	handle "github.com/MagicNetLab/ya-practicum-shortener/internal/app/server/handlers"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func RunServer(serverConfig configurator) error {
	defaultHost := serverConfig.GetDefaultHost()
	shortHost := serverConfig.GetShortHost()
	handlers := handle.GetHandlers()

	if defaultHost == shortHost {
		router := chi.NewRouter()
		for _, route := range handlers {
			if route.Method == http.MethodPost {
				router.Post(route.Path, route.Handler)
			} else {
				router.Get(route.Path, route.Handler)
			}
		}

		err := http.ListenAndServe(serverConfig.GetDefaultHost(), router)
		if err != nil {
			panic(err)
		}
	} else {
		defaultRouter := chi.NewRouter()
		defaultRouter.Post(handlers["default"].Path, handlers["default"].Handler)

		shortRouter := chi.NewRouter()
		shortRouter.Get(handlers["short"].Path, handlers["short"].Handler)

		go func() { log.Fatal(http.ListenAndServe(serverConfig.GetDefaultHost(), defaultRouter)) }()
		go func() { log.Fatal(http.ListenAndServe(serverConfig.GetShortHost(), shortRouter)) }()

		select {}
	}

	return nil
}
