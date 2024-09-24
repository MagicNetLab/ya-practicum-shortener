package config

import (
	"crypto/rand"
	"encoding/hex"
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
	SetDBConnectString(params string) error
	GetDBConnectString() string
	SetJWTSecret(secret string) error
	GetJWTSecret() string
	SetPProfHost(host string) error
	GetPProfHost() string
	IsValid() bool
}

type configParams struct {
	defaultHost     string
	defaultPort     string
	shortHost       string
	shortPort       string
	fileStoragePath string
	dbConnectString string
	JWTSecret       string
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
	sp.JWTSecret = secret
	return nil
}

func (sp *configParams) GetJWTSecret() string {
	return sp.JWTSecret
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

var servParams configParams

func GetParams() ParameterConfig {
	if servParams.IsValid() {
		return &servParams
	}

	// костыль для тестов. без этих значений приложение в принципе не должно запускаться
	_ = servParams.SetDefaultHost("localhost", "8080")
	_ = servParams.SetShortHost("localhost", "8080")
	_ = servParams.SetFileStoragePath("/tmp/short-url-db.json")
	_ = servParams.SetJWTSecret(getRandomSecret())

	envConf, err := env.Parse()
	if err == nil {
		if envConf.HasBaseHost() {
			host, hostErr := envConf.GetBaseHost()
			port, portErr := envConf.GetBasePort()

			if hostErr == nil && portErr == nil {
				err = servParams.SetDefaultHost(host, port)
				if err != nil {
					args := map[string]interface{}{"error": err.Error()}
					logger.Error("fail set default host from env", args)
				}
			}
		}

		if envConf.HasShortHost() {
			host, hostErr := envConf.GetShortHost()
			port, portErr := envConf.GetShortPort()

			if hostErr == nil && portErr == nil {
				err = servParams.SetShortHost(host, port)
				if err != nil {
					args := map[string]interface{}{"error": err.Error()}
					logger.Error("Fail set short host from env", args)
				}
			}
		}

		if envConf.HasFileStoragePath() {
			storagePath, storageErr := envConf.GetFileStoragePath()
			if storageErr == nil {
				err = servParams.SetFileStoragePath(storagePath)
				if err != nil {
					args := map[string]interface{}{"error": err.Error()}
					logger.Error("Fail set file storage path from env", args)
				}
			}
		}

		if envConf.HasDBConnectString() {
			dbConnectParams, dbParamsErr := envConf.GetDBConnectString()
			if dbParamsErr == nil {
				err = servParams.SetDBConnectString(dbConnectParams)
				if err != nil {
					args := map[string]interface{}{"error": err.Error()}
					logger.Error("Fail set db connect params from env", args)
				}
			}
		}

		if envConf.HasJWTSecret() {
			jwtSecret, jwtSecretErr := envConf.GetJWTSecret()
			if jwtSecretErr == nil {
				err = servParams.SetJWTSecret(jwtSecret)
				if err != nil {
					args := map[string]interface{}{"error": err.Error()}
					logger.Error("Fail set jwttoken secret from env", args)
				}
			}
		}

		err = servParams.SetPProfHost(envConf.GetPPROFHost())
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("Fail set pprof host from env", args)
		}
	}

	cliConf := flags.Parse()

	if cliConf.HasDefaultHost() {
		host, hostErr := cliConf.GetDefaultHost()
		port, portErr := cliConf.GetDefaultPort()
		if hostErr == nil && portErr == nil {
			err = servParams.SetDefaultHost(host, port)
			if err != nil {
				args := map[string]interface{}{"error": err.Error()}
				logger.Error("fail set default host from cli flags", args)
			}
		}
	}

	if cliConf.HasShortHost() {
		host, hostErr := cliConf.GetShortHost()
		port, portErr := cliConf.GetShortPort()
		if hostErr == nil && portErr == nil {
			err = servParams.SetShortHost(host, port)
			if err != nil {
				args := map[string]interface{}{"error": err.Error()}
				logger.Error("Fail set short host from cli flags", args)
			}
		}
	}

	if cliConf.HasFileStoragePath() {
		storagePath, storageErr := cliConf.GetFileStoragePath()
		if storageErr == nil {
			err = servParams.SetFileStoragePath(storagePath)
			if err != nil {
				args := map[string]interface{}{"error": err.Error()}
				logger.Error("Fail set file storage path from cli flags", args)
			}
		}
	}

	if cliConf.HasDBConnectString() {
		dbConnectParams, dbParamsErr := cliConf.GetDBConnectString()
		if dbParamsErr == nil {
			err = servParams.SetDBConnectString(dbConnectParams)
			if err != nil {
				args := map[string]interface{}{"error": err.Error()}
				logger.Error("Fail set db connect params from cli flags", args)
			}
		}
	}

	if cliConf.HasPProfHost() {
		pprofHost, pprofHostErr := cliConf.GetPProfHost()
		if pprofHostErr == nil {
			err = servParams.SetPProfHost(pprofHost)
			if err != nil {
				args := map[string]interface{}{"error": err.Error()}
				logger.Error("Fail set pprof host from cli flags", args)
			}
		}
	}

	return &servParams
}

func getRandomSecret() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}
