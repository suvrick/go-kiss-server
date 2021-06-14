package server

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/BurntSushi/toml"
)

// Config ...
type Config struct {
	BindAddr    string `toml:"bind_addr"`
	DatabaseURL string `toml:"database_url"`
}

//go:embed "config.toml"
var configToml []byte

// NewConfig ...
func NewConfig() *Config {

	config := &Config{
		BindAddr: ":8080",
	}

	if len(configToml) == 0 {
		log.Fatal("[NewConfig] >> error load config.toml file")
	}

	_, err := toml.DecodeReader(bytes.NewReader(configToml), config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
