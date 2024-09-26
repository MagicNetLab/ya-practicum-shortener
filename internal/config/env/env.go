package env

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

type Config struct {
	baseHost        []string `env:"SERVER_ADDRESS" envSeparator:":"`
	shortHost       []string `env:"BASE_URL" envSeparator:":"`
	fileStoragePath string   `env:"FILE_STORAGE_PATH"`
	dbConnectString string   `env:"DATABASE_DSN"`
	jwtSecret       string   `env:"JWT_SECRET"`
	pProfHost       string   `env:"PPROF_HOST" envDefault:"localhost:5000"`
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
	return e.fileStoragePath != ""
}

func (e Config) GetFileStoragePath() (string, error) {
	if !e.HasFileStoragePath() {
		return "", errors.New("file storage path not init")
	}

	return e.fileStoragePath, nil
}

func (e Config) HasDBConnectString() bool {
	return e.dbConnectString != ""
}

func (e Config) GetDBConnectString() (string, error) {
	if !e.HasDBConnectString() {
		return "", errors.New("db connect params not init")
	}

	return e.dbConnectString, nil
}

func (e Config) HasJWTSecret() bool {
	return e.jwtSecret != ""
}

func (e Config) GetJWTSecret() (string, error) {
	if !e.HasJWTSecret() {
		return "", errors.New("jwttoken secret not init")
	}
	return e.jwtSecret, nil
}

func (e Config) GetPPROFHost() string {
	return e.pProfHost
}

func Parse() (Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error(".env file not found", args)
	}

	baseHost := os.Getenv("SERVER_ADDRESS")
	if baseHost != "" && strings.Contains(baseHost, ":") {
		envConf.baseHost = strings.Split(baseHost, ":")
	}

	shortHost := os.Getenv("BASE_URL")
	if shortHost != "" && strings.Contains(shortHost, ":") {
		envConf.shortHost = strings.Split(shortHost, ":")
	}

	if fileStorage := os.Getenv("FILE_STORAGE_PATH"); fileStorage != "" {
		envConf.fileStoragePath = fileStorage
	}

	if dbParams := os.Getenv("DATABASE_DSN"); dbParams != "" {
		envConf.dbConnectString = dbParams
	}

	if JWTSecret := os.Getenv("JWT_SECRET"); JWTSecret != "" {
		envConf.jwtSecret = JWTSecret
	}

	if pProfHost := os.Getenv("PPROF_HOST"); pProfHost != "" {
		envConf.pProfHost = pProfHost
	}

	return envConf, nil
}
