package helper_test

import (
	"testing"

	"io/ioutil"
	"os"
	"strings"

	"github.com/shana0440/watchdog/helper"
)

func TestIgnoreGetShouldIgnoreComment(t *testing.T) {
	t.Run("IgnoreComment", func(t *testing.T) {
		os.Rename(".gitignore", ".gitignore.bak")
		defer func() {
			os.Rename(".gitignore.bak", ".gitignore")
		}()
		ioutil.WriteFile(".gitignore", []byte("# should be ignore"), 0644)
		defer func() {
			os.Remove(".gitignore")
		}()

		ignores := helper.IgnoreGit()
		for _, v := range ignores {
			if strings.Contains(v, "should be ignore") {
				t.Error("comment should be ignore!!")
			}
		}
	})

	t.Run("ConvertAbsToRelative", func(t *testing.T) {
		os.Rename(".gitignore", ".gitignore.bak")
		defer func() {
			os.Rename(".gitignore.bak", ".gitignore")
		}()
		ioutil.WriteFile(".gitignore", []byte("/node_modules"), 0644)
		defer func() {
			os.Remove(".gitignore")
		}()

		ignores := helper.IgnoreGit()
		for _, v := range ignores {
			if strings.HasPrefix(v, "/") {
				t.Error("/ should be convert to ./ !!")
			}
		}
	})
}
