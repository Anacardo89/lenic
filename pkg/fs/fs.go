package fs

import (
	"os"
	"path/filepath"
)

// gets the project root path
func ProjectRoot() (string, error) {
	// working dir
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	// checks for go.mod and goes up a dir ifdoesn't find it
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", os.ErrNotExist
}

// goes to the root of the project and makes a path to a file
func MakeFilePath(location, filename string) (string, error) {
	root, err := ProjectRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, location, filename), nil
}
