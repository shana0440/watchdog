package config

import (
	"flag"
)

type IgnoreFlags map[string]struct{}

func (f *IgnoreFlags) String() string {
	return "use to ignore files or dictionary"
}

func (f *IgnoreFlags) Set(value string) error {
	(*f)[value] = struct{}{}
	return nil
}

type App struct {
	Ignores IgnoreFlags
	Command string
}

var appConfig App

func init() {
	flag.Var(&appConfig.Ignores, "ignore", "ignore to watch")
	flag.StringVar(&appConfig.Command, "c", "", "the command when file change will execute")
}

func Parse() App {
	appConfig.Ignores = make(IgnoreFlags)
	appConfig.Ignores[".git"] = struct{}{}
	flag.Parse()
	return appConfig
}
