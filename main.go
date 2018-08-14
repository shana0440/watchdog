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
	ignores := make([]string, 0, len(config.Ignores))
	for k := range config.Ignores {
		ignores = append(ignores, k)
	}
	fmt.Println("ignore: ", strings.Join(ignores, ", "))

	dog, err := dog.NewDog(".", config.Ignores)
	if err != nil {
		log.Fatalln("Failed to create dog", err)
	}
	defer dog.Close()
	dog.Run(config.Command)
}
