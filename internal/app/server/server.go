package server

import (
	"log"
	"net/http"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"

	handle "github.com/MagicNetLab/ya-practicum-shortener/internal/app/server/handlers"
)

// Run запуск сервера
func Run(configurator configurator) {
	if configurator.IsEnableHTTPS() {
		runHTTPSServer(configurator)
	} else {
		runHTTPServer(configurator)
	}
}

func getListeners() listeners {
	handlers := handle.GetHandlers()
	l := make(listeners)
	for _, v := range handlers {
		l.append(v.Host, route{
			path:    v.Path,
			method:  v.Method,
			handler: v.Handler,
		})
	}
	return l
}

func runHTTPServer(configurator configurator) {
	listen := getListeners()

	for h, l := range listen {
		h := h
		l := l
		go func() { log.Fatal(http.ListenAndServe(h, l)) }()

		args := map[string]interface{}{"url": "http://" + h}
		logger.Info("listener starting", args)
	}

	if pprofHost := configurator.GetPProfHost(); pprofHost != "" {
		go func() {
			log.Println(http.ListenAndServe(pprofHost, nil))
		}()

		args := map[string]interface{}{"url": "http://" + pprofHost}
		logger.Info("pprof starting", args)
	}

	select {}
}

func runHTTPSServer(configurator configurator) {
	const (
		certFileName = "cert.pem"
		keyFileName  = "key.pem"
	)

	if !hasTLSCertsExists(certFileName, keyFileName) {
		err := createTLSCerts(certFileName, keyFileName)
		if err != nil {
			params := map[string]interface{}{"error": err}
			logger.Fatal("error creating tls certs", params)
		}
	}

	handlers := getListeners()
	for key, value := range handlers {
		host := key
		handler := value
		go func() {
			log.Fatal(http.ListenAndServeTLS(host, certFileName, keyFileName, handler))
		}()

		args := map[string]interface{}{"url": "https://" + host}
		logger.Info("listener starting", args)
	}

	if pprofHost := configurator.GetPProfHost(); pprofHost != "" {
		go func() {
			log.Println(http.ListenAndServeTLS(pprofHost, certFileName, keyFileName, nil))
		}()

		args := map[string]interface{}{"url": "https://" + pprofHost}
		logger.Info("pprof starting", args)
	}

	select {}

}
