package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alexadler/copycat/export"
	"github.com/gdamore/tcell/v2"
)

// HandleInput maneja los eventos de teclado
func (s *Selector) handleInput(event *tcell.EventKey) *tcell.EventKey {
	// Panel de vista previa activo
	if s.activePanel == 3 {
		return s.handlePreviewNavigation(event)
	}

	if s.searchMode {
		return s.handleSearchMode(event)
	}

	// Modo normal
	return s.handleNormalMode(event)
}

// HandleNormalMode maneja la entrada en modo normal
func (s *Selector) handleNormalMode(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	rune := event.Rune()

	switch {
	case key == tcell.KeyEscape:
		s.handleEscapeBack()
	case rune == 'q':
		s.app.Stop()
	case rune == '/':
		s.searchMode = true
		s.searchTerm = ""
		s.app.SetFocus(s.statusBar)
		s.statusMessage = "Buscar: "
		s.statusTime = float64(time.Now().Unix())
	case key == tcell.KeyUp || rune == 'k':
		s.position = max(0, s.position-1)
		s.updateUI()
	case key == tcell.KeyDown || rune == 'j':
		s.position = min(len(s.filtered)-1, s.position+1)
		s.updateUI()
	case rune == 'd':
		s.switchToDirectoriesPanel()
	case rune == 'f':
		s.switchToFilesPanel()
	case rune == 'p':
		s.activePanel = 3
		s.updateUI()
	case rune == 'i':
		s.includeSubdirectories = !s.includeSubdirectories
		s.updateUI()
	case rune == 's':
		s.handleSelectOnly()
		s.updateUI()
	case rune == 'a':
		s.toggleSelectAll()
		s.updateUI()
	case rune == 'o':
		s.exportAndOpen()
	case rune == 'c':
		s.exportAndCopyToClipboard()
	case key == tcell.KeyEnter:
		if s.activePanel == 1 && len(s.filtered) > 0 {
			item := s.filtered[s.position]
			if item == ".." {
				s.handleGoBackDirectory()
				return nil
			}
		}
		s.handleEnterDirectory()
	case key == tcell.KeyTab:
		if s.activePanel == 1 {
			s.switchToFilesPanel()
		} else if s.activePanel == 2 {
			s.switchToDirectoriesPanel()
		}
	case rune == 'h':
		s.handleEscapeBack()
	case rune == 'l':
		s.handleEnterDirectory()
	}

	return nil
}

// HandleSearchMode maneja la entrada en modo de búsqueda
func (s *Selector) handleSearchMode(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	rune := event.Rune()

	switch {
	case key == tcell.KeyEscape:
		s.searchMode = false
		s.updateFiltered()
		s.updateUI()
	case key == tcell.KeyBackspace || key == tcell.KeyBackspace2:
		if len(s.searchTerm) > 0 {
			s.searchTerm = s.searchTerm[:len(s.searchTerm)-1]
			s.updateFiltered()
			s.updateUI()
		}
	case key == tcell.KeyEnter:
		s.searchMode = false
		s.updateFiltered()
		s.updateUI()
	default:
		if key != tcell.KeyRune {
			break
		}
		s.searchTerm += string(rune)
		s.updateFiltered()
		s.updateUI()
	}

	return nil
}

// HandlePreviewNavigation maneja la navegación en el panel de vista previa
func (s *Selector) handlePreviewNavigation(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	rune := event.Rune()

	switch {
	case key == tcell.KeyUp || rune == 'k':
		s.previewPosition = max(0, s.previewPosition-1)
		if s.previewPosition < s.previewWindowStart {
			s.previewWindowStart = s.previewPosition
		}
	case key == tcell.KeyDown || rune == 'j':
		s.previewPosition = min(len(s.previewContent)-1, s.previewPosition+1)
		// Ajustar previewWindowStart si es necesario
		panelHeight := 20 // Valor predeterminado, debería calcularse dinámicamente
		if s.previewPosition >= s.previewWindowStart+panelHeight {
			s.previewWindowStart = s.previewPosition - panelHeight + 1
		}
	case rune == 'd':
		s.activePanel = 1
		s.updateUI()
	case rune == 'f':
		s.activePanel = 2
		s.updateUI()
	case key == tcell.KeyEscape:
		s.activePanel = 1
		s.updateUI()
	}

	return nil
}

