package main

import (
	"github.com/shana0440/watchdog/config"
	"github.com/shana0440/watchdog/dog"

	"fmt"
	"log"
	"strings"
)

func main() {
	config := config.Parse()

	dir := dog.NewDirectory(".", config.Ignores)
	fmt.Println("Ignores: ", strings.Join(dir.GetIgnoreItem(), ", "))
	cmd := dog.NewCommand()
	dog, err := dog.NewDog(dir, cmd)
	if err != nil {
		log.Fatalln("Failed to create dog", err)
	}
	defer dog.Close()
	err = dog.Run(config.Command)
	if err != nil {
		log.Fatalln("Failed to watch file", err)
	}
}
