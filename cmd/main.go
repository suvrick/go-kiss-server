package main

import (
	"fmt"
	"log"
	"os"

	"github.com/suvrick/go-kiss-server/server"
)

func main() {

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("workdir: %s", dir)

	config := server.NewConfig()

	log.Println("Start server")
	if err := server.Start(config); err != nil {
		log.Fatal(err)
	}

}
