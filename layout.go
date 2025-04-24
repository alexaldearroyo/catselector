package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Definición de colores exactamente como en el ejemplo Python
	directoryStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))  // Verde
	fileStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))  // Blanco
	selectedItemStyle  = lipgloss.NewStyle().Reverse(true)                    // Invertido
	markedItemStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))  // Amarillo
	headerStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))  // Cyan
	activeHeaderStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("6")) // Negro sobre Cyan
	dirIconStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))  // Amarillo
	fileIconStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))  // Azul
	keyHintStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))  // Magenta
	keyHintTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))  // Blanco
	magentaTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))  // Magenta
	infoCountStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))  // Azul
)

// RenderLayout renderiza todo el layout de la aplicación
func RenderLayout(m Model) string {
	width, height := 100, 30 // Valores por defecto, idealmente obtenerlos del terminal real

	// Calcular dimensiones de los paneles
	leftPanelWidth := width / 3
	middlePanelWidth := width / 3
	rightPanelWidth := width - leftPanelWidth - middlePanelWidth

	// Separadores para altura de contenido
	headerLinesUsed := 1
	panelHeight := height - 7 // Reservar espacio para header, footer y separador

	// Renderizar header y títulos de paneles
	header, headerLinesUsed, dirDisplayItems := renderHeader(m.directory, m.initialDir, width, int(m.activePanel), m.directories, m.includeSubdirs)

	// Contenido principal
	dirPanel := renderDirectoriesPanel(m, dirDisplayItems, leftPanelWidth, panelHeight, headerLinesUsed)
	filePanel := renderFilesPanel(m, leftPanelWidth, middlePanelWidth, panelHeight, headerLinesUsed)
	previewPanel := renderPreviewPanel(m, leftPanelWidth, middlePanelWidth, rightPanelWidth, panelHeight, headerLinesUsed)

	// Separadores verticales
	verticalSeparators := ""
	for y := 0; y < panelHeight; y++ {
		verticalSeparators += lipgloss.PlaceHorizontal(width, lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top,
				lipgloss.PlaceHorizontal(leftPanelWidth, lipgloss.Right, "│"),
				lipgloss.PlaceHorizontal(middlePanelWidth, lipgloss.Right, "│"),
			))
		if y < panelHeight-1 {
			verticalSeparators += "\n"
		}
	}

	// Renderizar los contadores de elementos
	itemCounts := renderItemCounts(m, leftPanelWidth, middlePanelWidth, rightPanelWidth, headerLinesUsed)

	// Separador horizontal
	separatorY := headerLinesUsed + 2 + panelHeight
	horizontalSeparator := strings.Repeat("─", width)

	// Footer con atajos de teclado
	footer := renderKeybindings(m, width, separatorY)

	// Ensamblar todas las partes
	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		itemCounts,
		lipgloss.JoinHorizontal(lipgloss.Top, dirPanel, filePanel, previewPanel),
		horizontalSeparator,
		footer,
	)
}

