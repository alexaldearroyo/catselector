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

// Crear el modelo inicial
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
	},
}

	// Iniciar el programa con el modelo
	p := tea.NewProgram(initialModel)

	// Ejecutar la aplicación
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

// El modelo de la aplicación para Bubble Tea
type model struct {
	position int
	items    []string
	selected map[string]bool
	selector core.Selector // Agregar el selector aquí
}

func (m model) Init() tea.Cmd {
	// Inicializar la posición y los elementos si es necesario
	return nil
}

// main.go

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
		// El manejo de las teclas ya se hace en input.go.
		m.position = core.HandleKeyPress(msg.String(), m.position, len(m.items), m.selected, m.items, &m.selector)

		// Sincronizar el estado del selector con el modelo
		m.selector.Position = m.position
		m.selector.Selection = m.selected
		m.items = m.selector.Filtered // Actualizar los items del modelo con los filtrados
	}
	return m, nil
}

func (m model) View() string {
	// Obtener los elementos de los directorios y la posición
	dir := m.selector.Directory
	items := m.selector.Filtered

	// Actualizar los archivos del directorio seleccionado
	m.selector.UpdateFilesForCurrentDirectory()
	files := m.selector.Files

	// Renderizar la vista con los archivos actualizados
	return core.DrawLayout(m.position, items, dir, files)
}
