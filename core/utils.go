package core

import (
	"os"
	"os/exec"
	"sort"
)

var rootDirectory string

func OpenTextFile(path string) {
	cmd := exec.Command("xdg-open", path)
	cmd.Start()
}

// GetRootDirectory devuelve el directorio desde donde se ejecuta la aplicación
func GetRootDirectory() string {
	if rootDirectory == "" {
		dir, err := os.Getwd()
		if err != nil {
			rootDirectory = "/"
		} else {
			rootDirectory = dir
		}
	}
	return rootDirectory
}

func PrepareDirItems(pwd string) []string {
	files, _ := os.ReadDir(pwd)
	var dirs []string
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, f.Name())
		}
	}
	sort.Strings(dirs)
	return append([]string{"."}, dirs...)
}

// Get the current directory
func GetCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		return "/" // Return root if there's an error getting the current directory
	}
	return dir
}