func renderHeader(directory, initialDir string, width int, activePanel int, directories []string, includeSubdirs bool) (string, int, []string) {
	headerPrefix := "Directory: "
	headerFull := headerPrefix + directory

	// Calcular espacio reservado para el título y modo
	reservedSpace := min(max(len("CopyScript")+len("Mode: EXCLUDE")+6, 30), width/2)
	maxHeaderWidth := max(width-reservedSpace, 20)

	headerLines := []string{}
	headerLinesUsed := 1

	if len(headerFull) <= maxHeaderWidth {
		// Mostrar todo en una sola línea
		lastSep := strings.LastIndex(directory, string(os.PathSeparator))
		if lastSep != -1 {
			dirPrefix := directory[:lastSep+1]
			dirSuffix := directory[lastSep+1:]
			headerLines = append(headerLines, headerPrefix+dirPrefix+markedItemStyle.Render(dirSuffix))
		} else {
			headerLines = append(headerLines, headerFull)
		}
		headerLinesUsed = 1
	} else {
		// Mostrar en múltiples líneas
		headerLines = append(headerLines, "Directory:")
		parts := strings.Split(directory, string(os.PathSeparator))
		currentLine := ""
		currentPos := 0
		currentLineNum := 1

		for i, part := range parts {
			if i > 0 {
				if currentPos+1 < width {
					currentLine += string(os.PathSeparator)
					currentPos++
				} else {
					headerLines = append(headerLines, currentLine)
					currentLine = ""
					currentPos = 0
					currentLineNum++
				}
			}

			if len(part) > width-1 {
				for j := 0; j < len(part); j += width-1 {
					end := min(j+width-1, len(part))
					chunk := part[j:end]
					isLast := i == len(parts)-1
					if isLast {
						currentLine += markedItemStyle.Render(chunk)
					} else {
						currentLine += chunk
					}
					headerLines = append(headerLines, currentLine)
					currentLine = ""
					currentPos = 0
					currentLineNum++
				}
			} else {
				if currentPos+len(part) <= width {
					if i == len(parts)-1 {
						currentLine += markedItemStyle.Render(part)
					} else {
						currentLine += part
					}
					currentPos += len(part)
				} else {
					headerLines = append(headerLines, currentLine)
					currentLine = ""
					currentPos = 0
					currentLineNum++
					if i == len(parts)-1 {
						currentLine += markedItemStyle.Render(part)
					} else {
						currentLine += part
					}
					currentPos += len(part)
				}
			}
		}

		if currentLine != "" {
			headerLines = append(headerLines, currentLine)
		}

		headerLinesUsed = len(headerLines)
	}

	// Títulos de paneles
	titles := []string{"Directories [d]", "Files [f]", "Preview"}
	if activePanel == 1 {
		titles[2] = "Preview Subdirectories"
	} else if activePanel == 2 {
		titles[2] = "Preview File"
	} else {
		titles[2] = "Preview"
	}

	panelTitleLine := ""
	for i, title := range titles {
		if activePanel == i+1 {
			panelTitleLine += activeHeaderStyle.Render(title)
		} else {
			panelTitleLine += headerStyle.Render(title)
		}

		if i < 2 {
			panelTitleLine += strings.Repeat(" ", width/3-len(title))
		}
	}

	headerLines = append(headerLines, panelTitleLine)
	headerLinesUsed += 1

	// Panel de directorios a mostrar
	var dirDisplayItems []string
	if directory == initialDir {
		dirDisplayItems = append([]string{"."}, directories...)
	} else {
		dirDisplayItems = append([]string{"..", "."}, directories...)
	}

	// Colocar el título de la app en la esquina superior derecha
	appTitle := "CopyScript"
	header := headerLines[0] + strings.Repeat(" ", width-len(headerLines[0])-len(appTitle)) + markedItemStyle.Render(appTitle)

	// Añadir resto de líneas del header
	for i := 1; i < len(headerLines); i++ {
		header += "\n" + headerLines[i]
	}

	// Incluir info sobre subdirectorios
	subdirMode := "Subdirectories: "
	includedText := "Included"
	if !includeSubdirs {
		includedText = "Not included"
	}

	// Creamos un string con toda la info de subdirectorios
	subdirInfo := keyHintTextStyle.Render(subdirMode) + magentaTextStyle.Render(includedText)
	header += "\n" + subdirInfo

	return header, headerLinesUsed + 1, dirDisplayItems
}

func renderDirectoriesPanel(m Model, dirDisplayItems []string, leftPanelWidth int, panelHeight int, headerLinesUsed int) string {
	dirContent := ""

	// Mostrar cada elemento del directorio
	for i, item := range dirDisplayItems {
		fullPath := item
		if item != ".." {
			fullPath = fmt.Sprintf("%s%c%s", m.directory, os.PathSeparator, item)
		} else {
			fullPath = fmt.Sprintf("%s", strings.TrimSuffix(m.directory, string(os.PathSeparator)))
		}

		isSelected := contains(m.selected, fullPath)

		// Verificar si este item tiene el foco
		hasFocus := (m.activePanel == directoriesPanel && i == m.position)

		icon := "📁"
		marker := "  "
		if isSelected {
			if m.includeSubdirs {
				marker = "* "
			} else {
				marker = "• "
			}
		}

		content := fmt.Sprintf("%s%s", marker, item)

		// Aplicar estilo según el estado
		if hasFocus {
			dirContent += selectedItemStyle.Render(icon + " " + content) + "\n"
		} else {
			if isSelected {
				dirContent += dirIconStyle.Render(icon) + " " + markedItemStyle.Render(content) + "\n"
			} else {
				dirContent += dirIconStyle.Render(icon) + " " + directoryStyle.Render(content) + "\n"
			}
		}
	}

	// TODO: Añadir scrollbar para directorios si es necesario

	return dirContent
}

func renderFilesPanel(m Model, leftPanelWidth, middlePanelWidth, panelHeight, headerLinesUsed int) string {
	fileContent := ""

	for i, item := range m.files {
		// Determinar directorio actual para la ruta completa
		fileDir := m.directory
		if m.currentPreviewDirectory != "" {
			fileDir = m.currentPreviewDirectory
		}

		fullPath := fmt.Sprintf("%s%c%s", fileDir, os.PathSeparator, item)

		// Determinar si está seleccionado
		isSelected := contains(m.selected, fullPath)

		marker := "*"
		if !isSelected {
			marker = " "
		}

		icon := GetFileIcon(fullPath)

		// Verificar si tiene el foco
		hasFocus := (m.activePanel == filesPanel && i == m.position)

		if hasFocus {
			fileContent += selectedItemStyle.Render(icon + " " + marker + item) + "\n"
		} else {
			if isSelected {
				fileContent += fileIconStyle.Render(icon) + " " + markedItemStyle.Render(marker+item) + "\n"
			} else {
				fileContent += fileIconStyle.Render(icon) + " " + fileStyle.Render(marker+item) + "\n"
			}
		}
	}

	// TODO: Añadir scrollbar para archivos si es necesario

	return fileContent
}

