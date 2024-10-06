package flags

import "errors"

const (
	defaultHostKey  = "a"
	shortHostKey    = "b"
	fileStoragePath = "f"
	dbConnectKey    = "d"
	jwtSecret       = "j"
	pProfKey        = "p"
)

type cliConf struct {
	defaultHost     string
	defaultPort     string
	shortHost       string
	shortPort       string
	fileStoragePath string
	dbConnectString string
	jwtSecret       string
	pProfHost       string
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

func (cc *cliConf) HasDBConnectString() bool {
	return cc.dbConnectString != ""
}

func (cc *cliConf) GetDBConnectString() (string, error) {
	if !cc.HasDBConnectString() {
		return "", errors.New("database connect param not set")
	}

	return cc.dbConnectString, nil
}

func (cc *cliConf) HasJWTSecret() bool {
	return cc.jwtSecret != ""
}

func (cc *cliConf) GetJWTSecret() (string, error) {
	if cc.jwtSecret == "" {
		return "", errors.New("jwttoken secret not set")
	}
	return cc.jwtSecret, nil
}

func (cc *cliConf) HasPProfHost() bool {
	return cc.pProfHost != ""
}

func (cc *cliConf) GetPProfHost() string {
	return cc.pProfHost
}