// SwitchToDirectoriesPanel cambia al panel de directorios
func (s *Selector) switchToDirectoriesPanel() {
	s.activePanel = 1

	// Construir elementos de directorio
	var dirItems []string
	if s.directory == s.initialDirectory {
		dirItems = append([]string{"."}, s.directories...)
	} else {
		dirItems = append([]string{"..", "."}, s.directories...)
	}

	prevDir := s.currentPreviewDirectory
	if prevDir != "" {
		// Si el directorio de vista previa es el directorio actual
		if filepath.Clean(prevDir) == filepath.Clean(s.directory) {
			s.position = indexOf(dirItems, ".")
		} else if filepath.Clean(prevDir) == filepath.Clean(filepath.Dir(s.directory)) {
			// Si el directorio de vista previa es el padre del directorio actual
			s.position = indexOf(dirItems, "..")
		} else {
			// Para subdirectorios, coincidir por nombre base
			dirName := filepath.Base(prevDir)
			idx := indexOf(dirItems, dirName)
			if idx >= 0 {
				s.position = idx
			} else {
				s.position = 0
			}
		}
	} else {
		s.position = 0
	}

	s.windowStart = 0
	s.updateFiltered()
	s.updateUI()
}

// SwitchToFilesPanel cambia al panel de archivos
func (s *Selector) switchToFilesPanel() {
	if s.activePanel == 1 && len(s.filtered) > 0 && s.position < len(s.filtered) {
		currentItem := s.filtered[s.position]
		if currentItem != ".." && currentItem != "." {
			s.currentPreviewDirectory = filepath.Join(s.directory, currentItem)
		} else if currentItem == ".." {
			s.currentPreviewDirectory = filepath.Dir(s.directory)
		} else {
			s.currentPreviewDirectory = s.directory
		}
	}

	s.activePanel = 2
	s.position = 0
	s.windowStart = 0
	s.updateFiltered()
	s.updateUI()
}

// ToggleSelectAll selecciona o deselecciona todos los elementos
func (s *Selector) toggleSelectAll() {
	allItems := make(map[string]bool)

	if s.activePanel == 1 {
		for _, item := range s.filtered {
			if item != ".." && item != "." {
				fullPath := filepath.Join(s.directory, item)
				allItems[fullPath] = true
			}
		}
	} else {
		fileDir := s.currentPreviewDirectory
		if fileDir == "" {
			fileDir = s.directory
		}
		for _, item := range s.files {
			fullPath := filepath.Join(fileDir, item)
			allItems[fullPath] = true
		}
	}

	// Verificar si todos los elementos ya están seleccionados
	allSelected := true
	for path := range allItems {
		if !s.selected[path] {
			allSelected = false
			break
		}
	}

	// Seleccionar o deseleccionar todos
	if allSelected {
		for path := range allItems {
			delete(s.selected, path)
		}
	} else {
		for path := range allItems {
			s.selected[path] = true
		}
	}

	s.updateFiltered()
}

// HandleSelectOnly selecciona o deselecciona el elemento actual
func (s *Selector) handleSelectOnly() {
	if len(s.filtered) == 0 {
		return
	}

	item := s.filtered[s.position]
	if item == ".." {
		return
	}

	var fullPath string
	if item == "." {
		fullPath = s.directory
	} else {
		fullPath = filepath.Join(s.directory, item)
	}

	if isDir(fullPath) {
		// Verificar si ya está seleccionado
		alreadySelected := s.selected[fullPath]
		if !alreadySelected {
			// Seleccionar directorios y su contenido
			s.selected[fullPath] = true

			if s.includeSubdirectories {
				// Seleccionar recursivamente
				filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return nil
					}
					s.selected[path] = true
					return nil
				})
			} else {
				// Seleccionar solo archivos inmediatos
				files, err := os.ReadDir(fullPath)
				if err == nil {
					for _, file := range files {
						if !file.IsDir() {
							filePath := filepath.Join(fullPath, file.Name())
							s.selected[filePath] = true
						}
					}
				}
			}
		} else {
			// Deseleccionar directorio y su contenido
			delete(s.selected, fullPath)

			// Deseleccionar contenido
			if s.includeSubdirectories {
				for path := range s.selected {
					if strings.HasPrefix(path, fullPath) {
						delete(s.selected, path)
					}
				}
			} else {
				// Deseleccionar solo archivos inmediatos
				files, err := os.ReadDir(fullPath)
				if err == nil {
					for _, file := range files {
						if !file.IsDir() {
							filePath := filepath.Join(fullPath, file.Name())
							delete(s.selected, filePath)
						}
					}
				}
			}
		}
	} else if s.activePanel == 2 {
		fileDir := s.currentPreviewDirectory
		if fileDir == "" {
			fileDir = s.directory
		}
		fullPath = filepath.Join(fileDir, item)

		// Alternar selección
		if s.selected[fullPath] {
			delete(s.selected, fullPath)
			s.manuallyDeselected[fullPath] = true
		} else {
			s.selected[fullPath] = true
			delete(s.manuallyDeselected, fullPath)
		}
	}

	s.updateFiltered()
}

