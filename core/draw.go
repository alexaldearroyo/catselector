package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func DrawLayout(position int, items []string, currentDir string, files []string, activePanel int, filePosition int) string {
	width, height := getTerminalSize()
	dirPrefix := ""
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

		// Añadir new_text en una línea separada después del pwd
		// Reemplazar con el texto de subdirectorios y selección
		selector := GetCurrentSelector()
		var subdirText, selectedText string

		// Determinar el estado de los subdirectorios
		if selector.IncludeMode {
			subdirText = White.Render("Mode: ") + Magenta.Render("Include")
		} else {
			subdirText = White.Render("Mode: ") + Magenta.Render("Normal")
		}

		// Contar archivos y directorios seleccionados
		selectedFiles, selectedDirs := countSelected(selector)
		selectedText = White.Render("Selected: ") +
			Magenta.Render(fmt.Sprintf("%d", selectedFiles)) +
			White.Render(" Files") +
			White.Render(", ") +
			Magenta.Render(fmt.Sprintf("%d", selectedDirs)) +
			White.Render(" Directories")

		// Texto completo con la parte de Selected alineada a la derecha
		infoText := subdirText
		// Calcular espacios para alinear a la derecha
		spaces := width - lipgloss.Width(subdirText) - lipgloss.Width(selectedText)
		if spaces > 0 {
			infoText += strings.Repeat(" ", spaces)
		}
		infoText += selectedText
		header += "\n" + infoText
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

		// Añadir un salto de línea y luego el texto informativo
		// Reemplazar con el texto de subdirectorios y selección
		selector := GetCurrentSelector()
		var subdirText, selectedText string

		// Determinar el estado de los subdirectorios
		if selector.IncludeMode {
			subdirText = White.Render("Mode: ") + Magenta.Render("Include")
		} else {
			subdirText = White.Render("Mode: ") + White.Render("Normal")
		}

		// Contar archivos y directorios seleccionados
		selectedFiles, selectedDirs := countSelected(selector)
		selectedText = White.Render("Selected: ") +
			Magenta.Render(fmt.Sprintf("%d", selectedFiles)) +
			White.Render(" Files") +
			White.Render(", ") +
			Magenta.Render(fmt.Sprintf("%d", selectedDirs)) +
			White.Render(" Directories")

		// Texto completo con la parte de Selected alineada a la derecha
		infoText := subdirText
		// Calcular espacios para alinear a la derecha
		spaces := width - lipgloss.Width(subdirText) - lipgloss.Width(selectedText)
		if spaces > 0 {
			infoText += strings.Repeat(" ", spaces)
		}
		infoText += selectedText
		header += "\n" + infoText
	}

	// Añadir un último salto de línea antes de los paneles
	header += "\n"

	// Panel layout
	panelWidth := width / 3

	renderLeft := func(text string, isActive bool, isCounter bool) string {
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

	// Obtener el selector actual para verificar el modo include
	// selector := GetCurrentSelector()
	includeModeText := ""
	// if selector.IncludeMode {
	// 	includeModeText = " [Include Mode]"
	// }

	// Contar elementos para cada panel
	var totalItems, totalFiles, totalSubdirs int
	var err error

	// Si estamos en el panel de directorios
	if activePanel == 1 && position >= 0 && position < len(items) {
		item := items[position]
		var selectedDir string
		if item == ".." {
			selectedDir = filepath.Dir(currentDir)
		} else if item == "." {
			selectedDir = currentDir
		} else {
			selectedDir = filepath.Join(currentDir, item)
		}
		totalItems, err = countItems(selectedDir)
		if err == nil {
			totalFiles, _ = countFiles(selectedDir)
			totalSubdirs, _ = countSubdirs(selectedDir)
		}
	} else {
		// Si estamos en el panel de archivos o no hay directorio seleccionado
		totalItems, err = countItems(currentDir)
		if err == nil {
			totalFiles, _ = countFiles(currentDir)
			totalSubdirs, _ = countSubdirs(currentDir)
		}
	}

	// Añadir contadores a los encabezados
	left := renderLeft("Directories"+includeModeText, activePanel == 1, false)
	middle := renderLeft("Files", activePanel == 2, false)
	right := renderLeft("Preview", activePanel == 3, false)

	// Añadir los contadores en una línea nueva
	leftCounter := renderLeft(fmt.Sprintf("%d items", totalItems), false, true)
	middleCounter := renderLeft(fmt.Sprintf("%d files", totalFiles), false, true)
	rightCounter := renderLeft(fmt.Sprintf("%d subdirs", totalSubdirs), false, true)

	// Combinar las cabeceras y contadores
	header += left + White.Render("│") + middle + White.Render("│") + right + "\n"
	header += leftCounter + White.Render("│") + middleCounter + White.Render("│") + rightCounter + "\n"

	// Panel izquierdo (Directories)
	selected := map[string]bool{}
	start := 0
	panelHeight := height - 6  // Ajustado para considerar la línea adicional
	active := activePanel == 1
	includeSubdirs := false

	leftPanel := renderLeftPanel(items, selected, currentDir, position, start, panelHeight, panelWidth, active, includeSubdirs)

	// Panel de archivos (Files)
	filePanel := renderFilePanel(files, position, panelWidth, height, panelHeight, activePanel, filePosition)

	// Panel derecho (Preview)
	rightPanel := renderPreviewPanel(currentDir, panelWidth, panelHeight, files, filePosition, activePanel, items, position)

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

	// Asegurarnos de que tengamos suficientes líneas
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
func renderFilePanel(files []string, position, panelWidth, height, panelHeight int, activePanel int, filePosition int) string {
	var b strings.Builder

	// Obtener el selector actual
	selector := GetCurrentSelector()

	// Verificar si el directorio actual o el directorio seleccionado está seleccionado
	currentDirSelected := selector.IsSelected(".")
	selectedDirSelected := false
	if selector.ActivePanel == 1 && selector.Position < len(selector.Filtered) {
		selectedItem := selector.Filtered[selector.Position]
		if selectedItem != ".." {
			selectedDirSelected = selector.IsSelected(selectedItem)
		}
	}

	for i := 0; i < panelHeight && i < len(files); i++ {
		file := files[i]
		icon := GetFileIcon(file)

		// Verificar si el archivo está seleccionado o si algún directorio padre está seleccionado
		isSelected := selector.IsFileSelected(file) || currentDirSelected || selectedDirSelected

		// Añadir asterisco si está seleccionado
		marker := "  "
		if isSelected {
			marker = " *"
		}

		line := icon + marker + file

		// Truncar el nombre del archivo si es demasiado largo
		maxWidth := panelWidth - 2 // Dejamos espacio para el scrollbar
		if lipgloss.Width(line) > maxWidth {
			// Calcular cuánto espacio tenemos para el nombre del archivo
			iconWidth := lipgloss.Width(icon + marker)
			availableWidth := maxWidth - iconWidth - 3 // 3 para "..."

			// Truncar el nombre del archivo
			if availableWidth > 0 {
				truncatedName := file
				if len(truncatedName) > availableWidth {
					truncatedName = truncatedName[:availableWidth] + "..."
				}
				line = icon + marker + truncatedName
			}
		}

		// Aplicar el estilo Focus si el panel está activo y este es el archivo seleccionado
		if activePanel == 2 && i == filePosition {
			// Rellenar hasta el ancho del panel
			padding := panelWidth - lipgloss.Width(line)
			if padding > 0 {
				line += strings.Repeat(" ", padding)
			}
			b.WriteString(Focus.Render(line) + "\n")
		} else if isSelected {
			// Aplicar estilo amarillo para archivos seleccionados
			b.WriteString(Yellow.Render(line) + "\n")
		} else {
			b.WriteString(White.Render(line) + "\n")
		}
	}

	return b.String()
}

// renderPreviewPanel muestra el contenido del archivo seleccionado o los subdirectorios
func renderPreviewPanel(dir string, width, height int, files []string, filePosition int, activePanel int, items []string, position int) string {
	var b strings.Builder

	// Obtener el selector actual
	selector := GetCurrentSelector()

	// Si el panel activo es el de directorios, mostrar subdirectorios
	if activePanel == 1 {
		// Determinar el directorio seleccionado
		var selectedDir string
		if position >= 0 && position < len(items) {
			item := items[position]
			if item == ".." {
				selectedDir = filepath.Dir(dir)
			} else if item == "." {
				selectedDir = dir
			} else {
				selectedDir = filepath.Join(dir, item)
			}
		} else {
			selectedDir = dir
		}

		// Verificar si el directorio existe y es accesible
		info, err := os.Stat(selectedDir)
		if err != nil || !info.IsDir() {
			return strings.Repeat(" ", width) + "\n"
		}

		// Leer los subdirectorios
		entries, err := os.ReadDir(selectedDir)
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
			icon := GetFileIcon(filepath.Join(selectedDir, subdir))

			// Verificar si el subdirectorio está seleccionado
			isSelected := selector.Selection[filepath.Join(selectedDir, subdir)]

			// Añadir el marcador correspondiente
			marker := "  "
			if isSelected {
				marker = " •"
			}

			line := icon + marker + subdir

			// Rellenar hasta el ancho del panel
			padding := width - lipgloss.Width(line)
			if padding > 0 {
				line += strings.Repeat(" ", padding)
			}

			// Aplicar el estilo correspondiente
			if isSelected {
				b.WriteString(Yellow.Render(line) + "\n")
			} else {
				b.WriteString(Green.Render(line) + "\n")
			}
		}

		// Rellenar con líneas vacías si es necesario
		for i := len(subdirs); i < height; i++ {
			b.WriteString(strings.Repeat(" ", width) + "\n")
		}
	} else {
		// Si estamos en el panel de archivos, mostrar el contenido del archivo seleccionado
		if len(files) > 0 && filePosition >= 0 && filePosition < len(files) {
			// Determinar el directorio actual para los archivos
			var currentDir string
			if position >= 0 && position < len(items) {
				item := items[position]
				if item == ".." {
					currentDir = filepath.Dir(dir)
				} else if item == "." {
					currentDir = dir
				} else {
					currentDir = filepath.Join(dir, item)
				}
			} else {
				currentDir = dir
			}

			filePath := filepath.Join(currentDir, files[filePosition])

			// Verificar si el archivo es binario o demasiado grande
			info, err := os.Stat(filePath)
			if err != nil {
				// Error al acceder al archivo
				showErrorMessage(&b, "No se puede acceder al archivo", filePath, width, height)
				return b.String()
			}

			// Verificar si es un archivo binario o demasiado grande
			if isBinaryFile(filePath) || info.Size() > 1024*1024 { // Más de 1MB
				showBinaryFileMessage(&b, filePath, info.Size(), width, height)
				return b.String()
			}

			// Leer el contenido del archivo
			content, err := os.ReadFile(filePath)
			if err == nil {
				// Convertir el contenido a string y limitar a las primeras líneas que quepan
				lines := strings.Split(string(content), "\n")

				// Limitar el número de líneas para evitar desbordamientos
				maxLines := height
				if len(lines) > maxLines {
					lines = lines[:maxLines]
				}

				for i := 0; i < len(lines); i++ {
					line := lines[i]

					// Sanitizar la línea para evitar caracteres problemáticos
					line = sanitizeLine(line)

					// Truncar la línea si es demasiado larga
					if lipgloss.Width(line) > width {
						line = line[:width-3] + "..."
					}

					// Rellenar hasta el ancho del panel
					padding := width - lipgloss.Width(line)
					if padding > 0 {
						line += strings.Repeat(" ", padding)
					}

					b.WriteString(White.Render(line) + "\n")
				}

				// Rellenar con líneas vacías si es necesario
				for i := len(lines); i < height; i++ {
					b.WriteString(strings.Repeat(" ", width) + "\n")
				}
			} else {
				// Error al leer el archivo
				showErrorMessage(&b, "No se puede leer el archivo", filePath, width, height)
			}
		} else {
			// Mensaje cuando no hay archivo seleccionado
			msg := "No hay archivo seleccionado"
			padding := width - lipgloss.Width(msg)
			if padding > 0 {
				msg += strings.Repeat(" ", padding)
			}
			b.WriteString(White.Render(msg) + "\n")

			// Rellenar con líneas vacías
			for i := 1; i < height; i++ {
				b.WriteString(strings.Repeat(" ", width) + "\n")
			}
		}
	}

	return b.String()
}

