package dog

import (
	"testing"

	"os"

	"github.com/stretchr/testify/assert"
)

func TestGetDirsShould(t *testing.T) {
	t.Run("ReturnRecursiveDirs", func(t *testing.T) {
		dir := NewDirectory(".", []string{}, []string{})
		os.MkdirAll("./should/return/recursive/dir", 0777)
		os.MkdirAll("./should/return/dir", 0777)
		defer os.RemoveAll("./should")

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
		dir := NewDirectory(".", []string{"should/return/dir", "./should/return/recursive"}, []string{})
		os.MkdirAll("./should/return/recursive/dir", 0777)
		os.MkdirAll("./should/return/dir", 0777)
		defer os.RemoveAll("./should")

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
		dir := NewDirectory(".", []string{"*.swp", "*~", "dir/*.db"}, []string{})
		assert.True(t, dir.ShouldIgnore("hello.swp"))
		assert.True(t, dir.ShouldIgnore("hello~"))
		assert.True(t, dir.ShouldIgnore("tmp/hello.swp"))
		assert.True(t, dir.ShouldIgnore("tmp/hello~"))
		assert.True(t, dir.ShouldIgnore("dir/tmp.db"))
	})

	t.Run("ReturnFalseWhenFileNotMatchIgnorePattern", func(t *testing.T) {
		dir := NewDirectory(".", []string{"*.swp", "*~", "dir/*.db"}, []string{})
		assert.False(t, dir.ShouldIgnore("hello.md"))
		assert.False(t, dir.ShouldIgnore("hello"))
		assert.False(t, dir.ShouldIgnore("tmp.db"))
	})

	t.Run("ReturnFalseWhenFileMatchMatchPattern", func(t *testing.T) {
		dir := NewDirectory(".", []string{}, []string{"*.go", "dir/*.db"})
		assert.False(t, dir.ShouldIgnore("hello.go"))
		assert.False(t, dir.ShouldIgnore("tmp/hello.go"))
		assert.False(t, dir.ShouldIgnore("dir/tmp.db"))
	})

	t.Run("ReturnTrueWhenFileNotMatchMatchPattern", func(t *testing.T) {
		dir := NewDirectory(".", []string{}, []string{"*.go", "dir/*.db"})
		assert.True(t, dir.ShouldIgnore("hello.php"))
		assert.True(t, dir.ShouldIgnore("tmp.db"))
	})

	t.Run("ReturnTrueWhenFileMatchMatchPatternAndIgnorePattern", func(t *testing.T) {
		dir := NewDirectory(".", []string{"*.go"}, []string{"*.go"})
		assert.True(t, dir.ShouldIgnore("hello.go"))
	})
}

func TestGetRelPathShould(t *testing.T) {
	t.Run("ReturnRelPathWithoutBasePath", func(t *testing.T) {
		actual := getRelPath(".", "./tmp")
		assert.Equal(t, "tmp", actual)
	})
}
