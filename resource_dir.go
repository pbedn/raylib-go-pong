package main

import (
	"os"
	"path/filepath"
)

// SearchAndSetResourceDir looks for folderName in the current working
// directory, the executable's directory, and up to three parent directories
// of the executable's directory. When found, it changes the working directory
// to that resource directory.
func SearchAndSetResourceDir(folderName string) bool {
	if folderName == "" {
		return false
	}

	workingDir, err := os.Getwd()
	if err == nil && changeToResourceDir(filepath.Join(workingDir, folderName)) {
		return true
	}

	executable, err := os.Executable()
	if err != nil {
		return false
	}

	appDir := filepath.Dir(executable)
	for level := 0; level <= 3; level++ {
		if changeToResourceDir(filepath.Join(appDir, folderName)) {
			return true
		}
		appDir = filepath.Dir(appDir)
	}

	return false
}

func changeToResourceDir(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return false
	}

	return os.Chdir(dir) == nil
}
