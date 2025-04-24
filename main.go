package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	initialDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	model := NewSelector(initialDir)
	p := tea.NewProgram(model)

	if err := p.Start(); err != nil {
		panic(err)
	}
}
