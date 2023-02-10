package main

import (
	"log"
	"ordent/internal/app"
	"ordent/internal/config"
)

func main() {
	log.Print("start program")

	config := getConfig()

	httpServer := app.InitHTTPServer(config)

	httpServer.ListenAndServe()
}

func getConfig() config.Config {
	config, err := config.Init()
	if err != nil {
		panic("Failed to read config: " + err.Error())
	}

	return *config
}
