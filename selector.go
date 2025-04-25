package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Comando personalizado para copiar al portapapeles
type CopyToClipboardCmd struct {
	Content string
}

// Crear función para generar el comando
func CopyToClipboard(content string) tea.Cmd {
	return func() tea.Msg {
		return CopyToClipboardCmd{Content: content}
	}
}

type panel int

const (
	directoriesPanel panel = iota + 1
	filesPanel
	previewPanel
)

type Model struct {
	directory             string
	initialDir            string
	directories           []string
	files                 []string
	position              int
	windowStart           int
	previewWindowStart    int
	previewPosition       int
	activePanel           panel
	previewContent        []string
	selected              []string
	manually_deselected   []string
	includeSubdirs        bool
	currentPreviewDirectory string
	statusMessage         string
	statusTime            int64
	width                 int
	height                int
}

// View implements tea.Model.
func (m Model) View() string {
	return RenderLayout(m)
}

func NewSelector(initialDir string) Model {
	model := Model{
		directory:            initialDir,
		initialDir:           initialDir,
		activePanel:          directoriesPanel,
		selected:             []string{},
		manually_deselected:  []string{},
		includeSubdirs:       false,
		statusMessage:        "",
		statusTime:           0,
		width:                100, // Valor por defecto hasta recibir el tamaño real
		height:               30,
	}
	model.loadDirectory()
	return model
}

func (m *Model) loadDirectory() {
	entries, err := os.ReadDir(m.directory)
	if err != nil {
		m.directories = nil
		m.files = nil
		return
	}

	dirs := []string{}
	files := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		} else {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(dirs)
	sort.Strings(files)

	m.directories = dirs
	m.files = files
	m.position = 0
	m.windowStart = 0
	m.updatePreview()
}

func (m *Model) updatePreview() {
	m.previewContent = []string{}

	if m.activePanel == directoriesPanel && len(m.directories) > 0 && m.position < len(m.directories) {
		// Actualizar el directorio de vista previa
		name := m.directories[m.position]
		subPath := filepath.Join(m.directory, name)
		m.currentPreviewDirectory = subPath

		// Obtener elementos para mostrar
		entries, err := os.ReadDir(subPath)
		if err == nil {
			for _, e := range entries {
				if e.IsDir() {
					fullPath := filepath.Join(subPath, e.Name())
					isSelected := containsPath(m.selected, fullPath) || isInSelectedDir(m.selected, fullPath, m.includeSubdirs)
					marker := ""
					if isSelected {
						if m.includeSubdirs {
							marker = "* "
						} else {
							marker = "• "
						}
					} else {
						marker = "  "
					}
					m.previewContent = append(m.previewContent, "📁 "+marker+e.Name())
				}
			}
		}
	} else if m.activePanel == filesPanel && len(m.files) > 0 && m.position < len(m.files) {
		name := m.files[m.position]
		filePath := filepath.Join(m.directory, name)
		// Intentar mostrar contenido para archivos de texto
		if isTextFile(filePath) {
			data, err := os.ReadFile(filePath)
			if err == nil {
				lines := strings.Split(string(data), "\n")
				for i, line := range lines {
					if i < 50 { // Limitamos a 50 líneas
						if len(line) > 50 {
							line = line[:47] + "..."
						}
						m.previewContent = append(m.previewContent, line)
					} else {
						m.previewContent = append(m.previewContent, "...")
						break
					}
				}
			} else {
				m.previewContent = append(m.previewContent, "Error reading file: "+err.Error())
			}
		} else {
			// Para archivos binarios
			m.previewContent = append(m.previewContent, "Binary file")
			m.previewContent = append(m.previewContent, "Preview not available")
		}
	}
}

// CountDirItems cuenta los directorios y archivos en un directorio
func (m *Model) CountDirItems(dir string) (int, int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, 0, err
	}

	dirCount := 0
	fileCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			dirCount++
		} else {
			fileCount++
		}
	}
	return dirCount, fileCount, nil
}