// isBinaryFile verifica si un archivo es binario
func isBinaryFile(filePath string) bool {
	// Abrir el archivo
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	// Leer los primeros 1024 bytes
	buf := make([]byte, 1024)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return false
	}

	// Verificar si contiene caracteres nulos o demasiados caracteres no imprimibles
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

	// Si hay más de 10 caracteres nulos o más del 30% son no imprimibles, considerarlo binario
	return nullCount > 10 || float64(nonPrintableCount)/float64(n) > 0.3
}

// showErrorMessage muestra un mensaje de error formateado
func showErrorMessage(b *strings.Builder, prefix, filePath string, width, height int) {
	// Mensaje de error más informativo pero seguro
	errorMsg := prefix
	if len(filePath) > width-20 {
		errorMsg += ": " + filePath[:width-25] + "..."
	} else {
		errorMsg += ": " + filePath
	}

	// Asegurar que el mensaje no exceda el ancho del panel
	if lipgloss.Width(errorMsg) > width {
		errorMsg = errorMsg[:width-3] + "..."
	}

	// Rellenar hasta el ancho del panel
	padding := width - lipgloss.Width(errorMsg)
	if padding > 0 {
		errorMsg += strings.Repeat(" ", padding)
	}

	b.WriteString(White.Render(errorMsg) + "\n")

	// Rellenar con líneas vacías
	for i := 1; i < height; i++ {
		b.WriteString(strings.Repeat(" ", width) + "\n")
	}
}

