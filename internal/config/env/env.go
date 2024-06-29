package env

import (
	"errors"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

type Config struct {
	baseHost    []string `env:"SERVER_ADDRESS" envSeparator:":"`
	shortHost   []string `env:"BASE_URL" envSeparator:":"`
	fileStorage string   `env:"FILE_STORAGE_PATH"`
}

var envConf = Config{}

func (e *Config) HasBaseHost() bool {
	return len(e.baseHost) == 2
}

func (e *Config) HasShortHost() bool {
	return len(e.shortHost) == 2
}

func (e *Config) GetBaseHost() (string, error) {
	if !e.HasBaseHost() {
		return "", errors.New("base host not init in env")
	}

	if e.baseHost[0] == "" {
		return "", errors.New("base host is empty")
	}

	return e.baseHost[0], nil
}

func (e *Config) GetBasePort() (string, error) {
	if !e.HasBaseHost() {
		return "", errors.New("base host not init in env")
	}

	if e.baseHost[1] == "" {
		return "", errors.New("base port is empty")
	}

	return e.baseHost[1], nil
}

func (e *Config) GetShortHostString() (string, error) {
	if !e.HasShortHost() {
		return "", errors.New("base host not init in env")
	}

	return strings.Join(e.shortHost, ":"), nil
}

func (e *Config) GetShortHost() (string, error) {
	if !e.HasShortHost() {
		return "", errors.New("base host not init in env")
	}

	return e.shortHost[0], nil
}

func (e *Config) GetShortPort() (string, error) {
	if !e.HasShortHost() {
		return "", errors.New("base host not init in env")
	}

	return e.shortHost[1], nil
}

func (e Config) HasFileStoragePath() bool {
	return e.fileStorage != ""
}

func (e Config) GetFileStoragePath() (string, error) {
	if !e.HasFileStoragePath() {
		return "", errors.New("file storage path not init")
	}

	return e.fileStorage, nil
}

func Parse() (Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf(".env file not found: %s", err)

		err := env.Parse(&envConf)
		if err != nil {
			log.Printf("Failed to parse .env file: %s", err)
			return envConf, err
		}

		return envConf, nil
	}

	baseHost := os.Getenv("SERVER_ADDRESS")
	if baseHost != "" && strings.Contains(baseHost, ":") {
		envConf.baseHost = strings.Split(baseHost, ":")
	}

	shortHost := os.Getenv("BASE_URL")
	if shortHost != "" && strings.Contains(shortHost, ":") {
		envConf.shortHost = strings.Split(shortHost, ":")
	}

	fileStorage := os.Getenv("FILE_STORAGE_PATH")
	if fileStorage != "" {
		envConf.fileStorage = fileStorage
	}

	return envConf, nil
}
