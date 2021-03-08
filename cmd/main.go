package main

import (
	"log"

	"github.com/suvrick/go-kiss-server/server"

	"github.com/BurntSushi/toml"
)

func main() {

	config := server.NewConfig()
	_, err := toml.DecodeFile("configs/apiserver.toml", config)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Start server")
	if err := server.Start(config); err != nil {
		log.Fatal(err)
	}
}
