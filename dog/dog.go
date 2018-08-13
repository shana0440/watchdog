package dog

import (
	"github.com/fsnotify/fsnotify"
	"github.com/shana0440/watchdog/config"

	"errors"
	"log"
	"os"
)

type Dog struct {
	DirHelper
	*CommandHelper
	watcher *fsnotify.Watcher
	reExec  chan struct{}
}

func NewDog(dir string, ignores config.IgnoreFlags) (*Dog, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	dog := &Dog{
		DirHelper:     DirHelper{ignores},
		CommandHelper: NewCommandHelper(),
		watcher:       w,
		reExec:        make(chan struct{}),
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
	dog.Exec(cmd)
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
			dog.ReExecute()
		case err := <-dog.CommandHelper.Error:
			return err
		}
	}
}

func (dog *Dog) Close() {
	dog.watcher.Close()
}