// showBinaryFileMessage muestra un mensaje para archivos binarios
func showBinaryFileMessage(b *strings.Builder, filePath string, size int64, width, height int) {
	// Crear un mensaje informativo
	sizeStr := formatFileSize(size)

	// Primera línea: nombre del archivo
	line1 := "Binary file"
	if lipgloss.Width(line1) > width {
		line1 = line1[:width-3] + "..."
	}
	padding1 := width - lipgloss.Width(line1)
	if padding1 > 0 {
		line1 += strings.Repeat(" ", padding1)
	}
	b.WriteString(White.Render(line1) + "\n")

	// Segunda línea: tamaño del archivo
	line2 := "Size: " + sizeStr
	padding2 := width - lipgloss.Width(line2)
	if padding2 > 0 {
		line2 += strings.Repeat(" ", padding2)
	}
	b.WriteString(White.Render(line2) + "\n")

	// Tercera línea: mensaje informativo
	line3 := "Preview not available"
	padding3 := width - lipgloss.Width(line3)
	if padding3 > 0 {
		line3 += strings.Repeat(" ", padding3)
	}
	b.WriteString(White.Render(line3) + "\n")

	// Rellenar con líneas vacías
	for i := 3; i < height; i++ {
		b.WriteString(strings.Repeat(" ", width) + "\n")
	}
}

