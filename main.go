package main

import (
	"github.com/shana0440/watchdog/config"
	"github.com/shana0440/watchdog/dog"

	"log"
)

func main() {
	config := config.Parse()
	log.Printf("config: %#v\n", config)

	dog, err := dog.NewDog(".", config.Ignores)
	if err != nil {
		log.Fatalln("Failed to create dog", err)
	}
	defer dog.Close()
	err = dog.Run(config.Command)
	if err != nil {
		log.Fatalln("Failed to run command", err)
	}
}
