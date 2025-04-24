package core

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenTextFile abre un archivo de texto con la aplicación predeterminada del sistema
func OpenTextFile(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", filePath)
	case "windows": // Windows
		cmd = exec.Command("cmd", "/c", "start", filePath)
	default: // Linux y otros
		cmd = exec.Command("xdg-open", filePath)
	}

	return cmd.Run()
}

// SetTerminalTitle establece el título de la ventana de la terminal
func SetTerminalTitle(title string) {
	switch runtime.GOOS {
	case "darwin", "linux": // macOS y Linux
		fmt.Printf("\033]0;%s\007", title)
	case "windows": // Windows
		// No hay una solución universal para Windows,
		// ya que depende del terminal específico utilizado
		cmd := exec.Command("title", title)
		_ = cmd.Run()
	}
}
