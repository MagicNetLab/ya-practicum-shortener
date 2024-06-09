package main

import (
	serv "github.com/MagicNetLab/ya-practicum-shortener/internal/app/server"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/server/config"
	"log"
)

func main() {
	conf := config.GetParams()
	if err := serv.RunServer(conf); err != nil {
		log.Fatalf("Failed starting server :%s", err)
	}
}
