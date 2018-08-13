package dog

import (
	"github.com/fsnotify/fsnotify"
	"github.com/shana0440/watchdog/config"

	"context"
	"errors"
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
		return nil, err
	}
	dog := &Dog{w, ignores, make(chan struct{})}
	dirs, err := dog.getDirRecursivelyWithoutIgnore(dir)
	if err != nil {
		return nil, err
	}
	dog.watchDirs(dirs)
	return dog, nil
}

func (dog *Dog) getDirRecursivelyWithoutIgnore(dir string) ([]string, error) {
	dirs := []string{dir}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if _, ok := dog.ignores[f.Name()]; f.IsDir() && !ok {
			newDir := fmt.Sprintf("%s/%s", dir, f.Name())
			dir, err := dog.getDirRecursivelyWithoutIgnore(newDir)
			if err != nil {
				return nil, err
			}
			dirs = append(dirs, dir...)
		}
	}
	return dirs, nil
}

func (dog *Dog) watchDirs(dirs []string) {
	for _, dir := range dirs {
		dog.watcher.Add(dir)
	}
}

func (dog *Dog) Run(cmd string) error {
	dog.executeCommand(cmd)
	for {
		select {
		case e, ok := <-dog.watcher.Events:
			if !ok {
				return errors.New("Watcher is closed")
			}
			if _, ok := dog.ignores[e.Name]; ok {
				continue
			}
			log.Println(e)
			if e.Op&fsnotify.Create == fsnotify.Create {
				stat, err := os.Stat(e.Name)
				if err == nil && stat.IsDir() {
					dog.watcher.Add(e.Name)
				}
			}
			dog.reExec <- struct{}{}
		}
	}
}

func (dog *Dog) executeCommand(cmd string) {
	go func(cmd string) {
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
	}(cmd)
}

func (dog *Dog) Close() {
	dog.watcher.Close()
}
