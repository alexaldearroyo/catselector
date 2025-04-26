package core

import (
	"os"
	"os/exec"
	"runtime"
	"sort"
)

var rootDirectory string

func OpenTextFile(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", path)
	case "windows": // Windows
		cmd = exec.Command("cmd", "/c", "start", "", path)
	default: // Linux y otros
		// Intentar xdg-open (estándar para la mayoría de distribuciones Linux)
		if _, err := exec.LookPath("xdg-open"); err == nil {
			cmd = exec.Command("xdg-open", path)
		} else if _, err := exec.LookPath("gnome-open"); err == nil {
			cmd = exec.Command("gnome-open", path)
		} else if _, err := exec.LookPath("kde-open"); err == nil {
			cmd = exec.Command("kde-open", path)
		} else {
			return os.ErrNotExist
		}
	}

	return cmd.Start()
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

	// Añadir ".." como primer item si no estamos en el directorio root
	rootDir := GetRootDirectory()
	if pwd != rootDir {
		return append([]string{"..", "."}, dirs...)
	}
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
