package main

import (
	"github.com/shana0440/watchdog/config"
	"github.com/shana0440/watchdog/dog"

	"log"
)

func main() {
	config := config.Parse()
	log.Printf("config: %#v\n", config)

	dog, _ := dog.NewDog(".", config.Ignores)
	defer dog.Close()
	dog.Run(config.Command)
}
