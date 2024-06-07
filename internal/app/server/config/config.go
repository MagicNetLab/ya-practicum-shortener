package config

import (
	"fmt"
	"strings"
)

const (
	defaultHostName = "localhost"
	defaultHostPort = "8080"
	shortHostName   = "localhost"
	shortHostPort   = "8080"
)

type ParamsInterface interface {
	SetDefaultHost(host string, port string) error
	SetShortHost(host string, port string) error
	GetDefaultHost() string
	GetShortHost() string
}

type params struct {
	defaultHost string
	defaultPort string
	shortHost   string
	shortPort   string
}

func (sp *params) initParams() error {
	sp.defaultHost = defaultHostName
	sp.defaultPort = defaultHostPort
	sp.shortHost = shortHostName
	sp.shortPort = shortHostPort

	return nil
}

func (sp *params) IsInit() bool {
	return sp.defaultHost != "" && sp.defaultPort != "" && sp.shortHost != "" && sp.shortPort != ""
}

func (sp *params) SetDefaultHost(host string, port string) error {
	sp.defaultHost = host
	sp.defaultPort = port

	return nil
}

func (sp *params) SetShortHost(host string, port string) error {
	sp.shortHost = host
	sp.shortPort = port

	return nil
}

func (sp *params) GetDefaultHost() string {
	p := []string{sp.defaultHost, sp.defaultPort}
	return strings.Join(p, ":")
}

func (sp *params) GetShortHost() string {
	p := []string{sp.shortHost, sp.shortPort}
	return strings.Join(p, ":")
}

var servParams params

func GetParams() ParamsInterface {
	if servParams.IsInit() {
		return &servParams
	}

	err := servParams.initParams()
	if err != nil {
		panic(err)
	}

	envConf, err = ReadEnv()
	if err == nil {
		if envConf.HasBaseHost() {
			host, hostErr := envConf.GetBaseHost()
			port, portErr := envConf.GetBasePort()

			if hostErr == nil && portErr == nil {
				err = servParams.SetDefaultHost(host, port)
				if err != nil {
					fmt.Println("Fail set default host from env")
				}
			}
		}

		if envConf.HasShortHost() {
			host, hostErr := envConf.GetShortHost()
			port, portErr := envConf.GetShortPort()

			if hostErr == nil && portErr == nil {
				err = servParams.SetShortHost(host, port)
				if err != nil {
					fmt.Println("Fail set short host from env")
				}
			}
		}
	}

	cliConf := ParseInitFlag()

	if cliConf.HasDefaultHost() {
		host, hostErr := cliConf.GetDefaultHost()
		port, portErr := cliConf.GetDefaultPort()
		if hostErr == nil && portErr == nil {
			err = servParams.SetDefaultHost(host, port)
			if err != nil {
				fmt.Println("Fail set default host from cli flags")
			}

		}
	}

	if cliConf.HasShortHost() {
		host, hostErr := cliConf.GetShortHost()
		port, portErr := cliConf.GetShortPort()
		if hostErr == nil && portErr == nil {
			err = servParams.SetShortHost(host, port)
			if err != nil {
				fmt.Println("Fail set short host from cli flags")
			}
		}
	}

	return &servParams
}
