package main

import (
	"catselector/core"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Check if --help or --version was requested
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
				switch arg {
				case "--help", "-h":
						printHelp()
						os.Exit(0)
				case "--version", "-v":
						printVersion()
						os.Exit(0)
				}
		}
}


	runApp()
}

func printHelp() {
	fmt.Println(`Cat Selector - Smart Concatenation Selector
A file browser utility that allows you to select multiple files and concatenate them for viewing or editing.

Usage:
  catsel            Start Cat Selector
  catsel --help     Show this help message
  catsel --version  Show version information

Controls:
  j / Down          Move down
  k / Up            Move up
  Enter / l         Enter directory
  Esc / h           Go to previous directory
  s                 Select or deselect
  a                 Select or deselect all
  i                 Toggle include subdirectories in selection
  o                 Concatenate and open selection in external editor
  c                 Concatenate and copy selection to clipboard
  Tab               Switch panel
  f                 Go to files panel
  d                 Go to directories panel
  q                 Quit
`)
}

func printVersion() {
	fmt.Println("Cat Selector version 1.0.3")
}

func runApp() {
	// Enter terminal's alternate screen mode
	fmt.Print("\033[?1049h")
	// Clear the screen
	fmt.Print("\033[H\033[2J")

	// Configure a handler to restore the terminal when exiting
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// Restore the terminal when exiting
		fmt.Print("\033[?1049l")
		os.Exit(0)
	}()

	// Create the initial model
	initialModel := model{
		position: 0,
		items:    core.PrepareDirItems(core.GetRootDirectory()),
		selected: make(map[string]bool),
		selector: core.Selector{
			Directory:   core.GetRootDirectory(),
			ActivePanel: 1,
			Position:    0,
			Selection:   make(map[string]bool),
			Filtered:    core.PrepareDirItems(core.GetRootDirectory()),
			Files:       []string{},
			IncludeMode: false,
			DirScroll:   0,
			FileScroll:  0,
		},
	}

	// Start the program with the model
	p := tea.NewProgram(initialModel)

	// Run the application
	err := p.Start()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Print("\033[?1049l")
		os.Exit(1)
	}

	fmt.Print("\033[?1049l")
}

// The application model for Bubble Tea
type model struct {
	position int
	items    []string
	selected map[string]bool
	selector core.Selector
}

func (m model) Init() tea.Cmd {
	// Initialize the position and items if necessary
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
		// The key handling is done in input.go.
		oldPosition := m.position
		m.position = core.HandleKeyPress(msg.String(), m.position, len(m.items), m.selected, m.items, &m.selector)

		// Synchronize the selector state with the model
		m.selector.Position = m.position
		m.selector.Selection = m.selected
		m.items = m.selector.Filtered // Update the model items with the filtered ones

		// Update the current selector in the core package
		core.SetCurrentSelector(&m.selector)

		// If the position changed in the directory panel, update the files
		if oldPosition != m.position && m.selector.ActivePanel == 1 {
			m.selector.UpdateFilesForCurrentDirectory()
		}
	}
	return m, nil
}

func (m model) View() string {
	// Get the directory elements and the position
	dir := m.selector.Directory
	items := m.selector.Filtered

	// Update the files of the selected directory
	m.selector.UpdateFilesForCurrentDirectory()
	files := m.selector.Files

	// Render the view with the updated files
	return core.DrawLayout(m.position, items, dir, files, m.selector.ActivePanel, m.selector.FilePosition)
}
