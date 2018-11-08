package dog

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Directoryer interface {
	GetDirs() []string
	ShouldIgnore(string) bool
}

type Directory struct {
	entry   string
	ignores map[string]struct{}
	matchs  []string
}

func NewDirectory(entryDir string, ignores []string, matchs []string) *Directory {
	ignoresMap := make(map[string]struct{})
	for _, path := range ignores {
		path = getRelPath(entryDir, path)
		ignoresMap[path] = struct{}{}
	}
	return &Directory{
		entry:   entryDir,
		ignores: ignoresMap,
		matchs:  matchs,
	}
}

// GetDirs will return recursive dirs under entry, excluded ignore dir
func (helper *Directory) GetDirs() []string {
	dirs := helper.getRecursiveDirs(helper.entry)
	return dirs
}

func (helper *Directory) getRecursiveDirs(dir string) []string {
	dirs := []string{dir}
	// dir must start with entry, and getRecursiveDirs only receive dirs under entry
	// so won't happend dir not exists error
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		child := fmt.Sprintf("%s/%s", dir, f.Name())
		child = getRelPath(helper.entry, child)
		if _, ok := helper.ignores[child]; !ok {
			childDirs := helper.getRecursiveDirs(child)
			dirs = append(dirs, childDirs...)
		}
	}
	return dirs
}

// ShouldIgnore will return file should be ignore or not
func (helper *Directory) ShouldIgnore(file string) bool {
	filename := filepath.Base(file)
	stat, err := os.Stat(filename)
	if err == nil {
		if stat.IsDir() {
			return false
		}
	}
	for pattern := range helper.ignores {
		matchPath, err := filepath.Match(pattern, file)
		matchFileName, err := filepath.Match(pattern, filename)
		if err == nil && (matchPath || matchFileName) {
			return true
		}
	}
	for _, pattern := range helper.matchs {
		matchPath, err := filepath.Match(pattern, file)
		matchFileName, err := filepath.Match(pattern, filename)
		if err == nil && (matchPath || matchFileName) {
			return false
		}
	}
	// ignore not match pattern if user specify match pattern
	// if user not specify, then allow not match pattern
	return len(helper.matchs) > 0
}

// getRelPath will make all path start with base and has consistency
func getRelPath(base, target string) string {
	path, err := filepath.Rel(base, target)
	if err != nil {
		panic(err)
	}
	return path
}
