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
	Selection   map[string]bool   // Items seleccionados (clave: ruta relativa al directorio actual)
	Filtered    []string          // Items filtrados para mostrar
	Files       []string          // Archivos en el subdirectorio actual
	History     []NavigationHistory // Historial de navegación
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

// Obtiene la clave de selección para un elemento, combinando el directorio actual con el nombre del elemento
func (s *Selector) GetSelectionKey(item string) string {
	if item == "." || item == ".." {
		return item
	}
	return filepath.Join(s.Directory, item)
}

// Obtiene la clave de selección para un archivo, teniendo en cuenta el directorio activo
func (s *Selector) GetFileSelectionKey(file string) string {
	// Si estamos en el panel de archivos, el archivo está en el directorio seleccionado
	if s.ActivePanel == 2 && s.Position < len(s.Filtered) {
		// Determinar el directorio actual para los archivos
		item := s.Filtered[s.Position]
		var currentDir string

		if item == ".." {
			currentDir = filepath.Dir(s.Directory)
		} else if item == "." {
			currentDir = s.Directory
		} else {
			currentDir = filepath.Join(s.Directory, item)
		}

		return filepath.Join(currentDir, file)
	}

	// Por defecto, asumimos que el archivo está en el directorio actual
	return filepath.Join(s.Directory, file)
}

// Verifica si un elemento está seleccionado
func (s *Selector) IsSelected(item string) bool {
	key := s.GetSelectionKey(item)
	// Si el elemento es un directorio y está seleccionado, todos sus archivos también están seleccionados
	if s.Selection[key] {
		return true
	}

	// Si el elemento es un archivo, verificar si su directorio padre está seleccionado
	parentDir := filepath.Dir(key)
	return s.Selection[parentDir]
}

// Verifica si un archivo está seleccionado
func (s *Selector) IsFileSelected(file string) bool {
	key := s.GetFileSelectionKey(file)
	// Verificar si el archivo está seleccionado directamente
	if s.Selection[key] {
		return true
	}

	// Verificar si el directorio padre está seleccionado
	parentDir := filepath.Dir(key)
	return s.Selection[parentDir]
}