// CountSelectedItems cuenta los archivos y directorios seleccionados
func (m *Model) CountSelectedItems() (int, int) {
	selectedFiles := 0
	selectedDirs := 0

	for _, path := range m.selected {
		fi, err := os.Stat(path)
		if err != nil {
			continue
		}

		if fi.IsDir() {
			selectedDirs++

			// Si includeSubdirs está activado, contar también los archivos dentro
			if m.includeSubdirs {
				err := filepath.Walk(path, func(subPath string, info os.FileInfo, err error) error {
					if err != nil {
						return nil // Ignorar errores y continuar
					}

					if subPath != path { // No contar el directorio raíz dos veces
						if info.IsDir() {
							selectedDirs++
						} else {
							selectedFiles++
						}
					}
					return nil
				})

				if err != nil {
					// Ignorar errores de recorrido
				}
			} else {
				// Solo contar archivos directamente en el directorio
				files, err := os.ReadDir(path)
				if err == nil {
					for _, file := range files {
						if !file.IsDir() {
							selectedFiles++
						}
					}
				}
			}
		} else {
			selectedFiles++
		}
	}

	return selectedFiles, selectedDirs
}

// isTextFile determina si un archivo es de texto basado en su extensión
func isTextFile(path string) bool {
	textExtensions := []string{".txt", ".py", ".js", ".java", ".c", ".cpp", ".h", ".html", ".css",
							   ".json", ".xml", ".md", ".sh", ".bat", ".ps1", ".yaml", ".yml",
							   ".ini", ".cfg", ".conf", ".rb", ".go"}

	ext := strings.ToLower(filepath.Ext(path))
	for _, textExt := range textExtensions {
		if ext == textExt {
			return true
		}
	}
	return false
}

// containsPath verifica si una ruta está en la lista de seleccionados
func containsPath(paths []string, target string) bool {
	for _, path := range paths {
		if path == target {
			return true
		}
	}
	return false
}

// isInSelectedDir verifica si una ruta está dentro de un directorio seleccionado
func isInSelectedDir(selectedPaths []string, target string, includeSubdirs bool) bool {
	if !includeSubdirs {
		return false
	}

	for _, path := range selectedPaths {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		if info.IsDir() && strings.HasPrefix(target, path+string(os.PathSeparator)) {
			return true
		}
	}
	return false
}

// purgeManuallyDeselected elimina entradas en manually_deselected que ya no son relevantes
func purgeManuallyDeselected(m *Model, dirPath string) {
	var newDeselected []string
	prefix := dirPath + string(os.PathSeparator)

	for _, path := range m.manually_deselected {
		if !strings.HasPrefix(path, prefix) {
			newDeselected = append(newDeselected, path)
		}
	}

	m.manually_deselected = newDeselected
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Actualizar el tamaño del terminal cuando cambia
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}
	return m, nil
}

// handleKeyPress maneja todas las pulsaciones de teclas
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Si estamos en el panel de vista previa, manejar navegación especial
	if m.activePanel == previewPanel {
		return m.handlePreviewNavigation(msg)
	}

	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "/":
		// TODO: Implementar modo de búsqueda
		return m, nil

	case "k", "up":
		if m.position > 0 {
			m.position--

			// Ajustar la ventana de desplazamiento si es necesario
			if m.position < m.windowStart {
				m.windowStart = m.position
			}
		}

	case "j", "down":
		if m.activePanel == directoriesPanel && m.position < len(m.directories)-1 {
			m.position++

			// Ajustar la ventana de desplazamiento si es necesario
			panelHeight := m.height - 7
			if m.position >= m.windowStart + panelHeight {
				m.windowStart = m.position - panelHeight + 1
			}
		} else if m.activePanel == filesPanel && m.position < len(m.files)-1 {
			m.position++

			// Ajustar la ventana de desplazamiento
			panelHeight := m.height - 7
			if m.position >= m.windowStart + panelHeight {
				m.windowStart = m.position - panelHeight + 1
			}
		}

	case "d":
		switchToDirectoriesPanel(&m)

	case "f":
		switchToFilesPanel(&m)

	case "s":
		handleSelectOnly(&m)

	case "a":
		m.toggleSelectAll()

	case "i":
		m.includeSubdirs = !m.includeSubdirs
		m.statusMessage = "Subdirectories: " + (map[bool]string{true: "Included", false: "Not included"})[m.includeSubdirs]
		m.statusTime = time.Now().Unix()

	case "o":
		// TODO: Implementar exportar y abrir en app externa
		// En Go necesitaremos implementar las funciones correspondientes de utils.py
		return m, nil

	case "c":
		// TODO: Implementar exportar y copiar al portapapeles
		// En Go necesitaremos implementar las funciones del sistema operativo
		return m, nil

	case "enter":
		if m.activePanel == directoriesPanel && len(m.directories) > 0 {
			if m.position >= 0 && m.position < len(m.directories) {
				// Si estamos en '..' volver al directorio padre
				if m.directories[m.position] == ".." {
					return handleGoBackDirectory(m), nil
				}
				// Si no, entrar al directorio
				return handleEnterDirectory(m), nil
			}
		}

	case "tab":
		if m.activePanel == directoriesPanel {
			switchToFilesPanel(&m)
		} else if m.activePanel == filesPanel {
			switchToDirectoriesPanel(&m)
		}

	case "esc", "h":
		// Manejar escape/volver atrás
		return handleEscapeBack(m), nil

	case "l":
		// Similar a Enter
		if m.activePanel == directoriesPanel && len(m.directories) > 0 {
			return handleEnterDirectory(m), nil
		}
	}

	m.updatePreview()
	return m, nil
}

