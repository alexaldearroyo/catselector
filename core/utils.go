package core

import (
	"os"
	"os/exec"
	"sort"
)

func OpenTextFile(path string) {
	cmd := exec.Command("xdg-open", path)
	cmd.Start()
}


func prepareDirItems(pwd string) []string {
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
