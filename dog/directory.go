package dog

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type Directoryer interface {
	GetDirs(...string) ([]string, error)
	IsIgnoreFile(string) bool
}

type Directory struct {
	entry   string
	ignores map[string]struct{}
}

func NewDirectory(entryDir string, ignores map[string]struct{}) *Directory {
	return &Directory{
		entry:   entryDir,
		ignores: ignores,
	}
}

// GetDirs will return recursive dirs under entry, excluded ignore dir
func (helper *Directory) GetDirs(dirs ...string) ([]string, error) {
	if len(dirs) == 0 {
		dirs = append(dirs, helper.entry)
	}
	currentDir := dirs[0]
	files, err := ioutil.ReadDir(currentDir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		newDir := fmt.Sprintf("%s/%s", currentDir, f.Name())
		if newDir[:2] == "./" {
			newDir = newDir[2:]
		}
		if _, ok := helper.ignores[newDir]; !ok {
			underDirs, err := helper.GetDirs(newDir)
			if err != nil {
				return nil, err
			}
			dirs = append(dirs, underDirs...)
		}
	}
	return dirs, nil
}

// IsIgnoreFile will return file should be ignore or not
func (helper *Directory) IsIgnoreFile(file string) bool {
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
