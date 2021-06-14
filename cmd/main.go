package main

import (
	"log"

	"github.com/suvrick/go-kiss-server/server"
)

func main() {

	config := server.NewConfig()
	countStart := 0

	for {
		if countStart == 10 {
			log.Println("Привышен максимальный лимит перезапуска сервера")
			break
		}

		log.Println("Start server")

		var errMsg string
		if err := server.Start(config, &errMsg); err != nil {
			log.Println(errMsg)
		}

		countStart++
	}

}
