package http

import (
	"time"
)

type Config struct {
	Host                  string
	Port                  uint32
	Timeout               time.Duration
	KeepAlive             time.Duration
	IdleConnTimeout       time.Duration
	TLSHandshakeTimeout   time.Duration
	ExpectContinueTimeout time.Duration
	MaxIdleConnsPerHost   int
}

func NewDefaultConfig() *Config {
	return &Config{
		Host:                  "localhost",
		Port:                  8080,
		MaxIdleConnsPerHost:   200,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		Timeout:               30 * time.Second,
		KeepAlive:             30 * time.Second,
	}
}
