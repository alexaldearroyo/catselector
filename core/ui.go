package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

// UpdateDirectoriesPanel actualiza el contenido del panel de directorios
func (s *Selector) updateDirectoriesPanel() {
	s.directoriesPanel.Clear()

	// Determinar elementos a mostrar
	var dirItems []string
	if s.directory == s.initialDirectory {
		dirItems = append([]string{"."}, s.directories...)
	} else {
		dirItems = append([]string{"..", "."}, s.directories...)
	}

	// Añadir elementos al panel
	for i, item := range dirItems {
		// Determinar icono y si está seleccionado
		icon := s.getItemIcon(item, true)
		fullPath := s.directory
		if item != "." {
			if item == ".." {
				fullPath = filepath.Dir(s.directory)
			} else {
				fullPath = filepath.Join(s.directory, item)
			}
		}
		isSelected := s.selected[fullPath]

		// Determinar marcador
		marker := "  "
		if isSelected {
			if s.includeSubdirectories {
				marker = "* "
			} else {
				marker = "• "
			}
		}

		// Crear texto del item con icono y marcador
		itemText := fmt.Sprintf("%s %s%s", icon, marker, item)

		if i == s.position && s.activePanel == 1 {
			// Agregarlo como seleccionado
			s.directoriesPanel.AddItem(itemText, "", 0, nil).SetSelectedBackgroundColor(tcell.ColorBlue)
		} else {
			// Agregar como item normal
			s.directoriesPanel.AddItem(itemText, "", 0, nil)
		}
	}
}

// UpdateFilesPanel actualiza el contenido del panel de archivos
func (s *Selector) updateFilesPanel() {
	s.filesPanel.Clear()

	// Obtener directorio para archivos
	fileDir := s.currentPreviewDirectory
	if fileDir == "" {
		fileDir = s.directory
	}

	// Intentar leer archivos del directorio
	files := s.files
	if len(files) == 0 && fileDir != "" {
		fileEntries, err := os.ReadDir(fileDir)
		if err == nil {
			for _, entry := range fileEntries {
				if !entry.IsDir() {
					files = append(files, entry.Name())
				}
			}
		}
	}

	// Añadir archivos al panel
	for i, file := range files {
		// Determinar icono y si está seleccionado
		icon := GetFileIcon(filepath.Join(fileDir, file))
		fullPath := filepath.Join(fileDir, file)

		// Verificar selección
		isSelected := false
		if _, found := s.manuallyDeselected[fullPath]; found {
			isSelected = false
		} else if _, found := s.selected[fullPath]; found {
			isSelected = true
		} else {
			// Comprobar si está seleccionado por un directorio padre
			for path := range s.selected {
				info, err := os.Stat(path)
				if err == nil && info.IsDir() {
					if s.includeSubdirectories {
						if strings.HasPrefix(fullPath, path+string(filepath.Separator)) {
							isSelected = true
							break
						}
					} else {
						if filepath.Dir(fullPath) == path {
							isSelected = true
							break
						}
					}
				}
			}
		}

		// Determinar marcador
		marker := " "
		if isSelected {
			marker = "*"
		}

		// Crear texto del item con icono y marcador
		itemText := fmt.Sprintf("%s %s%s", icon, marker, file)

		if i == s.position && s.activePanel == 2 {
			// Agregarlo como seleccionado
			s.filesPanel.AddItem(itemText, "", 0, nil).SetSelectedBackgroundColor(tcell.ColorBlue)
		} else {
			// Agregar como item normal
			s.filesPanel.AddItem(itemText, "", 0, nil)
		}
	}
}

// UpdatePreviewPanel actualiza el contenido del panel de vista previa
func (s *Selector) updatePreviewPanel() {
	s.previewPanel.Clear()

	// Obtener elemento seleccionado para vista previa
	var selectedItem string
	if s.activePanel == 1 {
		if len(s.filtered) > 0 && s.position < len(s.filtered) {
			selectedItem = s.filtered[s.position]
		}
	} else { // Panel Files
		if len(s.files) > 0 && s.position < len(s.files) {
			selectedItem = s.files[s.position]
		}
	}

	// Actualizar contenido de vista previa
	if selectedItem != "" {
		s.previewContent = s.getPreviewContent(selectedItem)
	} else {
		s.previewContent = []string{"Ningún elemento seleccionado"}
	}

	// Mostrar contenido de vista previa
	var displayText strings.Builder
	start := s.previewWindowStart
	end := min(start+30, len(s.previewContent)) // Mostrar máximo 30 líneas

	for i := start; i < end; i++ {
		line := s.previewContent[i]

		if i == s.previewPosition && s.activePanel == 3 {
			// Destacar línea seleccionada
			displayText.WriteString("[::b]")
			displayText.WriteString(line)
			displayText.WriteString("[::-]")
		} else {
			// Aplicar colores según contenido
			if strings.HasPrefix(line, "📁 ") {
				if strings.Contains(line, "*") {
					displayText.WriteString("[yellow]")
				} else {
					displayText.WriteString("[green]")
				}
			} else if strings.HasPrefix(line, "📄 ") ||
				      strings.HasPrefix(line, "📜 ") ||
				      strings.HasPrefix(line, "📝 ") {
				if strings.Contains(line, "*") {
					displayText.WriteString("[yellow]")
				} else {
					displayText.WriteString("[white]")
				}
			} else if strings.HasPrefix(line, "Error") {
				displayText.WriteString("[red]")
			}

			displayText.WriteString(line)
			displayText.WriteString("[white]")
		}

		displayText.WriteString("\n")
	}

	s.previewPanel.SetText(displayText.String())
}