// handlePreviewNavigation maneja la navegación en el panel de vista previa
func (m Model) handlePreviewNavigation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	panelHeight := m.height - 7

	switch msg.String() {
	case "k", "up":
		if m.previewPosition > 0 {
			m.previewPosition--

			// Ajustar la ventana de desplazamiento
			if m.previewPosition < m.previewWindowStart {
				m.previewWindowStart = m.previewPosition
			}
		}

	case "j", "down":
		if m.previewPosition < len(m.previewContent) - 1 {
			m.previewPosition++

			// Ajustar la ventana de desplazamiento
			if m.previewPosition >= m.previewWindowStart + panelHeight {
				m.previewWindowStart = m.previewPosition - panelHeight + 1
			}
		}

	case "d":
		m.activePanel = directoriesPanel

	case "f":
		m.activePanel = filesPanel
	}

	return m, nil
}

// switchToDirectoriesPanel cambia al panel de directorios
func switchToDirectoriesPanel(m *Model) {
	m.activePanel = directoriesPanel

	// Construir elementos de directorio para el filtrado
	var dirItems []string
	if m.directory == m.initialDir {
		dirItems = append([]string{"."}, m.directories...)
	} else {
		dirItems = append([]string{"..", "."}, m.directories...)
	}

	// Intentar posicionar en el directorio previo
	prevDir := m.currentPreviewDirectory
	if prevDir != "" {
		if filepath.Clean(prevDir) == filepath.Clean(m.directory) {
			// Si el directorio de vista previa es el directorio actual
			for i, item := range dirItems {
				if item == "." {
					m.position = i
					break
				}
			}
		} else if filepath.Clean(prevDir) == filepath.Clean(filepath.Dir(m.directory)) {
			// Si el directorio de vista previa es el directorio padre
			for i, item := range dirItems {
				if item == ".." {
					m.position = i
					break
				}
			}
		} else {
			// Para subdirectorios, buscar por nombre base
			dirName := filepath.Base(prevDir)
			for i, item := range dirItems {
				if item == dirName {
					m.position = i
					break
				}
			}
		}
	} else {
		m.position = 0
	}

	m.windowStart = 0
	m.updatePreview()
}

// switchToFilesPanel cambia al panel de archivos
func switchToFilesPanel(m *Model) {
	if m.activePanel == directoriesPanel && len(m.directories) > 0 && m.position < len(m.directories) {
		currentItem := m.directories[m.position]
		if currentItem != ".." && currentItem != "." {
			m.currentPreviewDirectory = filepath.Join(m.directory, currentItem)
		} else if currentItem == ".." {
			m.currentPreviewDirectory = filepath.Dir(m.directory)
		} else {
			m.currentPreviewDirectory = m.directory
		}
	}

	m.activePanel = filesPanel
	m.position = 0
	m.windowStart = 0
	m.updatePreview()
}

