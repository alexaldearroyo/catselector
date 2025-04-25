package main

import (
	"catexplore/core"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Entrar en el modo alternativo de la terminal
	fmt.Print("\033[?1049h")
	// Limpiar la pantalla
	fmt.Print("\033[H\033[2J")

	// Configurar un manejador para restaurar la terminal al salir
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// Restaurar la terminal al salir
		fmt.Print("\033[?1049l")
		os.Exit(0)
	}()

	// Start the application model
	p := tea.NewProgram(model{}) // Pass the model here

	// Display the interface for the first time
	err := p.Start()
	if err != nil {
		fmt.Println("Error: ", err)
		// Restaurar la terminal al salir con error
		fmt.Print("\033[?1049l")
		os.Exit(1)
	}

	// Restaurar la terminal al salir normalmente
	fmt.Print("\033[?1049l")
}

// The application model for Bubble Tea
type model struct{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Handle key messages
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c": // Exit the program
			return m, tea.Quit
		case "enter": // Handle Enter key
			fmt.Println("Enter key pressed")
		}
	}
	return m, nil
}

func (m model) View() string {
	// Call the layout drawing function and return the string
	return core.DrawLayout()
}
