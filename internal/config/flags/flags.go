package flags

import (
	"errors"
	"flag"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
	"os"
	"strings"
)

const (
	defaultHostKey  = "a"
	shortHostKey    = "b"
	fileStoragePath = "f"
)

// todo (Лучше объявить defaultHost и shortHost как часть структуры cliConf, а не глобальные переменные.)
var defaultHost = ""
var shortHost = ""
var fileStorage = ""

type CliConfigurator interface {
	HasDefaultHost() bool
	GetDefaultHost() (string, error)
	GetDefaultPort() (string, error)
	HasShortHost() bool
	GetShortHost() (string, error)
	GetShortPort() (string, error)
	HasFileStoragePath() bool
	GetFileStoragePath() (string, error)
}

type cliConf struct {
	defaultHost     string
	defaultPort     string
	shortHost       string
	shortPort       string
	fileStoragePath string
}

func (cc *cliConf) HasDefaultHost() bool {
	return cc.defaultHost != "" && cc.defaultPort != ""
}

func (cc *cliConf) GetDefaultHost() (string, error) {
	if cc.defaultHost == "" {
		return "", errors.New("default host not set")
	}

	return cc.defaultHost, nil
}

func (cc *cliConf) GetDefaultPort() (string, error) {
	if cc.defaultPort == "" {
		return "", errors.New("default port not set")
	}

	return cc.defaultPort, nil
}

func (cc *cliConf) HasShortHost() bool {
	return cc.shortHost != "" && cc.shortPort != ""
}

func (cc *cliConf) GetShortHost() (string, error) {
	if cc.shortHost == "" {
		return "", errors.New("short host not set")
	}

	return cc.shortHost, nil
}

func (cc *cliConf) GetShortPort() (string, error) {
	if cc.shortPort == "" {
		return "", errors.New("short port not set")
	}

	return cc.shortPort, nil
}

func (cc *cliConf) HasFileStoragePath() bool {
	return cc.fileStoragePath != ""
}

func (cc *cliConf) GetFileStoragePath() (string, error) {
	if !cc.HasFileStoragePath() {
		return "", errors.New("file storage path not set")
	}

	return cc.fileStoragePath, nil
}

var conf = cliConf{}

func Parse() CliConfigurator {
	flag.StringVar(&defaultHost, defaultHostKey, "", "Base address")
	flag.StringVar(&shortHost, shortHostKey, "", "short links host")
	flag.StringVar(&fileStorage, fileStoragePath, "", "file storage path")
	flag.Parse()

	dh := strings.Split(defaultHost, ":")
	if len(dh) == 2 {
		conf.defaultHost = dh[0]
		conf.defaultPort = dh[1]
	}

	sh := strings.Split(shortHost, ":")
	if len(sh) == 2 {
		conf.shortHost = sh[0]
		conf.shortPort = sh[1]
	}

	conf.fileStoragePath = fileStorage
	logger.Log.Infof("os Args: %s", strings.Join(os.Args, " // "))

	return &conf
}
