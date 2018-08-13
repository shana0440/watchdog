package dog

import (
	"context"
	"os"
	"os/exec"
)

type CommandHelper struct {
	reExecute chan struct{}
	Error     chan error
}

func NewCommandHelper() *CommandHelper {
	return &CommandHelper{
		reExecute: make(chan struct{}),
		Error:     make(chan error),
	}
}

func (helper *CommandHelper) ReExecute() {
	helper.reExecute <- struct{}{}
}

func (helper *CommandHelper) Exec(cmd string) {
	go func(cmd string) {
		for {
			ctx, cancel := context.WithCancel(context.Background())
			command := exec.CommandContext(ctx, "sh", "-c", cmd)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			err := command.Start()
			if err != nil {
				helper.Error <- err
			}
			<-helper.reExecute
			cancel()
		}
	}(cmd)
}
