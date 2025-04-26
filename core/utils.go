package core

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
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

// RenderLeft renderiza el texto del encabezado con el estilo apropiado
func RenderLeft(text string, isActive bool, isCounter bool, panelWidth int) string {
	padding := panelWidth - lipgloss.Width(text)
	if padding < 0 {
		padding = 0
	}
	if isActive {
		return ActiveHeader.Render(text + strings.Repeat(" ", padding))
	}
	if isCounter {
		return Blue.Render(text) + strings.Repeat(" ", padding)
	}
	return Cyan.Render(text) + strings.Repeat(" ", padding)
}

// IsBinaryFile verifica si un archivo es binario
func IsBinaryFile(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 1024)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return false
	}

	nullCount := 0
	nonPrintableCount := 0
	for i := 0; i < n; i++ {
		if buf[i] == 0 {
			nullCount++
		}
		if buf[i] < 32 && buf[i] != '\t' && buf[i] != '\n' && buf[i] != '\r' {
			nonPrintableCount++
		}
	}

	return nullCount > 10 || float64(nonPrintableCount)/float64(n) > 0.3
}

// CopyToClipboard copia texto al portapapeles del sistema
func CopyToClipboard(text string) bool {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("pbcopy")
	case "windows":
		// En Windows, necesitarías usar otro método, como el paquete clipboard
		return false
	default: // Linux y otros
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if _, err := exec.LookPath("xsel"); err == nil {
			cmd = exec.Command("xsel", "--clipboard", "--input")
		} else if _, err := exec.LookPath("wl-copy"); err == nil {
			cmd = exec.Command("wl-copy")
		} else {
			return false
		}
	}

	if cmd == nil {
		return false
	}

	pipe, err := cmd.StdinPipe()
	if err != nil {
		return false
	}

	if err := cmd.Start(); err != nil {
		return false
	}

	_, err = pipe.Write([]byte(text))
	if err != nil {
		return false
	}

	pipe.Close()
	return cmd.Wait() == nil
}

// ShowErrorMessage muestra un mensaje de error formateado
func ShowErrorMessage(b *strings.Builder, prefix, filePath string, width, height int) {
	errorMsg := prefix
	if len(filePath) > width-20 {
		errorMsg += ": " + filePath[:width-25] + "..."
	} else {
		errorMsg += ": " + filePath
	}

	if lipgloss.Width(errorMsg) > width {
		errorMsg = errorMsg[:width-3] + "..."
	}

	padding := width - lipgloss.Width(errorMsg)
	if padding > 0 {
		errorMsg += strings.Repeat(" ", padding)
	}

	b.WriteString(White.Render(errorMsg) + "\n")

	for i := 1; i < height; i++ {
		b.WriteString(strings.Repeat(" ", width) + "\n")
	}
}

// ShowBinaryFileMessage muestra un mensaje para archivos binarios
func ShowBinaryFileMessage(b *strings.Builder, filePath string, size int64, width, height int) {
	sizeStr := FormatFileSize(size)

	line1 := "Binary file"
	if lipgloss.Width(line1) > width {
		line1 = line1[:width-3] + "..."
	}
	padding1 := width - lipgloss.Width(line1)
	if padding1 > 0 {
		line1 += strings.Repeat(" ", padding1)
	}
	b.WriteString(White.Render(line1) + "\n")

	line2 := "Size: " + sizeStr
	padding2 := width - lipgloss.Width(line2)
	if padding2 > 0 {
		line2 += strings.Repeat(" ", padding2)
	}
	b.WriteString(White.Render(line2) + "\n")

	line3 := "Preview not available"
	padding3 := width - lipgloss.Width(line3)
	if padding3 > 0 {
		line3 += strings.Repeat(" ", padding3)
	}
	b.WriteString(White.Render(line3) + "\n")

	for i := 3; i < height; i++ {
		b.WriteString(strings.Repeat(" ", width) + "\n")
	}
}

// FormatFileSize formatea el tamaño de un archivo en una cadena legible
func FormatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// SanitizeLine elimina caracteres problemáticos de una línea
func SanitizeLine(line string) string {
	var result strings.Builder
	for _, r := range line {
		if r < 32 && r != '\t' && r != '\n' && r != '\r' {
			result.WriteRune(' ')
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
