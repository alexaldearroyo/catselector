package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func DrawLayout(position int, items []string, currentDir string, files []string) string {
	width, height := getTerminalSize()
	dirPrefix := "Directory: "
	titleText := "Cat Explorer"
	minSpacing := 2

	available := width - len(dirPrefix) - len(titleText) - minSpacing
	narrow := len(currentDir) > available

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

	// 2ª línea en pantalla estrecha
	if narrow {
		header += "\n" + DirectoryDir.Render(currentDir)
	} else {
		inLineSpaces := width - len(dirPrefix) - len(currentDir) - len(titleText)
		header = fmt.Sprintf(
			"%s%s%s%s",
			DirectoryText.Render(dirPrefix),
			DirectoryDir.Render(currentDir),
			strings.Repeat(" ", inLineSpaces),
			HeaderTitle.Render(titleText),
		)
	}

	header += "\n" // dejar una línea en blanco antes de los paneles

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

	// Panel izquierdo (Directories)
	selected := map[string]bool{}
	start := 0
	panelHeight := height - 5
	active := true
	includeSubdirs := false

	leftPanel := renderLeftPanel(items, selected, currentDir, position, start, panelHeight, panelWidth, active, includeSubdirs)

	// Panel de archivos (Files)
	filePanel := renderFilePanel(files, position, panelWidth, height, panelHeight)

	// Panel derecho (Preview)
	rightPanel := strings.Repeat(" ", panelWidth) + "\n"

	// Combinar los paneles horizontalmente
	var result strings.Builder
	result.WriteString(header)

	// Dividir los paneles en líneas
	leftLines := strings.Split(leftPanel, "\n")
	fileLines := strings.Split(filePanel, "\n")
	rightLines := strings.Split(rightPanel, "\n")

	// Encontrar el máximo número de líneas
	maxLines := len(leftLines)
	if len(fileLines) > maxLines {
		maxLines = len(fileLines)
	}
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}

	// Combinar las líneas horizontalmente
	for i := 0; i < maxLines; i++ {
		leftLine := ""
		if i < len(leftLines) {
			leftLine = leftLines[i]
		}
		fileLine := ""
		if i < len(fileLines) {
			fileLine = fileLines[i]
		}
		rightLine := ""
		if i < len(rightLines) {
			rightLine = rightLines[i]
		}

		// Asegurar que cada línea tenga el ancho correcto
		leftLine = lipgloss.PlaceHorizontal(panelWidth, lipgloss.Left, leftLine)
		fileLine = lipgloss.PlaceHorizontal(panelWidth, lipgloss.Left, fileLine)
		rightLine = lipgloss.PlaceHorizontal(panelWidth, lipgloss.Left, rightLine)

		result.WriteString(leftLine + fileLine + rightLine + "\n")
	}

	return result.String()
}

// Esta función debería manejar el renderizado de los archivos
func renderFilePanel(files []string, position, panelWidth, height, panelHeight int) string {
	var b strings.Builder

	for i := 0; i < panelHeight && i < len(files); i++ {
		file := files[i]
		hasFocus := i == position
		icon := "📄"
		line := icon + " " + file

		// Estilos
		if hasFocus {
			b.WriteString(Focus.Render(line) + "\n")
		} else {
			b.WriteString(White.Render(line) + "\n")
		}
	}

	return b.String()
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
