package flags

import (
	"errors"
	"flag"
	"strings"
)

const (
	defaultHostKey = "a"
	shortHostKey   = "b"
)

var defaultHost = ""
var shortHost = ""

type CliConfigurator interface {
	HasDefaultHost() bool
	GetDefaultHost() (string, error)
	GetDefaultPort() (string, error)
	HasShortHost() bool
	GetShortHost() (string, error)
	GetShortPort() (string, error)
}

type cliConf struct {
	defaultHost string
	defaultPort string
	shortHost   string
	shortPort   string
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

var conf = cliConf{}

func Parse() CliConfigurator {
	flag.StringVar(&defaultHost, defaultHostKey, "", "Base address")
	flag.StringVar(&shortHost, shortHostKey, "", "short links host")
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

	return &conf
}
