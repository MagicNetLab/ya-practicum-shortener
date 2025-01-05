package flagsreader

import (
	"flag"
	"strings"
)

// Parse сбор параметров установленных в cli при запуске приложения.
func Parse() CliConf {
	var conf CliConf
	var defaultHost string
	var shortHost string
	var fileStorage string
	var dbConnectString string
	var jwtSecretKey string
	var pProfHost string
	var enableHTTPS string
	var configFilePath string
	var trustedSubnetAddr string
	var grpcPort string

	flag.StringVar(&defaultHost, defaultHostKey, "", "Base address")
	flag.StringVar(&shortHost, shortHostKey, "", "short links host")
	flag.StringVar(&fileStorage, fileStoragePath, "", "file storage path")
	flag.StringVar(&dbConnectString, dbConnectKey, "", "database connect param")
	flag.StringVar(&jwtSecretKey, jwtSecret, "", "jwttoken secret")
	flag.StringVar(&pProfHost, pProfKey, "", "pprof host")
	flag.StringVar(&enableHTTPS, enableHTTPSKey, "", "enable https")
	flag.StringVar(&configFilePath, configFileKey, "", "config file path")
	flag.StringVar(&trustedSubnetAddr, trustedSubnet, "", "trusted subnet address")
	flag.StringVar(&grpcPort, grpcPortKey, "", "grpc port")
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
	conf.dbConnectString = dbConnectString
	conf.jwtSecret = jwtSecretKey
	conf.pProfHost = pProfHost
	if enableHTTPS != "" {
		conf.hasEnableHTTPS = true
		conf.enableHTTPS = enableHTTPS == "true"
	}
	conf.configFilePath = configFilePath
	conf.trustedSubnet = trustedSubnetAddr
	conf.grpcPort = grpcPort

	return conf
}
