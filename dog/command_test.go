package dog

import (
	"testing"

	"context"
	"os/exec"
	"reflect"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestExecShould(t *testing.T) {
	t.Run("ExecuteCommandAndBlocking", func(t *testing.T) {
		monkey.Patch(context.WithCancel, func(ctx context.Context) (context.Context, context.CancelFunc) {
			return ctx, func() {
				assert.Fail(t, "shouldn't execute cancel function")
			}
		})

		executed := false
		done := make(chan struct{})
		monkey.Patch(exec.CommandContext, func(ctx context.Context, name string, args ...string) *exec.Cmd {
			// name is "sh", args[0] is "-c", then is our command
			assert.Equal(t, args[1], "true")
			executed = true
			done <- struct{}{}
			monkey.Unpatch(exec.CommandContext)
			return exec.CommandContext(ctx, name, args...)
		})
		defer monkey.UnpatchAll()

		helper := NewCommand(false)
		helper.Exec("true")
		<-done
		assert.True(t, executed, "no execute command")
	})

	t.Run("ReExecuteCommandWhenCallReExceute", func(t *testing.T) {
		cancel := false
		done := make(chan struct{})
		monkey.Patch(context.WithCancel, func(ctx context.Context) (context.Context, context.CancelFunc) {
			return ctx, func() {
				cancel = true
				done <- struct{}{}
			}
		})
		defer monkey.UnpatchAll()

		helper := NewCommand(false)
		helper.Exec("true")
		helper.ReExecute()
		<-done
		assert.True(t, cancel, "no execute cancel function")
	})

	t.Run("ReceiveErrorWhenCommandWrong", func(t *testing.T) {
		done := make(chan struct{})
		monkey.Patch(exec.CommandContext, func(ctx context.Context, name string, args ...string) *exec.Cmd {
			done <- struct{}{}
			monkey.Unpatch(exec.CommandContext)
			return exec.CommandContext(ctx, name, args...)
		})
		defer monkey.UnpatchAll()

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
