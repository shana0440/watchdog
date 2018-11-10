package dog

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/iterable"
	"github.com/reactivex/rxgo/observable"
	"github.com/reactivex/rxgo/observer"
)

type Dog struct {
	Directoryer
	Commander
	Watcher
}

func NewDog(dir Directoryer, cmd Commander, watcher Watcher) *Dog {
	dog := Dog{
		Directoryer: dir,
		Commander:   cmd,
		Watcher:     watcher,
	}
	return &dog
}

func (dog *Dog) watch() {
	dirs := dog.GetDirs()
	for _, dir := range dirs {
		dog.Add(dir)
	}
}

func (dog *Dog) Run(cmd string) error {
	dog.watch()
	dog.Exec(cmd)
	// iterable.New only receive chan interface{}, <-chan interface{}, []interface{}
	it, _ := iterable.New(dog.EventsAndErrorsStream())
	sub := observable.
		From(it).
		Filter(func(item interface{}) bool {
			switch item.(type) {
			case fsnotify.Event:
				return true
			case error:
				log.Println(item)
				return false
			}
			return false
		}).
		Filter(func(item interface{}) bool {
			event := item.(fsnotify.Event)
			return !dog.ShouldIgnore(event.Name)
		}).
		Subscribe(observer.New(
			handlers.NextFunc(func(item interface{}) {
				event := item.(fsnotify.Event)
				log.Println(event)
				dog.addWatchWhenCreateDir(event)
				dog.ReExecute()
			}),
		))
	<-sub
	return nil
}

func (dog *Dog) EventsAndErrorsStream() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		event, err := false, false
		for {
			select {
			case data, ok := <-dog.Events():
				if ok {
					ch <- data
				} else {
					event = true
				}

			case data, ok := <-dog.Errors():
				if ok {
					ch <- data
				} else {
					err = true
				}
			}

			if event && err {
				close(ch)
				break
			}
		}
	}()
	return ch
}

func (dog *Dog) addWatchWhenCreateDir(event fsnotify.Event) error {
	if event.Op&fsnotify.Create == fsnotify.Create {
		stat, err := os.Stat(event.Name)
		if err != nil {
			return err
		}
		if stat.IsDir() {
			dog.Add(event.Name)
		}
	}
	return nil
}

func (dog *Dog) Close() {
	dog.Close()
}
