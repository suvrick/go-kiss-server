package models

import (
	"strings"
)

// Proxy ...
type Proxy struct {
	ProxyID int
	// host:port
	Host     string
	Username string
	Password string
}

// NewProxy ...
func NewProxy(url string) *Proxy {

	url = strings.Replace(url, "\n", "", -1)
	url = strings.Replace(url, "\r", "", -1)

	prxs := strings.Split(url, ":")

	if len(prxs) != 4 {
		return nil
	}

	return &Proxy{
		Host:     prxs[0] + ":" + prxs[1],
		Username: prxs[2],
		Password: prxs[3],
	}
}
