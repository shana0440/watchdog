package dog

import (
	"context"
	"os"
	"os/exec"
)

type Commander interface {
	Exec(string)
	Errors() chan error
	ReExecute()
}

type Command struct {
	reExecute    chan struct{}
	errors       chan error
	clearConsole bool
}

func NewCommand(clearConsole bool) *Command {
	return &Command{
		reExecute:    make(chan struct{}),
		errors:       make(chan error),
		clearConsole: clearConsole,
	}
}

func (helper *Command) Exec(cmd string) {
	go func(cmd string) {
		for {
			ctx, cancel := context.WithCancel(context.Background())
			command := exec.CommandContext(ctx, "sh", "-c", cmd)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			// sh will not happen error, cmd will, but will return via Wait(), not Start()
			command.Start()
			go func() {
				err := command.Wait()
				if err != nil {
					helper.errors <- err
				}
			}()
			<-helper.reExecute
			cancel()
			clearConsoleIfNeed(helper.clearConsole)
		}
	}(cmd)
}

func clearConsoleIfNeed(clear bool) {
	if clear {
		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()
	}
}

func (helper *Command) Errors() chan error {
	return helper.errors
}

func (helper *Command) ReExecute() {
	helper.reExecute <- struct{}{}
}
