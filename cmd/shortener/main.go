package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/repo"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/server"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func init() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	err := logger.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
		return
	}

	err = config.Initialize()
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Fatal("Failed to initialize config", args)
	}

	err = repo.Initialize(config.GetParams())
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Fatal("Failed to initialize storage", args)
	}
}

func main() {
	server.Run(config.GetParams())
}