func renderPreviewPanel(m Model, leftPanelWidth, middlePanelWidth, rightPanelWidth, panelHeight, headerLinesUsed int) string {
	previewContent := ""

	for i, line := range m.previewContent {
		actualLine := i
		highlight := (m.activePanel == previewPanel && actualLine == m.previewPosition)

		if strings.HasPrefix(line, "📁 ") {
			icon := "📁"
			content := strings.TrimPrefix(line, "📁 ")
			markerPresent := strings.HasPrefix(content, "* ")
			contentClean := content
			if markerPresent {
				contentClean = strings.TrimPrefix(content, "* ")
			}

			var fullPath string
			if m.currentPreviewDirectory != "" {
				fullPath = fmt.Sprintf("%s%c%s", m.currentPreviewDirectory, os.PathSeparator, contentClean)
			} else {
				fullPath = fmt.Sprintf("%s%c%s", m.directory, os.PathSeparator, contentClean)
			}

			isSelected := contains(m.selected, fullPath)

			style := fileStyle
			if isSelected {
				style = markedItemStyle
			}
			if highlight {
				style = selectedItemStyle
			}

			previewContent += dirIconStyle.Render(icon) + " " + style.Render(content) + "\n"
		} else if containsAnyPrefix(line, []string{"📄", "📜", "📝", "⚙️", "🖼️", "🎵", "🎬", "📦", "📕", "📘", "📗", "📙", "🚀", "🌐", "🐙"}) {
			icon := string([]rune(line)[0])
			rest := strings.TrimPrefix(line, icon+" ")
			isSelected := strings.HasPrefix(rest, "*")

			style := fileStyle
			if isSelected {
				style = markedItemStyle
			}
			if highlight {
				style = selectedItemStyle
			}

			previewContent += fileIconStyle.Render(icon) + " " + style.Render(rest) + "\n"
		} else {
			style := lipgloss.NewStyle()
			if strings.HasPrefix(strings.TrimSpace(line), "*") {
				style = markedItemStyle
			}
			if highlight {
				style = selectedItemStyle
			}

			previewContent += style.Render(line) + "\n"
		}
	}

	// TODO: Añadir scrollbar para preview si es necesario

	return previewContent
}

func renderItemCounts(m Model, leftPanelWidth, middlePanelWidth, rightPanelWidth, headerLinesUsed int) string {
	// Aquí implementaríamos la lógica de contadores como en el original
	// Por simplicidad, mostramos valores de ejemplo
	dirCountDisplay := infoCountStyle.Render("10 items")
	fileCountDisplay := infoCountStyle.Render("5 files")
	previewCountDisplay := infoCountStyle.Render("3 subdirs")

	// TODO: Calcular contadores reales basados en los datos del modelo

	itemCounts := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.PlaceHorizontal(leftPanelWidth, lipgloss.Left, dirCountDisplay),
		lipgloss.PlaceHorizontal(middlePanelWidth, lipgloss.Left, fileCountDisplay),
		lipgloss.PlaceHorizontal(rightPanelWidth, lipgloss.Left, previewCountDisplay),
	)

	return itemCounts
}

func renderKeybindings(m Model, width int, separatorY int) string {
	keyBindings := []string{
		"k/j: Up & Down", "s: Select", "o: Export & Open",
		"h/l: Go into", "a: Select all", "c: Copy to clipboard",
		"Esc/h: Go Back", "i: Include subdirs", "q: Quit",
	}

	keyBindingsText := ""

	for i, binding := range keyBindings {
		keyEnd := strings.Index(binding, ":")
		if keyEnd != -1 {
			key := binding[:keyEnd]
			desc := binding[keyEnd:]

			keyBindingsText += keyHintStyle.Render(key) + keyHintTextStyle.Render(desc)

			if i < len(keyBindings)-1 {
				if i % 3 == 2 {
					keyBindingsText += "\n"
				} else {
					keyBindingsText += strings.Repeat(" ", (width/3)-len(binding))
				}
			}
		}
	}

	// Creamos la variable pero la usamos directamente
	footer := keyBindingsText

	if m.statusMessage != "" && time.Now().Unix() - m.statusTime < 2 {
		footer += "\n" + markedItemStyle.Render(m.statusMessage)
	}

	return footer
}

// Funciones auxiliares

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsAnyPrefix(s string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix+" ") {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