// HandleGoBackDirectory maneja la acción de volver al directorio padre
func (s *Selector) handleGoBackDirectory() {
	if s.activePanel == 1 {
		currentItem := ""
		if len(s.filtered) > 0 && s.position < len(s.filtered) {
			currentItem = s.filtered[s.position]
		}
		shouldGoBack := currentItem != "."

		if shouldGoBack {
			parentDir := filepath.Dir(s.directory)
			if parentDir != s.directory && isSubdirOrSame(parentDir, s.initialDirectory) {
				s.directory = parentDir
				s.loadFiles()
				s.position = 0
				s.windowStart = 0
			}
		}
	}
}

// HandleEnterDirectory maneja la acción de entrar en un directorio
func (s *Selector) handleEnterDirectory() {
	if s.activePanel == 1 && len(s.filtered) > 0 {
		item := s.filtered[s.position]
		if item != ".." && item != "." {
			fullPath := filepath.Join(s.directory, item)
			if isDir(fullPath) {
				s.directory = fullPath
				s.loadFiles()
				s.position = 0
				s.windowStart = 0
			}
		}
	}
}

// HandleEscapeBack maneja la acción de volver atrás con la tecla Escape
func (s *Selector) handleEscapeBack() {
	if s.directory != s.initialDirectory {
		currentSubdir := filepath.Base(s.directory)
		parentDir := filepath.Dir(s.directory)

		if isSubdirOrSame(parentDir, s.initialDirectory) {
			s.directory = parentDir
			s.currentPreviewDirectory = ""
			s.loadFiles()

			// Buscar la posición del subdirectorio anterior
			for i, dir := range s.filtered {
				if dir == currentSubdir {
					s.position = i
					break
				}
			}

			if s.position < 0 || s.position >= len(s.filtered) {
				s.position = 0
			}

			s.windowStart = 0
		}
	}
}

