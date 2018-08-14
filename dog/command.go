package dog

import (
	"context"
	"os"
	"os/exec"
)

type Commander interface {
	Exec(string)
	ReExecute()
}

type Command struct {
	reExecute chan struct{}
	Error     chan error
}

func NewCommand() *Command {
	return &Command{
		reExecute: make(chan struct{}),
		Error:     make(chan error),
	}
}

func (helper *Command) Exec(cmd string) {
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
			command.Wait()
		}
	}(cmd)
}

func (helper *Command) ReExecute() {
	helper.reExecute <- struct{}{}
}
