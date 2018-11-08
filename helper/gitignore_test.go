package helper_test

import (
	"testing"

	"io/ioutil"
	"os"
	"strings"

	"github.com/shana0440/watchdog/helper"
)

func TestIgnoreGetShould(t *testing.T) {
	t.Run("IgnoreComment", func(t *testing.T) {
		cover := createGitIgnoreFile("# should be ignore")
		defer cover()
		ignores := helper.IgnoreGit()
		for _, v := range ignores {
			if strings.Contains(v, "should be ignore") {
				t.Error("comment should be ignore!!")
			}
		}
	})

	t.Run("ConvertAbsToRelative", func(t *testing.T) {
		cover := createGitIgnoreFile("/node_modules")
		defer cover()
		ignores := helper.IgnoreGit()
		for _, v := range ignores {
			if strings.HasPrefix(v, "/") {
				t.Error("/ should be convert to ./ !!")
			}
		}
	})

	t.Run("IgnoreFileOrDirectory", func(t *testing.T) {
		cover := createGitIgnoreFile("node_modules")
		defer cover()
		ignores := helper.IgnoreGit()
		if ignores[1] != "node_modules" {
			t.Error("should ignore node_modules")
		}
	})
}

func createGitIgnoreFile(content string) func() {
	os.Rename(".gitignore", ".gitignore.bak")
	ioutil.WriteFile(".gitignore", []byte(content), 0644)
	return func() {
		os.Remove(".gitignore")
		os.Rename(".gitignore.bak", ".gitignore")
	}
}
