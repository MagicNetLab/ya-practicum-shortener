package config

import (
	"errors"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

type EnvConfig struct {
	baseHost  []string `env:"SERVER_ADDRESS" envSeparator:":"`
	shortHost []string `env:"BASE_URL" envSeparator:":"`
}

var envConf = EnvConfig{}

func (e *EnvConfig) HasBaseHost() bool {
	return len(e.baseHost) == 2
}

func (e *EnvConfig) HasShortHost() bool {
	return len(e.shortHost) == 2
}

func (e *EnvConfig) GetBaseHost() (string, error) {
	if !e.HasBaseHost() {
		return "", errors.New("base host not init in env")
	}

	if e.baseHost[0] == "" {
		return "", errors.New("base host is empty")
	}

	return e.baseHost[0], nil
}

func (e *EnvConfig) GetBasePort() (string, error) {
	if e.HasBaseHost() {
		return "", errors.New("base host not init in env")
	}

	if e.baseHost[0] == "" {
		return "", errors.New("base port is empty")
	}

	return e.baseHost[1], nil
}

func (e *EnvConfig) GetShortHostString() (string, error) {
	if !e.HasShortHost() {
		return "", errors.New("base host not init in env")
	}

	return strings.Join(e.shortHost, ":"), nil
}

func (e *EnvConfig) GetShortHost() (string, error) {
	if !e.HasShortHost() {
		return "", errors.New("base host not init in env")
	}

	return e.shortHost[0], nil
}

func (e *EnvConfig) GetShortPort() (string, error) {
	if !e.HasShortHost() {
		return "", errors.New("base host not init in env")
	}

	return e.shortHost[1], nil
}

func ReadEnv() (EnvConfig, error) {
	if err := godotenv.Load(".env"); err != nil {
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

	return envConf, nil
}
