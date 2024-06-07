package main

import (
	serv "github.com/MagicNetLab/ya-practicum-shortener/internal/app/server"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/server/config"
)

func main() {
	conf := config.GetParams()
	err := serv.RunServer(conf)
	if err != nil {
		panic(err)
	}
}
