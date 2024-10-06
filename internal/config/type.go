package config

import "strings"

type configParams struct {
	defaultHost     string
	defaultPort     string
	shortHost       string
	shortPort       string
	fileStoragePath string
	dbConnectString string
	jwtSecret       string
	pProfOn         bool
	pProfHost       string
}

func (sp *configParams) SetFileStoragePath(path string) error {
	sp.fileStoragePath = path
	return nil
}

func (sp *configParams) GetFileStoragePath() string {
	return sp.fileStoragePath
}

func (sp *configParams) IsValid() bool {
	return sp.defaultHost != "" &&
		sp.defaultPort != "" &&
		sp.shortHost != "" &&
		sp.shortPort != "" &&
		sp.fileStoragePath != ""

}

func (sp *configParams) SetDefaultHost(host string, port string) error {
	sp.defaultHost = host
	sp.defaultPort = port

	return nil
}

func (sp *configParams) SetShortHost(host string, port string) error {
	sp.shortHost = host
	sp.shortPort = port

	return nil
}

func (sp *configParams) GetDefaultHost() string {
	p := []string{sp.defaultHost, sp.defaultPort}
	return strings.Join(p, ":")
}

func (sp *configParams) GetShortHost() string {
	p := []string{sp.shortHost, sp.shortPort}
	return strings.Join(p, ":")
}

func (sp *configParams) SetDBConnectString(params string) error {
	sp.dbConnectString = params

	return nil
}

func (sp *configParams) GetDBConnectString() string {
	return sp.dbConnectString
}

func (sp *configParams) SetJWTSecret(secret string) error {
	sp.jwtSecret = secret
	return nil
}

func (sp *configParams) GetJWTSecret() string {
	return sp.jwtSecret
}

func (sp *configParams) SetPProfOn(isOn bool) error {
	sp.pProfOn = isOn

	return nil
}

func (sp *configParams) IsPProfOn() bool {
	return sp.pProfOn
}

func (sp *configParams) SetPProfHost(host string) error {
	sp.pProfHost = host

	return nil
}

func (sp *configParams) GetPProfHost() string {
	return sp.pProfHost
}
