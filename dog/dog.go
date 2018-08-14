package dog

import (
	"github.com/fsnotify/fsnotify"
	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/iterable"
	"github.com/reactivex/rxgo/observable"
	"github.com/reactivex/rxgo/observer"

	"log"
	"os"
)

type Dog struct {
	Directoryer
	Commander
	watcher *fsnotify.Watcher
}

func NewDog(dir Directoryer, cmd Commander) (*Dog, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	dog := &Dog{
		Directoryer: dir,
		Commander:   cmd,
		watcher:     w,
	}
	dirs, err := dog.GetDirs()
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

func (dog *Dog) Run(cmd string) {
	dog.Exec(cmd)
	// iterable.New only receive chan interface{}, <-chan interface{}, []interface{}
	it, _ := iterable.New(toInterfaceChan(dog.watcher.Events))
	sub := observable.
		From(it).
		Filter(func(item interface{}) bool {
			event := item.(fsnotify.Event)
			return event.Op&fsnotify.Chmod != fsnotify.Chmod
		}).
		Filter(func(item interface{}) bool {
			event := item.(fsnotify.Event)
			return !dog.IsIgnoreFile(event.Name)
		}).
		Subscribe(observer.New(
			handlers.NextFunc(func(item interface{}) {
				event := item.(fsnotify.Event)
				log.Println(event)
				dog.AddWatchWhenCreateDir(event)
				dog.ReExecute()
			}),
		))
	<-sub
}

func (dog *Dog) AddWatchWhenCreateDir(event fsnotify.Event) error {
	if event.Op&fsnotify.Create == fsnotify.Create {
		stat, err := os.Stat(event.Name)
		if err != nil {
			return err
		}
		if stat.IsDir() {
			dog.watcher.Add(event.Name)
		}
	}
	return nil
}

func (dog *Dog) Close() {
	dog.watcher.Close()
}

func toInterfaceChan(ch chan fsnotify.Event) <-chan interface{} {
	channel := make(chan interface{})
	go func() {
		for {
			var data interface{} = <-ch
			channel <- data
		}
	}()
	return channel
}
