package core

import (
	"os/exec"
)

func OpenTextFile(path string) {
	cmd := exec.Command("xdg-open", path)
	cmd.Start()
}
