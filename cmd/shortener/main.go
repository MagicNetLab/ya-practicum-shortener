package main

import (
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/server"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"log"
)

func main() {
	conf := config.GetParams()
	if !conf.IsValid() {
		config.Testing()
		conf = config.GetParams()
	}

	if err := server.RunServer(conf); err != nil {
		log.Fatalf("Failed starting server :%s", err)
	}
}
