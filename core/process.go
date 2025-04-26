package core

import (
	"os"
	"path/filepath"
)

// Función recursiva para procesar directorios y archivos
func processDirectoryRecursive(selector *Selector, dirPath string, item string, selectState bool) {
	// Actualizar el estado de selección del directorio actual
	selectionKey := selector.GetSelectionKey(item)
	selector.Selection[selectionKey] = selectState

	// Si el modo include está activado, procesar recursivamente los subdirectorios
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
				} else {
					// Seleccionar archivos en el directorio actual
					fileKey := filepath.Join(dirPath, entry.Name())
					selector.Selection[fileKey] = selectState
				}
			}
		}
	} else {
		// Comportamiento original: solo seleccionar archivos en el directorio actual
		entries, err := os.ReadDir(dirPath)
		if err == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					fileKey := filepath.Join(dirPath, entry.Name())
					selector.Selection[fileKey] = selectState
				}
			}
		}
	}
}
