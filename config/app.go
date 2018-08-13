package config

import (
	"flag"
)

type ignoreFlags map[string]struct{}

func (f *ignoreFlags) String() string {
	return "use to ignore files or dictionary"
}

func (f *ignoreFlags) Set(value string) error {
	(*f)[value] = struct{}{}
	return nil
}

type App struct {
	Ignores ignoreFlags
	Command string
}

var appConfig App

func init() {
	flag.Var(&appConfig.Ignores, "ignore", "ignore to watch")
	flag.StringVar(&appConfig.Command, "c", "", "the command when file change will execute")
}

func Parse() App {
	appConfig.Ignores = make(ignoreFlags)
	appConfig.Ignores[".git"] = struct{}{}
	flag.Parse()
	return appConfig
}
