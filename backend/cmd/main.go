package main

import (
	"log"

	"github.com/suvrick/go-kiss-server/server"
)

func main() {

	config := server.NewConfig()

	log.Println("Start server")
	if err := server.Start(config); err != nil {
		log.Fatal(err)
	}
}
