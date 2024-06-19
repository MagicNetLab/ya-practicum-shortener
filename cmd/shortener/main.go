package main

import (
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/server"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
	"log"
)

func main() {
	err := logger.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	conf := config.GetParams()
	if !conf.IsValid() {
		logger.Log.Fatalln("Invalid server parameters. App is not running.")
		return
	}

	server.Run(conf)
}
