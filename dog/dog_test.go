package dog

import (
	"testing"

	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/mock"
)

type MockDirectory struct {
	mock.Mock
}

func (m *MockDirectory) GetDirs() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockDirectory) ShouldIgnore(file string) bool {
	args := m.Called(file)
	return args.Bool(0)
}

type MockCommand struct {
	mock.Mock
}

func (m *MockCommand) Exec(cmd string) {
	m.Called(cmd)
}

func (m *MockCommand) ReExecute() {
	m.Called()
}

func (m *MockCommand) Errors() chan error {
	return make(chan error)
}

type MockWatch struct {
	mock.Mock
}

func (m *MockWatch) Events() chan interface{} {
	args := m.Called()
	return args.Get(0).(chan interface{})
}

func (m *MockWatch) Add(path string) {
	m.Called(path)
}

func (m *MockWatch) Close() {
}

func TestRunShould(t *testing.T) {
	t.Run("ExecuteCommandAtFirst", func(t *testing.T) {
		dir, cmd, watch := new(MockDirectory), new(MockCommand), new(MockWatch)
		dog := NewDog(dir, cmd, watch)
		dir.On("GetDirs").Return([]string{"tmp", "dir"})
		cmd.On("Exec", "echo 1").Return()
		watch.On("Add", "tmp").Return()
		watch.On("Add", "dir").Return()
		ch := make(chan interface{})
		watch.On("Events").Return(ch)

		done := make(chan struct{})
		go func() {
			dog.Run("echo 1")
			done <- struct{}{}
		}()
		close(ch)
		<-done

		cmd.AssertCalled(t, "Exec", "echo 1")
		cmd.AssertNotCalled(t, "ReExecute")
	})

	t.Run("ReExecuteCommandWhenReceiveFsNotify", func(t *testing.T) {
		dir, cmd, watch := new(MockDirectory), new(MockCommand), new(MockWatch)
		dog := NewDog(dir, cmd, watch)
		dir.On("GetDirs").Return([]string{"tmp", "dir"})
		dir.On("ShouldIgnore", "hello").Return(false)
		cmd.On("Exec", "echo 1").Return()
		cmd.On("ReExecute").Return()
		watch.On("Add", "tmp").Return()
		watch.On("Add", "dir").Return()
		ch := make(chan interface{})
		watch.On("Events").Return(ch)

		done := make(chan struct{})
		go func() {
			dog.Run("echo 1")
			done <- struct{}{}
		}()
		ch <- fsnotify.Event{Name: "hello", Op: fsnotify.Write}
		close(ch)
		<-done

		cmd.AssertCalled(t, "Exec", "echo 1")
		cmd.AssertCalled(t, "ReExecute")
	})
}

func TestAddWatchWhenCreateDirShouldAddDirWhenDirCreated(t *testing.T) {
	dir, cmd, watch := new(MockDirectory), new(MockCommand), new(MockWatch)
	dog := NewDog(dir, cmd, watch)
	dir.On("GetDirs").Return([]string{})
	watch.On("Add", "tmp").Return()

	os.Mkdir("tmp", 0777)
	defer os.RemoveAll("tmp")
	dog.addWatchWhenCreateDir(fsnotify.Event{Name: "tmp", Op: fsnotify.Create})

	watch.AssertCalled(t, "Add", "tmp")
}
