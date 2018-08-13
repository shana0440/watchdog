package dog

import (
	"github.com/shana0440/watchdog/config"

	"fmt"
	"io/ioutil"
	"path/filepath"
)

type DirHelper struct {
	ignores config.IgnoreFlags
}

// GetDirs will return recursive dirs under dir, excluded ignore dir
func (helper *DirHelper) GetDirs(dir string) ([]string, error) {
	dirs := []string{dir}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		newDir := fmt.Sprintf("%s/%s", dir, f.Name())
		if _, ok := helper.ignores[newDir]; !ok {
			dir, err := helper.GetDirs(newDir)
			if err != nil {
				return nil, err
			}
			dirs = append(dirs, dir...)
		}
	}
	return dirs, nil
}

// IsIgnoreFile will return file should be ignore or not
func (helper *DirHelper) IsIgnoreFile(file string) bool {
	for pattern, _ := range helper.ignores {
		match, err := filepath.Match(pattern, file)
		if err == nil && match {
			return true
		}
	}
	return false
}
