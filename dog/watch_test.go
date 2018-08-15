package dog

import (
	"testing"

	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
)

func TestEventsShouldConvertEventToInterface(t *testing.T) {
	w := NewWatch()
	w.Add(".")
	defer w.Close()
	done := make(chan struct{})

	go func() {
		event := <-w.Events()
		switch event.(type) {
		case fsnotify.Event:
			done <- struct{}{}
		default:
			assert.Fail(t, "type should be fsnotify.Event")
		}
	}()

	os.MkdirAll("tmp", 0777)
	defer os.RemoveAll("tmp")
	<-done
}
