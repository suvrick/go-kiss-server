package model

import (
	"fmt"
	"strings"
	"time"
)

// Proxy ..
type Proxy struct {
	ID          int
	Host        string
	Login       string
	Password    string
	IsBusy      bool
	IsBad       bool
	DateAdd     string
	DateLastUse string
	URL         string
}

// NewProxy ...
func NewProxy(input string) *Proxy {

	input = strings.TrimRight(input, "\r")
	input = strings.TrimRight(input, "\n")

	proxy := &Proxy{
		URL: input,
	}

	sub := strings.Split(input, ":")
	if len(sub) != 4 {
		proxy.IsBad = true
		return proxy
	}

	//fmt.Println(sub)

	proxy.Host = fmt.Sprintf("%s:%s", sub[0], sub[1])
	proxy.Login = sub[2]
	proxy.Password = sub[3]
	proxy.DateAdd = time.Now().Format("2006-01-02")

	return proxy
}
