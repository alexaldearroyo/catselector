package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

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
					isSelected := contains(m.selected, fullPath)
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
		isSelected := contains(m.selected, filePath)
		marker := ""
		if isSelected {
			marker = "*"
		} else {
			marker = " "
		}
		icon := GetFileIcon(filePath)

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

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "j":
			if m.activePanel == directoriesPanel && m.position < len(m.directories)-1 {
				m.position++
			} else if m.activePanel == filesPanel && m.position < len(m.files)-1 {
				m.position++
			} else if m.activePanel == previewPanel && m.previewPosition < len(m.previewContent)-1 {
				m.previewPosition++
			}
		case "k":
			if m.position > 0 {
				m.position--
			} else if m.activePanel == previewPanel && m.previewPosition > 0 {
				m.previewPosition--
			}
		case "h":
			if m.directory != m.initialDir {
				m.directory = filepath.Dir(m.directory)
				m.loadDirectory()
			}
		case "l", "enter":
			if m.activePanel == directoriesPanel && len(m.directories) > 0 {
				subdir := filepath.Join(m.directory, m.directories[m.position])
				m.directory = subdir
				m.loadDirectory()
			}
		case "tab":
			m.activePanel++
			if m.activePanel > previewPanel {
				m.activePanel = directoriesPanel
			}
		case "s":
			// Seleccionar/deseleccionar
			if m.activePanel == directoriesPanel && len(m.directories) > 0 {
				path := filepath.Join(m.directory, m.directories[m.position])
				m.toggleSelection(path)
			} else if m.activePanel == filesPanel && len(m.files) > 0 {
				path := filepath.Join(m.directory, m.files[m.position])
				m.toggleSelection(path)
			}
		case "i":
			// Cambiar inclusión de subdirectorios
			m.includeSubdirs = !m.includeSubdirs
			m.statusMessage = "Subdirectories: " + (map[bool]string{true: "Included", false: "Not included"})[m.includeSubdirs]
			m.statusTime = time.Now().Unix()
		}
		m.updatePreview()
	}
	return m, nil
}

// toggleSelection añade o quita un elemento de la selección
func (m *Model) toggleSelection(path string) {
	for i, selected := range m.selected {
		if selected == path {
			// Eliminar de la selección
			m.selected = append(m.selected[:i], m.selected[i+1:]...)
			return
		}
	}
	// Añadir a la selección
	m.selected = append(m.selected, path)
}
