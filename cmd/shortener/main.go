package main

import (
	"log"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/server"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
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

	_, err = storage.GetStore()
	if err != nil {
		logger.Log.Fatalf("Failed to initialize storage: %v", err)
	}

	server.Run(conf)
}