// UpdateHeaderPanel actualiza el panel de cabecera
func (s *Selector) updateHeaderPanel() {
	var header strings.Builder

	// Mostrar directorio actual
	header.WriteString("[yellow]Directorio: [white]")

	// Dividir la ruta en partes para facilitar la lectura
	parts := strings.Split(filepath.ToSlash(s.directory), "/")
	displayPath := filepath.Join(parts...)

	header.WriteString(displayPath)
	header.WriteString("\n")

	// Mostrar modo de subdirectorios
	header.WriteString("[magenta]Subdirectorios: [white]")
	if s.includeSubdirectories {
		header.WriteString("[yellow]Incluidos[white]")
	} else {
		header.WriteString("[white]No incluidos[white]")
	}

	// Mostrar información de selección
	var selectedFiles, selectedDirs int
	countSelected(s, &selectedFiles, &selectedDirs)

	header.WriteString(" | [magenta]Seleccionados: [yellow]")
	header.WriteString(strconv.Itoa(selectedFiles))
	header.WriteString(" Archivos[white], [yellow]")
	header.WriteString(strconv.Itoa(selectedDirs))
	header.WriteString(" Directorios[white]")

	// Mostrar app title
	header.WriteString("[right][yellow]CopyScript[white]")

	s.headerPanel.SetText(header.String())
}

// UpdateKeybindingsPanel actualiza el panel de teclas de acceso rápido
func (s *Selector) updateKeybindingsPanel() {
	var keybindings strings.Builder

	// Teclas comunes
	keys := []struct {
		key  string
		desc string
	}{
		{"[magenta]k/j[white]", ": Arriba/Abajo"},
		{"[magenta]s[white]", ": Seleccionar"},
		{"[magenta]o[white]", ": Exportar y Abrir"},
		{"[magenta]h/l[white]", ": Navegar"},
		{"[magenta]a[white]", ": Seleccionar todo"},
		{"[magenta]c[white]", ": Copiar al portapapeles"},
		{"[magenta]Esc/h[white]", ": Volver"},
		{"[magenta]i[white]", ": Incluir subdirs"},
		{"[magenta]q[white]", ": Salir"},
	}

	// Primera línea de teclas
	for i := 0; i < 3 && i < len(keys); i++ {
		keybindings.WriteString(keys[i].key)
		keybindings.WriteString(keys[i].desc)
		keybindings.WriteString("  ")
	}
	keybindings.WriteString("\n")

	// Segunda línea de teclas
	for i := 3; i < 6 && i < len(keys); i++ {
		keybindings.WriteString(keys[i].key)
		keybindings.WriteString(keys[i].desc)
		keybindings.WriteString("  ")
	}
	keybindings.WriteString("\n")

	// Tercera línea de teclas
	for i := 6; i < len(keys); i++ {
		keybindings.WriteString(keys[i].key)
		keybindings.WriteString(keys[i].desc)
		keybindings.WriteString("  ")
	}

	s.keybindingsPanel.SetText(keybindings.String())
}

// UpdateStatusBar actualiza la barra de estado
func (s *Selector) updateStatusBar() {
	// Si estamos en modo búsqueda, mostrar prompt
	if s.searchMode {
		s.statusBar.SetText(fmt.Sprintf("Buscar: %s", s.searchTerm))
		return
	}

	// Si hay un mensaje de estado y no ha expirado, mostrarlo
	if s.statusMessage != "" && float64(time.Now().Unix())-s.statusTime < 2 {
		s.statusBar.SetText(s.statusMessage)
	} else {
		// De lo contrario, limpiar
		s.statusBar.SetText("")
		s.statusMessage = ""
	}
}
