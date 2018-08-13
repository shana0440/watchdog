package dog

import (
	"github.com/fsnotify/fsnotify"
	"github.com/shana0440/watchdog/config"

	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Dog struct {
	watcher *fsnotify.Watcher
	ignores config.IgnoreFlags
	reExec  chan struct{}
}

func NewDog(dir string, ignores config.IgnoreFlags) (*Dog, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln("Failed to create watcher", err)
	}
	dog := &Dog{w, ignores, make(chan struct{})}
	dirs := dog.getDirRecursivelyWithoutIgnore(dir)
	log.Println("current watch dirs: ", dirs)
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

func (dog *Dog) Run(cmd string) {
	dog.executeCommand(cmd)
	for {
		select {
		case e, ok := <-dog.watcher.Events:
			if !ok {
				log.Fatalln("Watcher is closed")
			}
			if _, ok := dog.ignores[e.Name]; ok {
				continue
			}
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
			dog.reExec <- struct{}{}
		}
	}
}

func (dog *Dog) executeCommand(cmd string) {
	go func() {
		for {
			args := strings.Split(cmd, " ")
			ctx, cancel := context.WithCancel(context.Background())
			command := exec.CommandContext(ctx, args[0], args[1:]...)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			err := command.Start()
			if err != nil {
				log.Println(err)
			}
			<-dog.reExec
			cancel()
		}
	}()
}

func (dog *Dog) Close() {
	dog.watcher.Close()
}
