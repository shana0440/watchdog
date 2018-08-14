package dog

import (
	"testing"

	"context"
	"os/exec"

	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
)

func TestCommandShouldExecute(t *testing.T) {
	cmd := "echo 1"
	shouldExec := make(chan struct{})
	monkey.Patch(context.WithCancel, func(ctx context.Context) (context.Context, context.CancelFunc) {
		return ctx, func() {
			assert.Fail(t, "shouldn't execute cancel function")
		}
	})

	monkey.Patch(exec.CommandContext, func(ctx context.Context, name string, args ...string) *exec.Cmd {
		assert.Equal(t, args[1], cmd)
		shouldExec <- struct{}{}
		monkey.Unpatch(exec.CommandContext)
		return exec.CommandContext(ctx, name, args...)
	})

	helper := NewCommand(false)
	helper.Exec(cmd)
	<-shouldExec
}

func TestReExecuteShouldReExecuteAndCancelPreviousCommand(t *testing.T) {
	cmd := "echo 1"
	shouldExec := make(chan struct{})
	monkey.Patch(context.WithCancel, func(ctx context.Context) (context.Context, context.CancelFunc) {
		return ctx, func() {
			shouldExec <- struct{}{}
		}
	})

	monkey.Patch(exec.CommandContext, func(ctx context.Context, name string, args ...string) *exec.Cmd {
		assert.Equal(t, args[1], cmd)
		shouldExec <- struct{}{}
		monkey.Unpatch(exec.CommandContext)
		return exec.CommandContext(ctx, name, args...)
	})

	helper := NewCommand(false)
	helper.Exec(cmd)
	<-shouldExec
	helper.ReExecute()
	<-shouldExec
}
