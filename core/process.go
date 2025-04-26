package core

import (
	"os"
	"path/filepath"
	"strings"
)

// Función recursiva para procesar directorios y archivos
func processDirectoryRecursive(selector *Selector, dirPath string, item string, selectState bool) {
	// Actualizar el estado de selección del directorio actual
	selectionKey := selector.GetSelectionKey(item)
	selector.Selection[selectionKey] = selectState

	// Si estamos deseleccionando, limpiar todos los archivos y subdirectorios
	// del directorio actual del mapa de selección
	if !selectState {
		// Borrar entradas de selección relacionadas con este directorio
		prefix := dirPath + string(os.PathSeparator)
		for path := range selector.Selection {
			if strings.HasPrefix(path, prefix) {
				selector.Selection[path] = false
			}
		}
		return
	}

	// Si estamos seleccionando y el modo include está activado,
	// procesar recursivamente los subdirectorios
	if selector.IncludeMode {
		// Leer el contenido del directorio
		entries, err := os.ReadDir(dirPath)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					// Procesar recursivamente el subdirectorio
					subDirPath := filepath.Join(dirPath, entry.Name())
					subItem := filepath.Join(item, entry.Name())
					processDirectoryRecursive(selector, subDirPath, subItem, selectState)
				}
				// No marcamos los archivos individualmente ya que el directorio padre
				// ya está seleccionado y la función IsFileSelected comprobará esto
			}
		}
	}
}
