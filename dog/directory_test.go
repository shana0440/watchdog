package dog

import (
	"github.com/shana0440/watchdog/config"
	"github.com/stretchr/testify/assert"

	"testing"

	"io/ioutil"
	"os"
)

func TestDirectoryShouldReturnRecursiveDirs(t *testing.T) {
	entryDir, _ := ioutil.TempDir("", "")
	nestedDir, _ := ioutil.TempDir(entryDir, "")
	nestedDir2, _ := ioutil.TempDir(nestedDir, "")
	ioutil.TempFile(entryDir, "")
	ignoreDir, _ := ioutil.TempDir(entryDir, "")

	ignores := make(config.IgnoreFlags)
	ignores[ignoreDir] = struct{}{}
	helper := NewDirectory(entryDir, ignores)
	dirs, _ := helper.GetDirs()

	for _, dir := range []string{entryDir, nestedDir, nestedDir2} {
		assert.Contains(t, dirs, dir)
	}
	assert.NotContains(t, dirs, ignoreDir)
	os.RemoveAll(entryDir)
}

func TestDirectoryShouldIgnoreWildcardExtension(t *testing.T) {
	ignores := make(config.IgnoreFlags)
	ignores["*.swp"] = struct{}{}
	helper := NewDirectory(".", ignores)
	if !helper.IsIgnoreFile("README.md.swp") {
		t.Error("should ignore README.md.swp")
	}
}

func TestDirectoryShouldIgnoreWildcardEndsWith(t *testing.T) {
	ignores := make(config.IgnoreFlags)
	ignores["*~"] = struct{}{}
	helper := NewDirectory(".", ignores)
	if !helper.IsIgnoreFile("README.md~") {
		t.Error("should ignore README.md~")
	}
}

func TestDirectoryShouldIgnoreNestedWildcardExtension(t *testing.T) {
	ignores := make(config.IgnoreFlags)
	ignores["*.swp"] = struct{}{}
	helper := NewDirectory(".", ignores)
	if !helper.IsIgnoreFile("doc/README.md.swp") {
		t.Error("should ignore README.md.swp")
	}
}
