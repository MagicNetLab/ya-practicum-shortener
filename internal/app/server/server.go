package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/repo"
	grpcServer "github.com/MagicNetLab/ya-practicum-shortener/internal/app/server/grpc"
	restServer "github.com/MagicNetLab/ya-practicum-shortener/internal/app/server/rest"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
	pb "github.com/MagicNetLab/ya-practicum-shortener/pkg/shortener_grpc"
)

const (
	certFileName = "cert.pem"
	keyFileName  = "key.pem"
)

// Run запуск сервера
func Run(configurator configurator) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	var runServers []*http.Server
	if configurator.IsEnableHTTPS() {
		runServers = runHTTPSServer(configurator)
	} else {
		runServers = runHTTPServer(configurator)
	}

	grpcServ := runGRPCServer(configurator)

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

	grpcServ.GracefulStop()

	logger.Info("Server exited", nil)
}

func getListeners() listeners {
	handlers := restServer.GetHandlers()
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

func runHTTPServer(configurator configurator) []*http.Server {
	var runServers []*http.Server
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

	return runServers
}

func runHTTPSServer(configurator configurator) []*http.Server {
	var runServers []*http.Server

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

	return runServers
}

func runGRPCServer(configurator configurator) *grpc.Server {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", ":"+configurator.GetGRPCPort())
	if err != nil {
		log.Fatal(err)
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			grpcServer.LoggerInterceptor,
			grpcServer.TrustedNetworkInterceptor,
			grpcServer.GuestInterceptor,
			grpcServer.AuthInterceptor,
		),
	}

	s := grpc.NewServer(opts...)
	pb.RegisterGrpcShortenerServer(s, &grpcServer.Server{})

	fmt.Println("grpc started listening on", listen.Addr().String())

	go func() {
		err := s.Serve(listen)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return s
}

func runPProfServer(configurator configurator) *http.Server {
	pprofHost := configurator.GetPProfHost()
	pprofServer := &http.Server{
		Addr:    pprofHost,
		Handler: nil,
	}

	if configurator.IsEnableHTTPS() {
		go func(s *http.Server) {
			err := s.ListenAndServeTLS(certFileName, keyFileName)
			if err != nil && err != http.ErrServerClosed {
				args := map[string]interface{}{"error": err}
				logger.Fatal("failed starting pprof server", args)
			}
		}(pprofServer)
	} else {
		go func(serv *http.Server) {
			err := serv.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				args := map[string]interface{}{"error": err.Error()}
				logger.Error("Failed start pprof server", args)
			}
		}(pprofServer)
	}

	args := map[string]interface{}{"url": "http://" + pprofHost}
	logger.Info("pprof starting", args)

	return pprofServer
}
