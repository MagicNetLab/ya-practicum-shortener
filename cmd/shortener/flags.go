package main

import (
	"flag"
	"strings"
)

const (
	DefaultHost    = "localhost"
	DefaultPort    = "8080"
	ShortHost      = "localhost"
	ShortPort      = "8081"
	DefaultAddress = "localhost:8080"
	ShortAddress   = "localhost:8081"
)

var baseHost string
var shortHost string

type AppFlagsStruct struct {
	baseHost  string
	basePort  string
	shortHost string
	shortPort string
}

var AppFlags AppFlagsStruct

func (f *AppFlagsStruct) GetBaseAddress() string {
	host := DefaultHost
	if f.baseHost != "" {
		host = f.baseHost
	}

	port := DefaultPort
	if f.basePort != "" {
		port = f.basePort
	}
	elems := []string{host, port}
	return strings.Join(elems, ":")
}

func (f *AppFlagsStruct) GetShortAddress() string {
	host := ShortHost
	if f.shortHost != "" {
		host = f.shortHost
	}

	port := ShortPort
	if f.shortPort != "" {
		port = f.shortPort
	}
	elems := []string{host, port}
	return strings.Join(elems, ":")
}

func ParseInitFlag() {
	flag.StringVar(&baseHost, "a", DefaultAddress, "Base address")
	flag.StringVar(&shortHost, "b", ShortAddress, "short links host")

	flag.Parse()

	bh := strings.Split(baseHost, ":")
	if len(bh) != 2 {
		bh = []string{"localhost", "8080"}
	}

	sh := strings.Split(shortHost, ":")
	if len(sh) != 2 {
		sh = []string{"localhost", "8081"}
	}

	AppFlags.baseHost = bh[0]
	AppFlags.basePort = bh[1]
	AppFlags.shortHost = sh[0]
	AppFlags.shortPort = sh[1]
}
