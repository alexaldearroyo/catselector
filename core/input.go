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
		if s.ActivePanel == 1 {
			position++
			if position >= itemCount {
				position = 0
			}
		} else if s.ActivePanel == 2 {
			s.FilePosition++
			if s.FilePosition >= len(s.Files) {
				s.FilePosition = 0
			}
		}
	case "up", "k":
		if s.ActivePanel == 1 {
			position--
			if position < 0 {
				position = itemCount - 1
			}
		} else if s.ActivePanel == 2 {
			s.FilePosition--
			if s.FilePosition < 0 {
				s.FilePosition = len(s.Files) - 1
			}
		}
	case "i":
		// Toggle del modo include
		s.IncludeMode = !s.IncludeMode
	case "tab":
		// Guardar el panel anterior
		previousPanel := s.ActivePanel

		// Cambiar solo entre los paneles de directorios y archivos (1 y 2)
		if s.ActivePanel == 1 {
			s.ActivePanel = 2
		} else {
			s.ActivePanel = 1
		}

		// Solo actualizar los archivos cuando cambiamos del panel de directorios al panel de archivos
		if previousPanel == 1 && s.ActivePanel == 2 {
			// Si venimos del panel de directorios, actualizar los archivos del directorio seleccionado
			if position >= 0 && position < len(items) {
				item := items[position]
				var selectedDir string
				if item == ".." {
					selectedDir = filepath.Dir(s.Directory)
				} else if item == "." {
					selectedDir = s.Directory
				} else {
					selectedDir = filepath.Join(s.Directory, item)
				}

				// Verificar si el directorio existe y es accesible
				if info, err := os.Stat(selectedDir); err == nil && info.IsDir() {
					// Actualizar la lista de archivos para el directorio seleccionado
					files, err := os.ReadDir(selectedDir)
					if err == nil {
						var fileList []string
						for _, file := range files {
							if !file.IsDir() { // Solo archivos
								fileList = append(fileList, file.Name())
							}
						}
						s.Files = fileList // Actualizamos los archivos
						s.FilePosition = 0 // Resetear la posición en el panel de archivos
					} else {
						s.Files = []string{} // Si hay error, limpiamos la lista de archivos
					}
				}
			}
		} else if s.ActivePanel == 2 && len(s.Files) > 0 {
			// Si ya estamos en el panel de archivos, solo resetear la posición
			s.FilePosition = 0
		}
	case "enter", "l":
		if s.ActivePanel == 1 && position >= 0 && position < len(items) {
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
				s.Filtered = PrepareDirItems(newDir)

				// Buscar la posición de "." en la nueva lista
				for i, item := range s.Filtered {
					if item == "." {
						position = i
						s.Position = i
						break
					}
				}

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
	case "s":
		if s.ActivePanel == 1 {
			// Toggle de selección del directorio actual
			if position >= 0 && position < len(items) {
				item := items[position]
				if item != ".." {
					// Obtener el selector actual
					selector := GetCurrentSelector()

					// Determinar el estado actual de selección
					isSelected := selector.IsSelected(item)

					// Procesar el directorio actual
					dirPath := filepath.Join(s.Directory, item)
					processDirectoryRecursive(selector, dirPath, item, !isSelected)

					// Actualizar la lista de archivos si es necesario
					if !isSelected {
						// Si estamos seleccionando, actualizar la lista de archivos
						UpdateFileList(selector, s.Directory, item)
					} else {
						// Si estamos deseleccionando, limpiar la lista de archivos
						s.Files = []string{}
						s.FilePosition = 0
					}
				}
			}
		} else if s.ActivePanel == 2 && s.FilePosition >= 0 && s.FilePosition < len(s.Files) {
			// Obtener el nombre del archivo seleccionado
			selectedFile := s.Files[s.FilePosition]
			// Cambiar el estado de selección
			fileKey := s.GetFileSelectionKey(selectedFile)
			s.Selection[fileKey] = !s.Selection[fileKey]
		}
	case "a":
		if s.ActivePanel == 1 {
			// Verificar si todos los directorios están seleccionados (excluyendo '..' y '.')
			allSelected := true
			for _, item := range items {
				if item != ".." && item != "." && !s.IsSelected(item) {
					allSelected = false
					break
				}
			}

			// Si todos están seleccionados, deseleccionar todos (excluyendo '..' y '.')
			// Si no todos están seleccionados, seleccionar todos (excluyendo '..' y '.')
			for _, item := range items {
				if item != ".." && item != "." {
					selectionKey := s.GetSelectionKey(item)
					s.Selection[selectionKey] = !allSelected
				}
			}
		} else if s.ActivePanel == 2 {
			// Verificar si todos los archivos están seleccionados
			allSelected := true
			for _, file := range s.Files {
				fileKey := s.GetFileSelectionKey(file)
				if !s.Selection[fileKey] {
					allSelected = false
					break
				}
			}

			// Si todos están seleccionados, deseleccionar todos
			// Si no todos están seleccionados, seleccionar todos
			for _, file := range s.Files {
				fileKey := s.GetFileSelectionKey(file)
				s.Selection[fileKey] = !allSelected
			}
		}
	}

	// Actualizar la posición en el selector
	s.Position = position

	// Actualizar los archivos cuando se navega
	s.UpdateFilesForCurrentDirectory()

	return position
}

// UpdateFileList actualiza la lista de archivos para un directorio
func UpdateFileList(selector *Selector, currentDir string, item string) {
	dirPath := filepath.Join(currentDir, item)
	files, err := os.ReadDir(dirPath)
	if err == nil {
		var fileList []string
		for _, file := range files {
			if !file.IsDir() { // Solo archivos
				fileList = append(fileList, file.Name())
			}
		}
		selector.Files = fileList
	}
}