// ExportAndOpen exporta los archivos seleccionados y los abre
func (s *Selector) exportAndOpen() {
	var selectedSlice []string
	var excludedSlice []string

	for path := range s.selected {
		selectedSlice = append(selectedSlice, path)
	}

	for path := range s.manuallyDeselected {
		excludedSlice = append(excludedSlice, path)
	}

	outputFile, err := export.GenerateTextFile(
		selectedSlice,
		excludedSlice,
		s.includeSubdirectories,
		s.initialDirectory,
		s.directory,
	)

	if err != nil {
		s.showStatusMessage(fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	// Abrir archivo generado
	if err := OpenTextFile(outputFile); err != nil {
		s.showStatusMessage(fmt.Sprintf("Error abriendo archivo: %s", err.Error()))
	} else {
		s.showStatusMessage(fmt.Sprintf("Archivo generado: %s", outputFile))
	}
}

// ExportAndCopyToClipboard exporta los archivos seleccionados y copia al portapapeles
func (s *Selector) exportAndCopyToClipboard() {
	var selectedSlice []string
	var excludedSlice []string

	for path := range s.selected {
		selectedSlice = append(selectedSlice, path)
	}

	for path := range s.manuallyDeselected {
		excludedSlice = append(excludedSlice, path)
	}

	outputFile, err := export.GenerateTextFile(
		selectedSlice,
		excludedSlice,
		s.includeSubdirectories,
		s.initialDirectory,
		s.directory,
	)

	if err != nil {
		s.showStatusMessage(fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	// Leer contenido del archivo
	content, err := os.ReadFile(outputFile)
	if err != nil {
		s.showStatusMessage(fmt.Sprintf("Error leyendo archivo: %s", err.Error()))
		return
	}

	// Copiar al portapapeles según el sistema operativo
	var success bool
	switch runtime.GOOS {
	case "darwin": // macOS
		cmd := exec.Command("pbcopy")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			s.showStatusMessage("Error: No se pudo copiar al portapapeles")
			return
		}

		if err := cmd.Start(); err != nil {
			s.showStatusMessage("Error: No se pudo iniciar pbcopy")
			return
		}

		if _, err := stdin.Write(content); err != nil {
			s.showStatusMessage("Error: No se pudo escribir en pbcopy")
			return
		}

		if err := stdin.Close(); err != nil {
			s.showStatusMessage("Error: No se pudo cerrar stdin")
			return
		}

		if err := cmd.Wait(); err != nil {
			s.showStatusMessage("Error: pbcopy falló")
			return
		}

		success = true

	case "windows": // Windows
		// En un entorno real, usaríamos la API de Windows o una biblioteca para esto
		cmd := exec.Command("clip")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			s.showStatusMessage("Error: No se pudo copiar al portapapeles")
			return
		}

		if err := cmd.Start(); err != nil {
			s.showStatusMessage("Error: No se pudo iniciar clip")
			return
		}

		if _, err := stdin.Write(content); err != nil {
			s.showStatusMessage("Error: No se pudo escribir en clip")
			return
		}

		if err := stdin.Close(); err != nil {
			s.showStatusMessage("Error: No se pudo cerrar stdin")
			return
		}

		if err := cmd.Wait(); err != nil {
			s.showStatusMessage("Error: clip falló")
			return
		}

		success = true

	default: // Linux/Unix
		// Intentar xclip primero
		cmd := exec.Command("xclip", "-selection", "clipboard")
		stdin, err := cmd.StdinPipe()
		if err == nil {
			if cmd.Start() == nil {
				if _, err := stdin.Write(content); err == nil {
					if stdin.Close() == nil {
						if cmd.Wait() == nil {
							success = true
						}
					}
				}
			}
		}

		// Si xclip falla, intentar xsel
		if !success {
			cmd = exec.Command("xsel", "--clipboard", "--input")
			stdin, err := cmd.StdinPipe()
			if err == nil {
				if cmd.Start() == nil {
					if _, err := stdin.Write(content); err == nil {
						if stdin.Close() == nil {
							if cmd.Wait() == nil {
								success = true
							}
						}
					}
				}
			}
		}

		// Si todo falla, intentar wl-copy para Wayland
		if !success {
			cmd = exec.Command("wl-copy")
			stdin, err := cmd.StdinPipe()
			if err == nil {
				if cmd.Start() == nil {
					if _, err := stdin.Write(content); err == nil {
						if stdin.Close() == nil {
							if cmd.Wait() == nil {
								success = true
							}
						}
					}
				}
			}
		}
	}

	// Eliminar el archivo temporal
	os.Remove(outputFile)

	// Calcular archivos seleccionados
	var selectedFiles, selectedDirs int
	countSelected(s, &selectedFiles, &selectedDirs)

	// Mostrar mensaje
	if success {
		s.showStatusMessage(fmt.Sprintf("%d archivos copiados al portapapeles", selectedFiles))
	} else {
		s.showStatusMessage("Error: No se pudo copiar al portapapeles")
	}
}

// ShowStatusMessage muestra un mensaje en la barra de estado
func (s *Selector) showStatusMessage(message string) {
	s.statusMessage = message
	s.statusTime = float64(time.Now().Unix())
	s.statusBar.SetText(message)
}

// UpdateUI actualiza la interfaz de usuario
func (s *Selector) updateUI() {
	// Actualizar título según el panel activo
	switch s.activePanel {
	case 1:
		s.directoriesPanel.SetTitle("[cyan]Directorios [d][white]")
		s.filesPanel.SetTitle("Archivos [f]")
		s.previewPanel.SetTitle("Vista previa")
	case 2:
		s.directoriesPanel.SetTitle("Directorios [d]")
		s.filesPanel.SetTitle("[cyan]Archivos [f][white]")
		s.previewPanel.SetTitle("Vista previa")
	case 3:
		s.directoriesPanel.SetTitle("Directorios [d]")
		s.filesPanel.SetTitle("Archivos [f]")
		s.previewPanel.SetTitle("[cyan]Vista previa[white]")
	}

	// Actualizar panel de directorios
	s.updateDirectoriesPanel()

	// Actualizar panel de archivos
	s.updateFilesPanel()

	// Actualizar panel de vista previa
	s.updatePreviewPanel()

	// Actualizar panel de cabecera
	s.updateHeaderPanel()

	// Actualizar panel de teclas
	s.updateKeybindingsPanel()

	// Actualizar barra de estado
	s.updateStatusBar()
}

// Función auxiliar para contar seleccionados
func countSelected(s *Selector, selectedFiles, selectedDirs *int) {
	*selectedFiles = 0
	*selectedDirs = 0

	for path := range s.selected {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		if info.IsDir() {
			(*selectedDirs)++
			if s.includeSubdirectories {
				filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
					if err != nil {
						return nil
					}
					if !info.IsDir() && !s.manuallyDeselected[p] {
						(*selectedFiles)++
					} else if info.IsDir() && p != path {
						(*selectedDirs)++
					}
					return nil
				})
			} else {
				files, err := os.ReadDir(path)
				if err == nil {
					for _, file := range files {
						if !file.IsDir() {
							fullPath := filepath.Join(path, file.Name())
							if !s.manuallyDeselected[fullPath] {
								(*selectedFiles)++
							}
						}
					}
				}
			}
		} else {
			(*selectedFiles)++
		}
	}
}

// Funciones auxiliares

// isSubdirOrSame comprueba si subdir es un subdirectorio de o igual a parent
func isSubdirOrSame(subdir, parent string) bool {
	subdir = filepath.Clean(subdir)
	parent = filepath.Clean(parent)
	return subdir == parent || strings.HasPrefix(subdir, parent+string(filepath.Separator))
}

// max devuelve el mayor de dos enteros
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min devuelve el menor de dos enteros
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// indexOf devuelve el índice de un elemento en un slice
func indexOf(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1
}
