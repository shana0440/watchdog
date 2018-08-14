package dog

import (
	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"

	"testing"

	"context"
	"os/exec"
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

	helper := NewCommandHelper()
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

	helper := NewCommandHelper()
	helper.Exec(cmd)
	<-shouldExec
	helper.ReExecute()
	<-shouldExec
}
