package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/shana0440/watchdog/config"

	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	config := config.Parse()
	log.Println(config)
	dirs := getAllWatchDir(".")
	log.Println(dirs)

	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln("Failed to create watcher", err)
	}

	if err != nil {
		log.Fatalln("Failed to add folder", err)
	}

	for _, d := range dirs {
		w.Add(d)
	}
	for {
		select {
		case e, ok := <-w.Events:
			if !ok {
				log.Fatalln("Watcher is closed")
			}
			if _, ok := config.Ignores[e.Name]; !ok {
				if strings.Contains(e.Name, ".git") {
					continue
				}
				log.Println(e)
				if e.Op&fsnotify.Create == fsnotify.Create {
					stat, err := os.Stat(e.Name)
					if err != nil {
						log.Fatalln("failed to get file state", err)
					}
					if stat.IsDir() {
						w.Add(e.Name)
					}
				}
			}
		}
	}

	defer w.Close()
}

func getAllWatchDir(dir string) []string {
	log.Println(dir)
	dirs := []string{dir}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln("Failed to read dir", err)
	}
	for _, f := range files {
		if f.IsDir() {
			newDir := fmt.Sprintf("%s/%s", dir, f.Name())
			dirs = append(dirs, getAllWatchDir(newDir)...)
		}
	}
	return dirs
}
