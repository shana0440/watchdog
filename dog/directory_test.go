package dog

import (
	"testing"

	"io/ioutil"
	"os"

	"github.com/stretchr/testify/assert"
)

func TestDirectoryShould(t *testing.T) {
	t.Run("ReturnRecursiveDirs", func(t *testing.T) {
		entryDir, _ := ioutil.TempDir("", "")
		nestedDir, _ := ioutil.TempDir(entryDir, "")
		nestedDir2, _ := ioutil.TempDir(nestedDir, "")
		ioutil.TempFile(entryDir, "")
		ignoreDir, _ := ioutil.TempDir(entryDir, "")

		ignores := make(map[string]struct{})
		ignores[ignoreDir] = struct{}{}
		helper := NewDirectory(entryDir, ignores)
		dirs, _ := helper.GetDirs()

		for _, dir := range []string{entryDir, nestedDir, nestedDir2} {
			assert.Contains(t, dirs, dir)
		}
		assert.NotContains(t, dirs, ignoreDir)
		os.RemoveAll(entryDir)
	})
	t.Run("IgnoreWildcardExtension", func(t *testing.T) {
		ignores := make(map[string]struct{})
		ignores["*.swp"] = struct{}{}
		helper := NewDirectory(".", ignores)
		if !helper.IsIgnoreFile("README.md.swp") {
			t.Error("should ignore README.md.swp")
		}
	})
	t.Run("IgnoreWildcardEndsWith", func(t *testing.T) {
		ignores := make(map[string]struct{})
		ignores["*~"] = struct{}{}
		helper := NewDirectory(".", ignores)
		if !helper.IsIgnoreFile("README.md~") {
			t.Error("should ignore README.md~")
		}
	})
	t.Run("IgnoreNestedWildcardExtension", func(t *testing.T) {
		ignores := make(map[string]struct{})
		ignores["*.swp"] = struct{}{}
		helper := NewDirectory(".", ignores)
		if !helper.IsIgnoreFile("doc/README.md.swp") {
			t.Error("should ignore README.md.swp")
		}
	})
}
