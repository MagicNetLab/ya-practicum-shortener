package jsonreader

import (
	"encoding/json"
	"os"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

func Parse(fileName string) Configurator {
	var config Configurator

	file, err := os.Open(fileName)
	if err != nil {
		return config
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("failed to parse json config file", args)
		return Configurator{}
	}

	return config
}
