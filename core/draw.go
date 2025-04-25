package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
		// Dividir el directorio en partes
		parts := strings.Split(currentDir, "/")
		lastPart := parts[len(parts)-1]
		parentDir := strings.Join(parts[:len(parts)-1], "/")
		if parentDir != "" {
			parentDir += "/"
		}
		header += "\n" + DirectoryText.Render(parentDir) + DirectoryDir.Render(lastPart)
	} else {
		// Dividir el directorio en partes
		parts := strings.Split(currentDir, "/")
		lastPart := parts[len(parts)-1]
		parentDir := strings.Join(parts[:len(parts)-1], "/")
		if parentDir != "" {
			parentDir += "/"
		}
		inLineSpaces := width - len(dirPrefix) - len(currentDir) - len(titleText)
		header = fmt.Sprintf(
			"%s%s%s%s%s",
			DirectoryText.Render(dirPrefix),
			DirectoryText.Render(parentDir),
			DirectoryDir.Render(lastPart),
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

	header += left + White.Render("│") + middle + White.Render("│") + right + "\n"

	// Panel izquierdo (Directories)
	selected := map[string]bool{}
	start := 0
	panelHeight := height - 5
	active := true
	includeSubdirs := false

	leftPanel := renderLeftPanel(items, selected, currentDir, position, start, panelHeight, panelWidth, active, includeSubdirs)

	// Panel de archivos (Files)
	filePanel := renderFilePanel(files, position, panelWidth, height, panelHeight)

	// Panel derecho (Preview Subdirectories)
	var selectedDir string
	if position < len(items) {
		item := items[position]
		if item == ".." {
			selectedDir = filepath.Dir(currentDir)
		} else if item == "." {
			selectedDir = currentDir
		} else {
			selectedDir = filepath.Join(currentDir, item)
		}
	}
	rightPanel := renderPreviewPanel(selectedDir, panelWidth, panelHeight)

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

	// Asegurarnos de que tengamos suficientes líneas para llenar el buffer
	for maxLines < panelHeight {
		leftLines = append(leftLines, "")
		fileLines = append(fileLines, "")
		rightLines = append(rightLines, "")
		maxLines++
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

		// Añadir líneas verticales blancas entre los paneles
		result.WriteString(leftLine + White.Render("│") + fileLine + White.Render("│") + rightLine + "\n")
	}

	return result.String()
}

// Esta función debería manejar el renderizado de los archivos
func renderFilePanel(files []string, position, panelWidth, height, panelHeight int) string {
	var b strings.Builder

	for i := 0; i < panelHeight && i < len(files); i++ {
		file := files[i]
		icon := GetFileIcon(file)
		line := icon + "  " + file

		// Truncar el nombre del archivo si es demasiado largo
		maxWidth := panelWidth - 2 // Dejamos espacio para el scrollbar
		if lipgloss.Width(line) > maxWidth {
			// Calcular cuánto espacio tenemos para el nombre del archivo
			iconWidth := lipgloss.Width(icon + "  ")
			availableWidth := maxWidth - iconWidth - 3 // 3 para "..."

			// Truncar el nombre del archivo
			if availableWidth > 0 {
				truncatedName := file
				if len(truncatedName) > availableWidth {
					truncatedName = truncatedName[:availableWidth] + "..."
				}
				line = icon + "  " + truncatedName
			}
		}

		// Solo usar el estilo White para los archivos
		b.WriteString(White.Render(line) + "\n")
	}

	return b.String()
}

// renderPreviewPanel muestra los subdirectorios del directorio seleccionado
func renderPreviewPanel(dir string, width, height int) string {
	var b strings.Builder

	// Verificar si el directorio existe y es accesible
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return strings.Repeat(" ", width) + "\n"
	}

	// Leer los subdirectorios
	entries, err := os.ReadDir(dir)
	if err != nil {
		return strings.Repeat(" ", width) + "\n"
	}

	// Filtrar solo directorios y ordenarlos
	var subdirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			subdirs = append(subdirs, entry.Name())
		}
	}
	sort.Strings(subdirs)

	// Mostrar los subdirectorios
	for i := 0; i < height && i < len(subdirs); i++ {
		subdir := subdirs[i]
		icon := GetFileIcon(filepath.Join(dir, subdir))
		line := icon + "  " + subdir

		// Rellenar hasta el ancho del panel
		padding := width - lipgloss.Width(line)
		if padding > 0 {
			line += strings.Repeat(" ", padding)
		}

		b.WriteString(Green.Render(line) + "\n")
	}

	return b.String()
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

		icon := GetFileIcon(fullPath)
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


	return b.String()
}
