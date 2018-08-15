package dog

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type Directoryer interface {
	GetDirs() []string
	ShouldIgnore(string) bool
}

type Directory struct {
	entry   string
	ignores map[string]struct{}
}

func NewDirectory(entryDir string, ignores []string) *Directory {
	ignoresMap := make(map[string]struct{})
	for _, path := range ignores {
		path = getRelPath(entryDir, path)
		ignoresMap[path] = struct{}{}
	}
	return &Directory{
		entry:   entryDir,
		ignores: ignoresMap,
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
	for pattern := range helper.ignores {
		matchPath, err := filepath.Match(pattern, file)
		matchFileName, err := filepath.Match(pattern, filename)
		if err == nil && (matchPath || matchFileName) {
			return true
		}
	}
	return false
}

// getRelPath will make all path start with base and has consistency
func getRelPath(base, target string) string {
	path, err := filepath.Rel(base, target)
	if err != nil {
		panic(err)
	}
	return path
}
