package config

import (
	"errors"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

type EnvConfType struct {
	baseHost  []string `env:"SERVER_ADDRESS" envSeparator:":"`
	shortHost []string `env:"BASE_URL" envSeparator:":"`
}

var envConf = EnvConfType{}

func (e *EnvConfType) HasBaseHost() bool {
	return len(envConf.baseHost) == 2
}

func (e *EnvConfType) HasShortHost() bool {
	return len(envConf.shortHost) == 2
}

func (e *EnvConfType) GetBaseHost() (string, error) {
	if len(envConf.baseHost) != 2 {
		return "", errors.New("base host not init in env")
	}

	if envConf.baseHost[0] == "" {
		return "", errors.New("base host is empty")
	}

	return envConf.baseHost[0], nil
}

func (e *EnvConfType) GetBasePort() (string, error) {
	if len(envConf.baseHost) != 2 {
		return "", errors.New("base host not init in env")
	}

	if envConf.baseHost[0] == "" {
		return "", errors.New("base port is empty")
	}

	return envConf.baseHost[1], nil
}

func (e *EnvConfType) GetShortHostString() (string, error) {
	if len(envConf.shortHost) != 2 {
		return "", errors.New("base host not init in env")
	}

	return strings.Join(envConf.shortHost, ":"), nil
}

func (e *EnvConfType) GetShortHost() (string, error) {
	if len(envConf.shortHost) != 2 {
		return "", errors.New("base host not init in env")
	}

	return envConf.shortHost[0], nil
}

func (e *EnvConfType) GetShortPort() (string, error) {
	if len(envConf.shortHost) != 2 {
		return "", errors.New("base host not init in env")
	}

	return envConf.shortHost[1], nil
}

func ReadEnv() (EnvConfType, error) {
	if err := godotenv.Load(".env"); err != nil {
		err := env.Parse(&envConf)
		if err != nil {
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