// toggleSelectAll selecciona o deselecciona todos los elementos actuales
func (m *Model) toggleSelectAll() {
	var allItems []string

	if m.activePanel == directoriesPanel {
		for _, item := range m.directories {
			if item != "." && item != ".." {
				fullPath := filepath.Join(m.directory, item)
				allItems = append(allItems, fullPath)
			}
		}
	} else if m.activePanel == filesPanel {
		fileDir := m.directory
		if m.currentPreviewDirectory != "" {
			fileDir = m.currentPreviewDirectory
		}

		for _, item := range m.files {
			fullPath := filepath.Join(fileDir, item)
			allItems = append(allItems, fullPath)
		}
	}

	// Verificar si todos los elementos están seleccionados
	allSelected := true
	for _, path := range allItems {
		if !containsPath(m.selected, path) {
			allSelected = false
			break
		}
	}

	// Si todos están seleccionados, deseleccionar; de lo contrario, seleccionar todos
	if allSelected {
		for _, path := range allItems {
			m.removeFromSelection(path)
		}
	} else {
		for _, path := range allItems {
			m.addToSelection(path)
		}
	}
}

// handleSelectOnly maneja la selección/deselección de elementos individuales
func handleSelectOnly(m *Model) {
	if m.activePanel == directoriesPanel && len(m.directories) > 0 && m.position < len(m.directories) {
		item := m.directories[m.position]
		if item != ".." {
			fullPath := filepath.Join(m.directory, item)

			if item == "." {
				fullPath = m.directory
			}

			fi, err := os.Stat(fullPath)
			if err != nil {
				return
			}

			if fi.IsDir() {
				// Verificar si ya está seleccionado
				alreadySelected := false

				// Comprobar si el directorio mismo está seleccionado
				if containsPath(m.selected, fullPath) {
					alreadySelected = true
				} else {
					// Comprobar si hay archivos del directorio seleccionados
					entries, err := os.ReadDir(fullPath)
					if err == nil {
						for _, entry := range entries {
							if !entry.IsDir() {
								entryPath := filepath.Join(fullPath, entry.Name())
								if containsPath(m.selected, entryPath) {
									alreadySelected = true
									break
								}
							}
						}
					}
				}

				if alreadySelected {
					// Deseleccionar directorio y su contenido
					if containsPath(m.selected, fullPath) {
						m.removeFromSelection(fullPath)
					}

					if m.includeSubdirs {
						// Deseleccionar todos los archivos recursivamente
						filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
							if err != nil {
								return nil // Ignorar errores
							}

							if path != fullPath {
								m.removeFromSelection(path)
							}

							return nil
						})
					} else {
						// Deseleccionar solo los archivos directos
						entries, err := os.ReadDir(fullPath)
						if err == nil {
							for _, entry := range entries {
								entryPath := filepath.Join(fullPath, entry.Name())
								m.removeFromSelection(entryPath)
							}
						}
					}

					purgeManuallyDeselected(m, fullPath)
				} else {
					// Seleccionar directorio y su contenido
					if m.includeSubdirs {
						m.addToSelection(fullPath)
						filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
							if err != nil {
								return nil // Ignorar errores
							}

							if path != fullPath {
								m.addToSelection(path)
							}

							return nil
						})
					} else {
						absPath := filepath.Clean(fullPath)
						m.addToSelection(absPath)
					}

					purgeManuallyDeselected(m, fullPath)
				}

				if item == "." {
					m.updatePreview()
				}
			}
		}
	} else if m.activePanel == filesPanel && len(m.files) > 0 && m.position < len(m.files) {
		item := m.files[m.position]
		fileDir := m.directory
		if m.currentPreviewDirectory != "" {
			fileDir = m.currentPreviewDirectory
		}

		fullPath := filepath.Join(fileDir, item)

		// Verificar si está auto-seleccionado por pertenecer a un directorio
		autoSelected := false

		// Comprobar si está en la selección directa
		if containsPath(m.selected, fullPath) {
			autoSelected = true
		} else {
			// Comprobar si está en un directorio seleccionado
			for _, path := range m.selected {
				fi, err := os.Stat(path)
				if err == nil && fi.IsDir() && strings.HasPrefix(fullPath, path+string(os.PathSeparator)) {
					autoSelected = true
					break
				}
			}
		}

		// No considerar auto-seleccionado si está en la lista de deselección manual
		if containsPath(m.manually_deselected, fullPath) {
			autoSelected = false
		}

		if autoSelected {
			// Deseleccionar manualmente este archivo
			m.manually_deselected = append(m.manually_deselected, fullPath)
			m.removeFromSelection(fullPath)
		} else {
			// Eliminar de deselección manual si está
			for i, path := range m.manually_deselected {
				if path == fullPath {
					m.manually_deselected = append(m.manually_deselected[:i], m.manually_deselected[i+1:]...)
					break
				}
			}

			// Añadir a selección
			m.addToSelection(fullPath)
		}
	}
}

