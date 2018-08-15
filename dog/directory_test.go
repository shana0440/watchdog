package dog

import (
	"testing"

	"os"

	"github.com/stretchr/testify/assert"
)

func TestGetDirsShould(t *testing.T) {
	t.Run("ReturnRecursiveDirs", func(t *testing.T) {
		dir := NewDirectory(".", []string{})
		os.MkdirAll("./should/return/recursive/dir", 0777)
		os.MkdirAll("./should/return/dir", 0777)

		actual := dir.GetDirs()

		assert.ElementsMatch(
			t,
			[]string{
				".",
				"should",
				"should/return",
				"should/return/dir",
				"should/return/recursive",
				"should/return/recursive/dir",
			},
			actual,
		)
	})

	t.Run("IgnoreExcludeDirs", func(t *testing.T) {
		dir := NewDirectory(".", []string{"should/return/dir", "./should/return/recursive"})
		os.MkdirAll("./should/return/recursive/dir", 0777)
		os.MkdirAll("./should/return/dir", 0777)

		actual := dir.GetDirs()

		assert.ElementsMatch(
			t,
			[]string{
				".",
				"should",
				"should/return",
			},
			actual,
		)
	})
}

func TestShouldIgnoreShould(t *testing.T) {
	t.Run("ReturnTrueWhenFileMatchIgnorePattern", func(t *testing.T) {
		dir := NewDirectory(".", []string{"*.swp", "*~", "dir/*.db"})
		assert.True(t, dir.ShouldIgnore("hello.swp"))
		assert.True(t, dir.ShouldIgnore("hello~"))
		assert.True(t, dir.ShouldIgnore("tmp/hello.swp"))
		assert.True(t, dir.ShouldIgnore("tmp/hello~"))
		assert.True(t, dir.ShouldIgnore("dir/tmp.db"))
	})

	t.Run("ReturnFalseWhenFileNotMatchIgnorePattern", func(t *testing.T) {
		dir := NewDirectory(".", []string{"*.swp", "*~", "dir/*.db"})
		assert.False(t, dir.ShouldIgnore("hello.md"))
		assert.False(t, dir.ShouldIgnore("hello"))
		assert.False(t, dir.ShouldIgnore("tmp.db"))
	})
}

func TestGetRelPathShould(t *testing.T) {
	t.Run("ReturnRelPathWithoutBasePath", func(t *testing.T) {
		actual := getRelPath(".", "./tmp")
		assert.Equal(t, "tmp", actual)
	})
}
