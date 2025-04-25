package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// Function that generates the layout as a string instead of printing it
func DrawLayout() string {
	width, _ := getTerminalSize()
	dir           := getCurrentDirectory()

	dirPrefix  := "Directory: "
	titleText  := "Cat Explorer"
	minSpacing := 2

	available := width - len(dirPrefix) - len(titleText) - minSpacing
	narrow    := len(dir) > available

	// 1ª línea: prefijo + espacios + título
	spaces := width - len(dirPrefix) - len(titleText)
	if spaces < minSpacing {
			spaces = minSpacing
	}
	header := fmt.Sprintf(
			"%s%s%s",
			DirectoryText.Render(dirPrefix),
			strings.Repeat(" ", spaces),
			HeaderTitle.Render(titleText),
	)

	// 2ª línea  en pantalla estrecha
	if narrow {
			// 2ª línea: directorio en sí
			header += "\n" + DirectoryDir.Render(dir)
	} else {
			// todo en una línea si cabe
			inLineSpaces := width - len(dirPrefix) - len(dir) - len(titleText)
			header = fmt.Sprintf(
					"%s%s%s%s",
					DirectoryText.Render(dirPrefix),
					DirectoryDir.Render(dir),
					strings.Repeat(" ", inLineSpaces),
					HeaderTitle.Render(titleText),
			)
	}

	header += "\n" // deja una línea en blanco antes de los paneles

	// Panel layout
	panelWidth := width / 3

	renderLeft := func(text string) string {
		padding := panelWidth - lipgloss.Width(text)
		if padding < 0 {
			padding = 0
		}
		return Cyan.Render(text) + strings.Repeat(" ", padding)
	}

	left := renderLeft("Directories")
	middle := renderLeft("Files")
	right := renderLeft("Preview Subdirectories")

	header += left + middle + right + "\n"

		// Datos para el panel izquierdo
		items := prepareDirItems(dir)
		selected := map[string]bool{}
		position := 0
		start := 0
		_, height := getTerminalSize()
		panelHeight := height - 5
		active := true
		includeSubdirs := false

		leftPanel := renderLeftPanel(items, selected, dir, position, start, panelHeight, panelWidth, active, includeSubdirs)

		return header + leftPanel
}

// Helper function to split directory path into multiple lines
func splitDirectory(dir string, maxWidth int) []string {
	var parts []string
	current := dir

	for len(current) > maxWidth {
		// Find the last separator before maxWidth
		splitIndex := strings.LastIndex(current[:maxWidth], "/")
		if splitIndex == -1 {
			splitIndex = maxWidth - 1
		}
		parts = append(parts, current[:splitIndex+1])
		current = current[splitIndex+1:]
	}
	if len(current) > 0 {
		parts = append(parts, current)
	}

	return parts
}

// Get the terminal size
func getTerminalSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width, height = 80, 24 // Default terminal size if error occurs
	}
	return width, height
}

// Get the current directory
func getCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		return "/" // Return root if there's an error getting the current directory
	}
	return dir
}

func renderLeftPanel(items []string, selected map[string]bool, directory string, position, start, height, width int, active bool, includeSubdirs bool) string {
	var b strings.Builder

	end := start + height
	if end > len(items) {
		end = len(items)
	}

	for i := start; i < end; i++ {
		item := items[i]
		fullPath := filepath.Join(directory, item)
		if item == ".." {
			fullPath = filepath.Dir(directory)
		}
		absPath, _ := filepath.Abs(fullPath)
		isSelected := selected[absPath]
		hasFocus := active && i == position

		marker := "  "
		if isSelected {
			if includeSubdirs {
				marker = "* "
			} else {
				marker = "• "
			}
		}
		content := marker + item
		maxWidth := width - 3
		if lipgloss.Width(content) > maxWidth {
			content = content[:maxWidth-3] + "..."
		}

		icon := "📁"
		line := icon + content

		// Rellenar hasta el ancho del panel
		padding := width - lipgloss.Width(line)
		if padding > 0 {
			line += strings.Repeat(" ", padding)
		}

		// Estilos
		if hasFocus {
			b.WriteString(Focus.Render(line) + "\n")
		} else if isSelected {
			b.WriteString(Selected.Render(line) + "\n")
		} else {
			b.WriteString(Green.Render(line) + "\n")
		}
	}

	// Scrollbar
	total := len(items)
	if total > height {
		barX := width - 1
		ratio := float64(start) / float64(total-height)
		thumb := int(ratio * float64(height-1))
		for y := 0; y < height; y++ {
			ch := "│"
			if y == thumb {
				ch = "█"
			}
			b.WriteString(lipgloss.PlaceHorizontal(width, lipgloss.Left, strings.Repeat(" ", barX)+ch) + "\n")
		}
	}

	return b.String()
}