// handleGoBackDirectory maneja la navegación hacia atrás en el directorio
func handleGoBackDirectory(m Model) Model {
	if m.activePanel == directoriesPanel {
		currentItem := ""
		if len(m.directories) > 0 && m.position < len(m.directories) {
			currentItem = m.directories[m.position]
		}

		shouldGoBack := currentItem != "."

		if shouldGoBack {
			parentDir := filepath.Dir(m.directory)
			if parentDir != m.directory && strings.HasPrefix(filepath.Clean(parentDir), filepath.Clean(m.initialDir)) {
				m.directory = parentDir
				m.loadDirectory()
				m.position = 0
				m.windowStart = 0
			}
		}
	}
	return m
}

// handleEnterDirectory maneja la navegación dentro de un directorio
func handleEnterDirectory(m Model) Model {
	if m.activePanel == directoriesPanel && len(m.directories) > 0 && m.position < len(m.directories) {
		item := m.directories[m.position]
		if item != ".." && item != "." {
			fullPath := filepath.Join(m.directory, item)
			if fi, err := os.Stat(fullPath); err == nil && fi.IsDir() {
				m.directory = fullPath
				m.loadDirectory()
				m.position = 0
				m.windowStart = 0
			}
		}
	}
	return m
}

// handleEscapeBack maneja la acción de volver atrás (Escape)
func handleEscapeBack(m Model) Model {
	if m.directory != m.initialDir {
		currentSubdir := filepath.Base(m.directory)
		parentDir := filepath.Dir(m.directory)

		if strings.HasPrefix(filepath.Clean(parentDir), filepath.Clean(m.initialDir)) {
			m.directory = parentDir
			m.currentPreviewDirectory = ""
			m.loadDirectory()

			// Intentar encontrar la posición del subdirectorio anterior
			for i, dir := range m.directories {
				if dir == currentSubdir {
					m.position = i
					break
				}
			}

			m.windowStart = 0
		}
	}
	return m
}

// toggleSelection añade o quita un elemento de la selección
func (m *Model) toggleSelection(path string) {
	_, err := os.Stat(path)
	if err != nil {
		return
	}

	isSelected := containsPath(m.selected, path)

	if isSelected {
		m.removeFromSelection(path)
	} else {
		m.addToSelection(path)
	}
}

// addToSelection añade un elemento a la selección, gestionando subdirectorios
func (m *Model) addToSelection(path string) {
	fi, err := os.Stat(path)
	if err != nil {
		return
	}

	isDir := fi.IsDir()

	// Añadir el elemento actual
	if !containsPath(m.selected, path) {
		m.selected = append(m.selected, path)
	}

	// Si es un directorio y se incluyen subdirectorios, añadir todo su contenido
	if isDir && m.includeSubdirs {
		err := filepath.Walk(path, func(subPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if subPath != path && !containsPath(m.selected, subPath) {
				m.selected = append(m.selected, subPath)
			}

			return nil
		})

		if err != nil {
			// Ignorar errores de recorrido
		}
	}
}

// removeFromSelection elimina un elemento de la selección
func (m *Model) removeFromSelection(path string) {
	for i, selected := range m.selected {
		if selected == path {
			// Eliminar de la selección
			m.selected = append(m.selected[:i], m.selected[i+1:]...)
			break
		}
	}

	// Si es un directorio, también eliminar sus subdirectorios y archivos si están seleccionados
	fi, err := os.Stat(path)
	if err == nil && fi.IsDir() {
		prefix := path + string(os.PathSeparator)

		// Eliminar todos los elementos que comienzan con la ruta del directorio
		var newSelected []string
		for _, selected := range m.selected {
			if !strings.HasPrefix(selected, prefix) {
				newSelected = append(newSelected, selected)
			}
		}
		m.selected = newSelected
	}
}
