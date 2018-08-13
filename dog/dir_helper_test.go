package dog

import (
	"github.com/shana0440/watchdog/config"
	"github.com/stretchr/testify/assert"

	"testing"

	"io/ioutil"
	"os"
)

func TestDirHelperShouldReturnRecursiveDirs(t *testing.T) {
	entryDir, _ := ioutil.TempDir("", "")
	nestedDir, _ := ioutil.TempDir(entryDir, "")
	nestedDir2, _ := ioutil.TempDir(nestedDir, "")
	ioutil.TempFile(entryDir, "")
	ignoreDir, _ := ioutil.TempDir(entryDir, "")

	ignores := make(config.IgnoreFlags)
	ignores[ignoreDir] = struct{}{}
	helper := DirHelper{ignores}
	dirs, _ := helper.GetDirs(entryDir)

	for _, dir := range []string{entryDir, nestedDir, nestedDir2} {
		assert.Contains(t, dirs, dir)
	}
	assert.NotContains(t, dirs, ignoreDir)
	os.RemoveAll(entryDir)
}

func TestDirHelperShouldIgnoreWildcardExtension(t *testing.T) {
	ignores := make(config.IgnoreFlags)
	ignores["*.swp"] = struct{}{}
	helper := DirHelper{ignores}
	if !helper.IsIgnoreFile("README.md.swp") {
		t.Error("should ignore README.md.swp")
	}
}

func TestDirHelperShouldIgnoreWildcardEndsWith(t *testing.T) {
	ignores := make(config.IgnoreFlags)
	ignores["*~"] = struct{}{}
	helper := DirHelper{ignores}
	if !helper.IsIgnoreFile("README.md~") {
		t.Error("should ignore README.md~")
	}
}

func TestDirHelperShouldIgnoreNestedWildcardExtension(t *testing.T) {
	ignores := make(config.IgnoreFlags)
	ignores["*.swp"] = struct{}{}
	helper := DirHelper{ignores}
	if !helper.IsIgnoreFile("doc/README.md.swp") {
		t.Error("should ignore README.md.swp")
	}
}
