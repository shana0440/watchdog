package dog

import (
	"github.com/fsnotify/fsnotify"
)

type Watcher interface {
	Events() chan interface{}
	Add(string)
	Close()
}

type Watch struct {
	*fsnotify.Watcher
}

func NewWatch() *Watch {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	return &Watch{
		Watcher: w,
	}
}

func (w *Watch) Events() chan interface{} {
	channel := make(chan interface{})
	go func() {
		for {
			var data interface{} = <-w.Watcher.Events
			channel <- data
		}
	}()
	return channel
}

func (w *Watch) Add(path string) {
	w.Watcher.Add(path)
}

func (w *Watch) Close() {
	w.Watcher.Close()
}
