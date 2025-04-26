// core/styles.go
package core

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	// Directories
	Green = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	// Panel Headers
	Cyan = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	// Files
	White = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	// Selected item
	// Marked items
	Marked = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Background(lipgloss.Color("0"))
	// Headers
	Header = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Background(lipgloss.Color("0"))
	// Active panel header
	ActiveHeader = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("6"))
	// Directory icon
	DirectoryIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Background(lipgloss.Color("0"))
	// File icon
	FileIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Background(lipgloss.Color("0"))
	// Key hints
	KeyHints = lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Background(lipgloss.Color("0"))
	// Key hint text
	KeyHintText = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Background(lipgloss.Color("0"))
	// Magenta text
	Magenta = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	// Yellow text for selected files
	Yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	// Blue text for counters
	Blue = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
)

// Styles for the layout (combined into one block)
var (
	// Directory text and the directory display
	DirectoryText = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	DirectoryDir  = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	// Title on the right (Cat Explorer)
	HeaderTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
)


var (
	Focus = lipgloss.NewStyle().
	Foreground(lipgloss.Color("0")).
	Background(lipgloss.Color("7")) // blanco
	Selected = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
)
