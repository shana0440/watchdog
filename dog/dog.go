package dog

import (
	"github.com/fsnotify/fsnotify"
	"github.com/shana0440/watchdog/config"

	"context"
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Dog struct {
	DirHelper
	watcher *fsnotify.Watcher
	reExec  chan struct{}
}

func NewDog(dir string, ignores config.IgnoreFlags) (*Dog, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	dog := &Dog{
		DirHelper: DirHelper{ignores},
		watcher:   w,
		reExec:    make(chan struct{}),
	}
	dirs, err := dog.GetDirs(dir)
	if err != nil {
		return nil, err
	}
	dog.watchDirs(dirs)
	return dog, nil
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
			if dog.IsIgnoreFile(e.Name) {
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
