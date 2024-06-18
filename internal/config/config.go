package config

import (
	"fmt"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config/env"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config/flags"
	"strings"
)

type ParameterConfig interface {
	SetDefaultHost(host string, port string) error
	SetShortHost(host string, port string) error
	GetDefaultHost() string
	GetShortHost() string
	IsValid() bool
}

// TODO разделить структуру на 2 для defaultHost и shortHost
type params struct {
	defaultHost string
	defaultPort string
	shortHost   string
	shortPort   string
}

func (sp *params) IsValid() bool {
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

func GetParams() ParameterConfig {
	if servParams.IsValid() {
		return &servParams
	}

	// todo default values
	_ = servParams.SetDefaultHost("localhost", "8080")
	_ = servParams.SetShortHost("localhost", "8080")

	envConf, err := env.Parse()
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

	cliConf := flags.Parse()

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
