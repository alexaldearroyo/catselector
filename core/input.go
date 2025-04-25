package core

import (
	"os"
	"path/filepath"
)

func CaptureInput(key string) string {
	// Aquí manejas las teclas que recibes
	switch key {
	case "j", "Down":
		return "down"
	case "k", "Up":
		return "up"
	default:
		return ""
	}
}

func HandleKeyPress(key string, position, itemCount int, selected map[string]bool, items []string, s *Selector) int {
	switch key {
	case "down", "j":
		position++
		if position >= itemCount {
			position = 0
		}
	case "up", "k":
		position--
		if position < 0 {
			position = itemCount - 1
		}
	case "enter", "l":
		if position >= 0 && position < len(items) {
			item := items[position]
			var newDir string
			if item == ".." {
				newDir = filepath.Dir(s.Directory)
			} else if item == "." {
				newDir = s.Directory
			} else {
				newDir = filepath.Join(s.Directory, item)
			}

			// Verificar si el directorio existe y es accesible
			if info, err := os.Stat(newDir); err == nil && info.IsDir() {
				// Guardar el estado actual en el historial antes de cambiar
				s.History = append(s.History, NavigationHistory{
					Directory: s.Directory,
					Position:  position,
				})

				s.Directory = newDir
				s.Position = 0
				s.Filtered = PrepareDirItems(newDir)
				position = 0
				items = s.Filtered // Actualizar los items con los nuevos
			}
		}
	case "esc", "h":
		// Navegar hacia atrás si no estamos en el directorio raíz de la aplicación
		rootDir := GetRootDirectory()
		if s.Directory != rootDir && len(s.History) > 0 {
			// Obtener el último estado del historial
			lastState := s.History[len(s.History)-1]

			// Verificar si el directorio del historial existe y es accesible
			if info, err := os.Stat(lastState.Directory); err == nil && info.IsDir() {
				s.Directory = lastState.Directory
				s.Filtered = PrepareDirItems(lastState.Directory)
				position = lastState.Position
				s.Position = lastState.Position
				items = s.Filtered

				// Eliminar el último estado del historial
				s.History = s.History[:len(s.History)-1]
			}
		}
	}

	// Actualizar la posición en el selector
	s.Position = position

	// Actualizar los archivos cuando se navega
	s.UpdateFilesForCurrentDirectory()

	// Actualizar la selección
	if position >= 0 && position < len(items) {
		selectedItem := items[position]
		selected[selectedItem] = !selected[selectedItem]
	}

	return position
}
