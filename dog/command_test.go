package dog

import (
	"testing"

	"os"
	"os/exec"
	"reflect"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestExecShould(t *testing.T) {
	t.Run("ExecuteCommandAndBlocking", func(t *testing.T) {
		defer monkey.UnpatchAll()

		executed := false
		done := make(chan struct{})
		monkey.Patch(exec.Command, func(name string, args ...string) *exec.Cmd {
			// name is "sh", args[0] is "-c", then is our command
			monkey.Unpatch(exec.Command)
			assert.Equal(t, args[1], "true")
			executed = true
			done <- struct{}{}
			return exec.Command(name, args...)
		})

		cmd := NewCommand(false)
		cmd.Exec("true")
		<-done
		assert.True(t, executed, "no execute command")
	})

	t.Run("ReExecuteCommandWhenCallReExceute", func(t *testing.T) {
		defer monkey.UnpatchAll()

		killed := false
		done := make(chan struct{})
		monkey.Patch(exec.Command, func(name string, args ...string) *exec.Cmd {
			monkey.Unpatch(exec.Command)
			cmd := exec.Command(name, args...)
			monkey.PatchInstanceMethod(reflect.TypeOf(cmd.Process), "Kill", func(p *os.Process) error {
				killed = true
				done <- struct{}{}
				return nil
			})
			return cmd
		})

		cmd := NewCommand(false)
		cmd.Exec("true")
		cmd.ReExecute()
		<-done
		assert.True(t, killed, "no execute process kill function")
	})

	t.Run("ReceiveErrorWhenCommandWrong", func(t *testing.T) {
		defer monkey.UnpatchAll()

		done := make(chan struct{})
		monkey.Patch(exec.Command, func(name string, args ...string) *exec.Cmd {
			monkey.Unpatch(exec.Command)
			done <- struct{}{}
			return exec.Command(name, args...)
		})

		helper := NewCommand(false)
		helper.Exec("invalid")
		<-done
		err := <-helper.Errors()
		assert.Error(t, err)
	})
}

func TestClearConsoleIfNeedShould(t *testing.T) {
	t.Run("ClearConsoleWhenTrue", func(t *testing.T) {
		c := exec.Command("clear")
		monkey.Patch(exec.Command, func(cmd string, _ ...string) *exec.Cmd {
			assert.Equal(t, "clear", cmd)
			return c
		})
		done := false
		monkey.PatchInstanceMethod(reflect.TypeOf(c), "Run", func(cmd *exec.Cmd) error {
			done = true
			return nil
		})
		defer monkey.UnpatchAll()
		clearConsoleIfNeed(true)
		assert.True(t, done, "no execute clear command")
	})

	t.Run("NotClearConsoleWhenFalse", func(t *testing.T) {
		c := exec.Command("clear")
		monkey.Patch(exec.Command, func(cmd string, _ ...string) *exec.Cmd {
			assert.Equal(t, "clear", cmd)
			return c
		})
		done := false
		monkey.PatchInstanceMethod(reflect.TypeOf(c), "Run", func(cmd *exec.Cmd) error {
			done = true
			return nil
		})
		defer monkey.UnpatchAll()
		clearConsoleIfNeed(false)
		assert.False(t, done, "should not execute clear command")
	})
}
