package config

import (
	"flag"
	"io/ioutil"
	"strings"
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
	for _, f := range gitignore() {
		appConfig.Ignores[f] = struct{}{}
	}
	flag.Parse()
	return appConfig
}

func gitignore() []string {
	data, err := ioutil.ReadFile(".gitignore")
	if err != nil {
		return []string{}
	}
	lines := strings.Split(strings.Trim(string(data), "\n"), "\n")
	res := make([]string, 0, len(lines))
	for _, line := range lines {
		if line[0] == '/' {
			res = append(res, line[1:])
		} else {
			res = append(res, line)
		}
	}
	return res
}
