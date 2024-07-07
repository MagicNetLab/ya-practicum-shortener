package config

import (
	"strings"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config/env"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config/flags"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

type ParameterConfig interface {
	SetDefaultHost(host string, port string) error
	SetShortHost(host string, port string) error
	GetDefaultHost() string
	GetShortHost() string
	SetFileStoragePath(path string) error
	GetFileStoragePath() string
	IsValid() bool
}

// TODO разделить структуру на 2 для defaultHost и shortHost
type configParams struct {
	defaultHost     string
	defaultPort     string
	shortHost       string
	shortPort       string
	fileStoragePath string
}

func (sp *configParams) SetFileStoragePath(path string) error {
	sp.fileStoragePath = path
	return nil
}

func (sp *configParams) GetFileStoragePath() string {
	return sp.fileStoragePath
}

func (sp *configParams) IsValid() bool {
	return sp.defaultHost != "" && sp.defaultPort != "" && sp.shortHost != "" && sp.shortPort != "" && sp.fileStoragePath != ""
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

var servParams configParams

func GetParams() ParameterConfig {
	if servParams.IsValid() {
		return &servParams
	}

	// todo default values
	_ = servParams.SetDefaultHost("localhost", "8080")
	_ = servParams.SetShortHost("localhost", "8080")
	_ = servParams.SetFileStoragePath("/tmp/short-url-db.json")

	envConf, err := env.Parse()
	if err == nil {
		if envConf.HasBaseHost() {
			host, hostErr := envConf.GetBaseHost()
			port, portErr := envConf.GetBasePort()

			if hostErr == nil && portErr == nil {
				err = servParams.SetDefaultHost(host, port)
				if err != nil {
					logger.Log.Errorf("Fail set default host from env: %s", err)
				}
			}
		}

		if envConf.HasShortHost() {
			host, hostErr := envConf.GetShortHost()
			port, portErr := envConf.GetShortPort()

			if hostErr == nil && portErr == nil {
				err = servParams.SetShortHost(host, port)
				if err != nil {
					logger.Log.Errorf("Fail set short host from env: %s", err)
				}
			}
		}

		if envConf.HasFileStoragePath() {
			storagePath, storageErr := envConf.GetFileStoragePath()
			if storageErr == nil {
				err = servParams.SetFileStoragePath(storagePath)
				if err != nil {
					logger.Log.Errorf("Fail set file storage path from env: %s", err)
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
				logger.Log.Errorf("Fail set default host from cli flags: %s", err)
			}
		}
	}

	if cliConf.HasShortHost() {
		host, hostErr := cliConf.GetShortHost()
		port, portErr := cliConf.GetShortPort()
		if hostErr == nil && portErr == nil {
			err = servParams.SetShortHost(host, port)
			if err != nil {
				logger.Log.Errorf("Fail set short host from cli flags: %s", err)
			}
		}
	}

	if cliConf.HasFileStoragePath() {
		storagePath, storageErr := cliConf.GetFileStoragePath()
		if storageErr == nil {
			err = servParams.SetFileStoragePath(storagePath)
			if err != nil {
				logger.Log.Errorf("Fail set file storage path from cli flags: %s", err)
			}
		}
	}

	return &servParams
}