// formatFileSize formatea el tamaño de un archivo en una cadena legible
func formatFileSize(size int64) string {
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

// sanitizeLine elimina caracteres problemáticos de una línea
func sanitizeLine(line string) string {
	// Reemplazar caracteres de control y otros caracteres problemáticos
	var result strings.Builder
	for _, r := range line {
		if r < 32 && r != '\t' && r != '\n' && r != '\r' {
			// Reemplazar caracteres de control con un espacio
			result.WriteRune(' ')
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
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

	// Obtener el selector actual
	selector := GetCurrentSelector()

	for i := start; i < end; i++ {
		item := items[i]
		fullPath := filepath.Join(directory, item)
		if item == ".." {
			fullPath = filepath.Dir(directory)
		}
		isSelected := false
		if item != ".." {
			itemKey := selector.GetSelectionKey(item)
			isSelected = selector.Selection[itemKey]
		}
		hasFocus := active && i == position

		marker := "  "
		if isSelected {
			marker = " •"
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
			b.WriteString(Yellow.Render(line) + "\n")
		} else {
			b.WriteString(Green.Render(line) + "\n")
		}
	}

	return b.String()
}

// Variable global para mantener el selector actual
var currentSelector *Selector

// Función para establecer el selector actual
func SetCurrentSelector(s *Selector) {
	currentSelector = s
}

// Función para obtener el selector actual
func GetCurrentSelector() *Selector {
	if currentSelector == nil {
		currentSelector = &Selector{
			Selection: make(map[string]bool),
		}
	}
	return currentSelector
}

// countItems cuenta el número total de elementos (archivos + subdirectorios) en un directorio
func countItems(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	return len(entries), nil
}

// countFiles cuenta el número de archivos en un directorio
func countFiles(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			count++
		}
	}
	return count, nil
}

// countSubdirs cuenta el número de subdirectorios en un directorio
func countSubdirs(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			count++
		}
	}
	return count, nil
}

// countSelected cuenta el número de archivos y directorios seleccionados
func countSelected(selector *Selector) (int, int) {
	selectedFiles := 0
	selectedDirs := 0

	for key, selected := range selector.Selection {
		if !selected {
			continue
		}

		// Verificar si es un directorio o un archivo
		fileInfo, err := os.Stat(key)
		if err != nil {
			continue
		}

		if fileInfo.IsDir() {
			selectedDirs++
		} else {
			selectedFiles++
		}
	}

	return selectedFiles, selectedDirs
}
