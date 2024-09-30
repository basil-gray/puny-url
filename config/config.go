package config

import (
	"strconv"
	"time"
)

type Config struct {
	EvictInterval time.Duration
	Ttl           time.Duration
	Port          int
}

var Args Config

func BaseURL() string {
	return "http://localhost:" + strconv.Itoa(Args.Port) + "/"
}
