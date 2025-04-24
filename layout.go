// layout.go
package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	dirStyle      = lipgloss.NewStyle().MarginRight(1)
	fileStyle     = lipgloss.NewStyle().MarginRight(1)
	previewStyle  = lipgloss.NewStyle().MarginLeft(1)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	focusStyle    = lipgloss.NewStyle().Background(lipgloss.Color("236"))
)

func renderDirectoryPanel(m Model, height, width int) string {
	lines := []string{}
	items := m.directories
	if m.directory != m.initialDir {
		items = append([]string{"..", "."}, items...)
	} else {
		items = append([]string{"."}, items...)
	}

	for i := 0; i < height && i < len(items); i++ {
		item := items[i]
		icon := "📁"
		fullPath := filepath.Join(m.directory, item)
		isSelected := false
		for _, sel := range m.selected {
			if sel == fullPath {
				isSelected = true
				break
			}
		}
		marker := "  "
		if isSelected {
			if m.includeSubdirs {
				marker = "* "
			} else {
				marker = "• "
			}
		}

		content := marker + item
		if len(content) > width-2 {
			content = content[:width-5] + "..."
		}

		line := icon + " " + content
		if m.activePanel == directoriesPanel && m.position == i {
			line = focusStyle.Render(line)
		}
		lines = append(lines, line)
	}

	// Scrollbar mock (provisional, sin scroll real aún)
	if len(items) > height {
		lines[height-1] = lines[height-1][:width-1] + "↓"
	}

	return strings.Join(lines, "\n")
}

func (m Model) View() string {
	header := headerStyle.Render("ExtraCat - " + m.directory)

	// DIRECTORIOS
	dirs := []string{"Directories:"}
	for i, d := range m.directories {
		line := d
		if i == m.position && m.activePanel == directoriesPanel {
			line = focusStyle.Render("* " + line)
		} else {
			line = "  " + line
		}
		dirs = append(dirs, line)
	}
	dirView := dirStyle.Render(strings.Join(dirs, "\n"))

	// ARCHIVOS
	files := []string{"Files:"}
	for i, f := range m.files {
		line := f
		if i == m.position && m.activePanel == filesPanel {
			line = focusStyle.Render("* " + line)
		} else {
			line = "  " + line
		}
		files = append(files, line)
	}
	fileView := fileStyle.Render(strings.Join(files, "\n"))

	// PREVIEW
	preview := []string{"Preview:"}
	for _, line := range m.previewContent {
		preview = append(preview, line)
	}
	previewView := previewStyle.Render(strings.Join(preview, "\n"))

	// COMPOSICIÓN FINAL
	layout := lipgloss.JoinHorizontal(lipgloss.Top, dirView, fileView, previewView)

	return fmt.Sprintf("%s\n%s", header, layout)
}
