package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/repo"
	handle "github.com/MagicNetLab/ya-practicum-shortener/internal/app/server/handlers"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

var runServers []*http.Server

// Run запуск сервера
func Run(configurator configurator) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	if configurator.IsEnableHTTPS() {
		runHTTPSServer(configurator)
	} else {
		runHTTPServer(configurator)
	}

	<-ctx.Done()

	logger.Info("Server stopped", nil)
	shutDownCTX, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Закрываем БД + сбрасываем логи из буфер
	defer func() {
		err := repo.Close()
		if err != nil {
			args := map[string]interface{}{"errpr": err.Error()}
			logger.Error("Failed close store", args)
		}
		logger.Sync()
	}()

	for _, server := range runServers {
		if err := server.Shutdown(shutDownCTX); err != nil {
			args := map[string]interface{}{"errpr": err.Error()}
			logger.Error("Failed shutdown server", args)
		}
	}

	logger.Info("Server exited", nil)
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
		server := &http.Server{
			Addr:    h,
			Handler: l,
		}

		s := server
		go func(s *http.Server) {
			err := s.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				args := map[string]interface{}{"error": err.Error()}
				logger.Fatal("Failed start http server", args)
			}
		}(s)

		runServers = append(runServers, server)

		args := map[string]interface{}{"url": "http://" + h}
		logger.Info("listener starting", args)
	}

	if pprofHost := configurator.GetPProfHost(); pprofHost != "" {
		pprofServer := &http.Server{
			Addr:    pprofHost,
			Handler: nil,
		}

		runServers = append(runServers, pprofServer)

		go func(serv *http.Server) {
			err := serv.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				args := map[string]interface{}{"error": err.Error()}
				logger.Error("Failed start pprof server", args)
			}
		}(pprofServer)

		args := map[string]interface{}{"url": "http://" + pprofHost}
		logger.Info("pprof starting", args)
	}
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
		server := &http.Server{
			Addr:    host,
			Handler: handler,
		}

		runServers = append(runServers, server)

		serv := server
		go func(s *http.Server) {
			err := s.ListenAndServeTLS(certFileName, keyFileName)
			if err != nil && err != http.ErrServerClosed {
				args := map[string]interface{}{"error": err}
				logger.Fatal("failed starting server", args)
			}

		}(serv)

		args := map[string]interface{}{"url": "https://" + host}
		logger.Info("listener starting", args)
	}

	if pprofHost := configurator.GetPProfHost(); pprofHost != "" {
		pprofServer := &http.Server{
			Addr:    pprofHost,
			Handler: nil,
		}

		runServers = append(runServers, pprofServer)

		go func(s *http.Server) {
			err := s.ListenAndServeTLS(certFileName, keyFileName)
			if err != nil && err != http.ErrServerClosed {
				args := map[string]interface{}{"error": err}
				logger.Fatal("failed starting pprof server", args)
			}
		}(pprofServer)

		args := map[string]interface{}{"url": "https://" + pprofHost}
		logger.Info("pprof starting", args)
	}
}
