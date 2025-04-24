package main

import (
	"os"
	"path/filepath"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
)

type panel int

const (
	directoriesPanel panel = iota + 1
	filesPanel
	previewPanel
)

type Model struct {
	directory          string
	initialDir         string
	directories        []string
	files              []string
	position           int
	activePanel        panel
	previewContent     []string
	selected           []string // ← añade esto
	includeSubdirs     bool     // ← y esto también
}

func NewSelector(initialDir string) Model {
	model := Model{
		directory:       initialDir,
		initialDir:      initialDir,
		activePanel:     directoriesPanel,
		selected: []string{},
		includeSubdirs:  false,                 // ← por defecto
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
	m.updatePreview()
}

func (m *Model) updatePreview() {
	m.previewContent = []string{}

	if m.activePanel == directoriesPanel && len(m.directories) > 0 && m.position < len(m.directories) {
		m.previewContent = append(m.previewContent, "[Subdirectory Preview]")
		name := m.directories[m.position]
		subPath := filepath.Join(m.directory, name)
		entries, err := os.ReadDir(subPath)
		if err == nil {
			for _, e := range entries {
				if e.IsDir() {
					m.previewContent = append(m.previewContent, "📁 "+e.Name())
				}
			}
		}
	} else if m.activePanel == filesPanel && len(m.files) > 0 && m.position < len(m.files) {
		m.previewContent = append(m.previewContent, "[File Preview]")
		name := m.files[m.position]
		m.previewContent = append(m.previewContent, name)
	}
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
			}
		case "k":
			if m.position > 0 {
				m.position--
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
		}
		m.updatePreview()
	}
	return m, nil
}
