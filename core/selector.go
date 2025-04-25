package core

import (
	"os"
	"path/filepath"
)

// Estructura para mantener el historial de navegación
type NavigationHistory struct {
	Directory string
	Position  int
}

// Estructura Selector con los campos necesarios
type Selector struct {
	Directory   string            // Directorio actual
	ActivePanel int               // Panel activo: 1 - Directorios, 2 - Archivos, 3 - Vista previa
	Position    int               // Posición actual en el panel de directorios
	FilePosition int              // Posición actual en el panel de archivos
	Selection   map[string]bool   // Items seleccionados
	Filtered    []string          // Items filtrados para mostrar
	Files       []string          // Archivos en el subdirectorio actual
	History     []NavigationHistory // Historial de navegación
	IncludeMode bool              // Modo de inclusión de subdirectorios
}

// Método para actualizar los archivos del directorio seleccionado
func (s *Selector) UpdateFilesForCurrentDirectory() {
	// Si estamos en el panel de directorios, actualizar los archivos del directorio seleccionado
	if s.ActivePanel == 1 && s.Position < len(s.Filtered) {
		item := s.Filtered[s.Position]
		var dir string
		if item == ".." {
			dir = filepath.Dir(s.Directory)
		} else if item == "." {
			dir = s.Directory
		} else {
			dir = filepath.Join(s.Directory, item)
		}

		// Actualizar la lista de archivos para el directorio seleccionado
		files, err := os.ReadDir(dir)
		if err == nil {
			var fileList []string
			for _, file := range files {
				if !file.IsDir() { // Solo archivos
					fileList = append(fileList, file.Name())
				}
			}
			s.Files = fileList // Actualizamos los archivos
		} else {
			s.Files = []string{} // Si hay error, limpiamos la lista de archivos
		}
	}
	// No actualizamos los archivos cuando estamos en el panel de archivos
}
