package env

import (
	"errors"
	"strings"
)

type configurator struct {
	baseHost        []string `env:"SERVER_ADDRESS" envSeparator:":"`
	shortHost       []string `env:"BASE_URL" envSeparator:":"`
	fileStoragePath string   `env:"FILE_STORAGE_PATH"`
	dbConnectString string   `env:"DATABASE_DSN"`
	jwtSecret       string   `env:"JWT_SECRET"`
	pProfHost       string   `env:"PPROF_HOST" envDefault:"localhost:5000"`
}

func (e configurator) HasBaseHost() bool {
	return len(e.baseHost) == 2
}

func (e configurator) HasShortHost() bool {
	return len(e.shortHost) == 2
}

func (e configurator) GetBaseHost() (string, error) {
	if !e.HasBaseHost() {
		return "", errors.New("base host not init in env")
	}

	if e.baseHost[0] == "" {
		return "", errors.New("base host is empty")
	}

	return e.baseHost[0], nil
}

func (e configurator) GetBasePort() (string, error) {
	if !e.HasBaseHost() {
		return "", errors.New("base host not init in env")
	}

	if e.baseHost[1] == "" {
		return "", errors.New("base port is empty")
	}

	return e.baseHost[1], nil
}

func (e configurator) GetShortHostString() (string, error) {
	if !e.HasShortHost() {
		return "", errors.New("base host not init in env")
	}

	return strings.Join(e.shortHost, ":"), nil
}

func (e configurator) GetShortHost() (string, error) {
	if !e.HasShortHost() {
		return "", errors.New("base host not init in env")
	}

	return e.shortHost[0], nil
}

func (e configurator) GetShortPort() (string, error) {
	if !e.HasShortHost() {
		return "", errors.New("base host not init in env")
	}

	return e.shortHost[1], nil
}

func (e configurator) HasFileStoragePath() bool {
	return e.fileStoragePath != ""
}

func (e configurator) GetFileStoragePath() (string, error) {
	if !e.HasFileStoragePath() {
		return "", errors.New("file storage path not init")
	}

	return e.fileStoragePath, nil
}

func (e configurator) HasDBConnectString() bool {
	return e.dbConnectString != ""
}

func (e configurator) GetDBConnectString() (string, error) {
	if !e.HasDBConnectString() {
		return "", errors.New("db connect params not init")
	}

	return e.dbConnectString, nil
}

func (e configurator) HasJWTSecret() bool {
	return e.jwtSecret != ""
}

func (e configurator) GetJWTSecret() (string, error) {
	if !e.HasJWTSecret() {
		return "", errors.New("jwttoken secret not init")
	}
	return e.jwtSecret, nil
}

func (e configurator) GetPPROFHost() string {
	return e.pProfHost
}
