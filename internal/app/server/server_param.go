package server

import "strings"

const (
	DefaultHostName = "localhost"
	DefaultHostPort = "8080"
	ShortHostName   = "localhost"
	ShortHostPort   = "8080"
)

type configurator interface {
	SetDefaultHost(host string, port string) error
	SetShortHost(host string, port string) error
	GetDefaultHost() string
	GetShortHost() string
}

type ServerParams struct {
	defaultHost string
	defaultPort string
	shortHost   string
	shortPort   string
}

func (sp *ServerParams) Init() error {
	sp.defaultHost = DefaultHostName
	sp.defaultPort = DefaultHostPort
	sp.shortHost = ShortHostName
	sp.shortPort = ShortHostPort

	return nil
}

func (sp *ServerParams) SetDefaultHost(host string, port string) error {
	sp.defaultHost = host
	sp.defaultPort = port

	return nil
}

func (sp *ServerParams) SetShortHost(host string, port string) error {
	sp.shortHost = host
	sp.shortPort = port

	return nil
}

func (sp *ServerParams) GetDefaultHost() string {
	p := []string{sp.defaultHost, sp.defaultPort}
	return strings.Join(p, ":")
}

func (sp *ServerParams) GetShortHost() string {
	p := []string{sp.shortHost, sp.shortPort}
	return strings.Join(p, ":")
}
