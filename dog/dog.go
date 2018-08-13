package dog

import (
	"github.com/fsnotify/fsnotify"
	"github.com/shana0440/watchdog/config"

	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Dog struct {
	watcher *fsnotify.Watcher
	ignores config.IgnoreFlags
}

func NewDog(dir string, ignores config.IgnoreFlags) (*Dog, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln("Failed to create watcher", err)
	}
	dog := &Dog{w, ignores}
	dirs := dog.getDirRecursivelyWithoutIgnore(dir)
	log.Println(dirs)
	dog.watchDirs(dirs)
	return dog, nil
}

func (dog *Dog) getDirRecursivelyWithoutIgnore(dir string) []string {
	dirs := []string{dir}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln("Failed to read dir", err)
	}
	for _, f := range files {
		if _, ok := dog.ignores[f.Name()]; f.IsDir() && !ok {
			newDir := fmt.Sprintf("%s/%s", dir, f.Name())
			dirs = append(dirs, dog.getDirRecursivelyWithoutIgnore(newDir)...)
		}
	}
	return dirs
}

func (dog *Dog) watchDirs(dirs []string) {
	for _, dir := range dirs {
		dog.watcher.Add(dir)
	}
}

func (dog *Dog) Run() {
	for {
		select {
		case e, ok := <-dog.watcher.Events:
			if !ok {
				log.Fatalln("Watcher is closed")
			}
			if _, ok := dog.ignores[e.Name]; !ok {
				log.Println(e)
				if e.Op&fsnotify.Create == fsnotify.Create {
					stat, err := os.Stat(e.Name)
					if err != nil {
						log.Println("Failed to get file state", err)
						continue
					}
					if stat.IsDir() {
						dog.watcher.Add(e.Name)
					}
				}
			}
		}
	}
}

func (dog *Dog) Close() {
	dog.watcher.Close()
}
