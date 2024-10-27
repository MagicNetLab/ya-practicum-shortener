package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/server"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/storage"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	err := logger.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	conf := config.GetParams()
	if !conf.IsValid() {
		logger.Fatal("Invalid server parameters. App is not running.", nil)
		return
	}

	_, err = storage.GetStore()
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Fatal("Failed to initialize storage", args)
	}

	server.Run(conf)
}
