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
	DataAdd     string
	DataLastUse string
}

// NewProxy ...
func NewProxy(input string) *Proxy {

	input = strings.TrimRight(input, "\r")

	proxy := &Proxy{}
	sub := strings.Split(input, ":")
	if len(sub) != 4 {
		proxy.IsBad = true
		return proxy
	}

	//fmt.Println(sub)

	proxy.Host = fmt.Sprintf("%s:%s", sub[0], sub[1])
	proxy.Login = sub[2]
	proxy.Password = sub[3]
	proxy.DataAdd = time.Now().String()

	return proxy
}
