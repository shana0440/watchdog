package helper

import (
	"io/ioutil"
	"strings"
)

func IgnoreGit() []string {
	targets := make([]string, 0)
	targets = append(targets, ".git")

	bytes, err := ioutil.ReadFile(".gitignore")
	if err == nil {
		gitignore := strings.Split(strings.Trim(string(bytes), "\n"), "\n")
		for _, v := range gitignore {
			if strings.HasPrefix(v, "/") {
				targets = append(targets, "."+v)
			} else if !strings.HasPrefix(v, "#") {
				targets = append(targets, v)
			}
		}
	}
	return targets
}
