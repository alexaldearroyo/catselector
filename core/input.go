package core

import (
	"catselector/export"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CaptureInput(key string) string {
	switch key {
	case "j", "Down":
		return "down"
	case "k", "Up":
		return "up"
	default:
		return ""
	}
}

// Estructura para mantener los resultados de búsqueda separados
type SearchResults struct {
	Directories []string
	Files       []string
}

// Función para buscar recursivamente y separar resultados
func searchRecursively(rootDir string, query string) SearchResults {
	var results SearchResults
	query = strings.ToLower(query)

	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Obtener el nombre relativo desde el directorio raíz
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return nil
		}

		// Si el nombre contiene la consulta, añadirlo a los resultados correspondientes
		if strings.Contains(strings.ToLower(relPath), query) {
			if info.IsDir() {
				results.Directories = append(results.Directories, relPath)
			} else {
				results.Files = append(results.Files, relPath)
			}
		}

		return nil
	})

	return results
}

func HandleKeyPress(key string, position, itemCount int, selected map[string]bool, items []string, s *Selector) int {
	// Si estamos en modo búsqueda
	if s.SearchMode {
		switch key {
		case "esc":
			// Salir del modo búsqueda
			s.SearchMode = false
			s.SearchQuery = ""
			s.Filtered = s.OriginalItems
			s.Files = []string{} // Limpiar los archivos
			return position
		case "enter":
			// Salir del modo búsqueda y mover el cursor al primer resultado
			s.SearchMode = false
			s.SearchQuery = ""

			// Si hay directorios, mover al panel de directorios
			if len(s.Filtered) > 0 {
				s.ActivePanel = 1 // Panel de directorios
				s.Position = 0    // Primera posición
				s.DirScroll = 0   // Resetear scroll
				return 0
			}
			// Si no hay directorios pero hay archivos, mover al panel de archivos
			if len(s.Files) > 0 {
				s.ActivePanel = 2 // Panel de archivos
				s.FilePosition = 0 // Primera posición
				s.FileScroll = 0   // Resetear scroll
				return position
			}
			// Si no hay resultados, volver a la vista normal
			s.Filtered = s.OriginalItems
			s.Files = []string{}
			return position

		case "backspace":
			// Eliminar último carácter de la búsqueda
			if len(s.SearchQuery) > 0 {
				s.SearchQuery = s.SearchQuery[:len(s.SearchQuery)-1]
				if s.SearchQuery == "" {
					s.Filtered = s.OriginalItems
					s.Files = []string{}
				} else {
					results := searchRecursively(GetRootDirectory(), s.SearchQuery)
					s.Filtered = results.Directories
					s.Files = results.Files
				}
			}
			return position
		default:
			// Añadir carácter a la búsqueda
			if len(key) == 1 {
				s.SearchQuery += key
				results := searchRecursively(GetRootDirectory(), s.SearchQuery)
				s.Filtered = results.Directories
				s.Files = results.Files
			}
			return position
		}
	}

	// Manejo normal de teclas
	switch key {
	case "/":
		// Entrar en modo búsqueda
		s.SearchMode = true
		s.SearchQuery = ""
		s.OriginalItems = items
		return position
	case "esc", "h":
		// Si estamos en un resultado de búsqueda, volver a la vista normal
		if len(s.Filtered) != len(s.OriginalItems) || len(s.Files) > 0 {
			s.Filtered = s.OriginalItems
			s.Files = []string{}
			s.DirScroll = 0
			s.FileScroll = 0
			return 0
		}
		// Comportamiento normal de ESC/h
		rootDir := GetRootDirectory()
		if s.Directory != rootDir && len(s.History) > 0 {
			lastState := s.History[len(s.History)-1]
			if info, err := os.Stat(lastState.Directory); err == nil && info.IsDir() {
				s.Directory = lastState.Directory
				s.Filtered = PrepareDirItems(lastState.Directory)
				position = lastState.Position
				s.Position = lastState.Position
				items = s.Filtered
				s.History = s.History[:len(s.History)-1]
				s.DirScroll = 0
			}
		}
	case "q":
		// Restore the terminal and exit
		fmt.Print("\033[?1049l")
		os.Exit(0)
	case "down", "j":
		if s.ActivePanel == 1 {
			position++
			if position >= itemCount {
				position = 0
				s.DirScroll = 0
			} else {
				// Calcular el número de líneas visibles basado en el tamaño de la terminal
				_, height := getTerminalSize()
				visibleLines := height - 9 // 9 líneas para headers y otros elementos

				// Si la posición actual está fuera del área visible, ajustar el scroll
				if position >= s.DirScroll + visibleLines {
					s.DirScroll = position - visibleLines + 1
				}
			}
		} else if s.ActivePanel == 2 {
			s.FilePosition++
			if s.FilePosition >= len(s.Files) {
				s.FilePosition = 0
				s.FileScroll = 0
			} else {
				// Calcular el número de líneas visibles basado en el tamaño de la terminal
				_, height := getTerminalSize()
				visibleLines := height - 9 // 9 líneas para headers y otros elementos

				// Si la posición actual está fuera del área visible, ajustar el scroll
				if s.FilePosition >= s.FileScroll + visibleLines {
					s.FileScroll = s.FilePosition - visibleLines + 1
				}
			}
		}
	case "up", "k":
		if s.ActivePanel == 1 {
			position--
			if position < 0 {
				position = itemCount - 1
				// Ajustar el scroll para que el último elemento esté visible
				_, height := getTerminalSize()
				visibleLines := height - 9
				s.DirScroll = max(0, position - visibleLines + 1)
			} else if position < s.DirScroll {
				// Ajustar el scroll para mantener visible el elemento actual
				s.DirScroll = position
			}
		} else if s.ActivePanel == 2 {
			s.FilePosition--
			if s.FilePosition < 0 {
				s.FilePosition = len(s.Files) - 1
				// Ajustar el scroll para que el último elemento esté visible
				_, height := getTerminalSize()
				visibleLines := height - 9
				s.FileScroll = max(0, s.FilePosition - visibleLines + 1)
			} else if s.FilePosition < s.FileScroll {
				// Ajustar el scroll para mantener visible el elemento actual
				s.FileScroll = s.FilePosition
			}
		}
	case "i":
		// Toggle the include mode
		s.IncludeMode = !s.IncludeMode
	case "o":
		// Export and open in external application
		selectedPaths := getSelectedPaths(s.Selection)
		outputFile := export.GenerateTextFile(
			selectedPaths,
			[]string{}, // Empty excluded paths
			s.IncludeMode,
			GetRootDirectory(),
			s.Directory,
		)
		if outputFile != "" {
			// Open the file without exiting the alternative mode
			err := OpenTextFile(outputFile)

			// Show success or error message
			if err == nil {
				s.StatusMessage = "Opened file: " + filepath.Base(outputFile)
			} else {
				s.StatusMessage = "Error opening file"
			}
			s.StatusTime = time.Now().Unix()
		}
	case "c":
		// Export and copy to clipboard and delete file
		selectedPaths := getSelectedPaths(s.Selection)
		outputFile := export.GenerateTextFile(
			selectedPaths,
			[]string{}, // Empty excluded paths
			s.IncludeMode,
			GetRootDirectory(),
			s.Directory,
		)

		if outputFile != "" {
			// Read the content of the file
			content, err := os.ReadFile(outputFile)
			if err == nil {
				// Copy to clipboard according to the operating system
				success := CopyToClipboard(string(content))

				// Delete the temporary file
				os.Remove(outputFile)

				// Count selected files for the status message
				selectedFiles := countSelectedFiles(s)

				// Prepare message and save status
				msg := ""
				if success {
					msg = fmt.Sprintf("%d files copied to clipboard", selectedFiles)
				} else {
					msg = "Error copying to clipboard"
				}
				s.StatusMessage = msg
				s.StatusTime = time.Now().Unix()
			}
		}
	case "tab":
		// Save the previous panel
		previousPanel := s.ActivePanel

		// Change only between the directory and file panels (1 and 2)
		if s.ActivePanel == 1 {
			s.ActivePanel = 2
		} else {
			s.ActivePanel = 1
		}

		// Only update the files when changing from the directory panel to the file panel
		if previousPanel == 1 && s.ActivePanel == 2 {
			// If we come from the directory panel, update the files of the selected directory
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

				// Check if the directory exists and is accessible
				if info, err := os.Stat(selectedDir); err == nil && info.IsDir() {
					// Update the file list for the selected directory
					files, err := os.ReadDir(selectedDir)
					if err == nil {
						var fileList []string
						for _, file := range files {
							if !file.IsDir() { // Solo archivos
								fileList = append(fileList, file.Name())
							}
						}
						s.Files = fileList // Update the files
						s.FilePosition = 0 // Reset the position in the file panel
						s.FileScroll = 0   // Reset the scroll position
					} else {
						s.Files = []string{} // If there is an error, clear the file list
					}
				}
			}
		} else if s.ActivePanel == 2 && len(s.Files) > 0 {
			// If we are already in the file panel, only reset the position
			s.FilePosition = 0
			s.FileScroll = 0
		}
	case "f":
		// Change to the file panel
		s.ActivePanel = 2
		if len(s.Files) > 0 {
			s.FilePosition = 0
			s.FileScroll = 0
		}
	case "d":
		// Change to the directory panel
		s.ActivePanel = 1
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

			// Check if the directory exists and is accessible
			if info, err := os.Stat(newDir); err == nil && info.IsDir() {
				// Save the current state in the history before changing
				s.History = append(s.History, NavigationHistory{
					Directory: s.Directory,
					Position:  position,
				})

				s.Directory = newDir
				s.Filtered = PrepareDirItems(newDir)

				// Search for the position of "." in the new list
				for i, item := range s.Filtered {
					if item == "." {
						position = i
						s.Position = i
						break
					}
				}

				items = s.Filtered // Update the items with the new ones
				s.DirScroll = 0    // Reset the scroll position
			}
		}
	case "s":
		if s.ActivePanel == 1 {
			// Toggle the selection of the current directory
			if position >= 0 && position < len(items) {
				item := items[position]
				if item != ".." {
					// Get the current selector
					selector := GetCurrentSelector()

					// Determine the current selection state
					isSelected := selector.IsSelected(item)

					// Process the current directory
					dirPath := filepath.Join(s.Directory, item)
					processDirectoryRecursive(selector, dirPath, item, !isSelected)

					// Update the file list if necessary
					if !isSelected {
						// If we are selecting, update the file list
						UpdateFileList(selector, s.Directory, item)
					} else {
						// If we are deselecting, clear the file list
						s.Files = []string{}
						s.FilePosition = 0
						s.FileScroll = 0
					}
				}
			}
		} else if s.ActivePanel == 2 && s.FilePosition >= 0 && s.FilePosition < len(s.Files) {
			// Get the name of the selected file
			selectedFile := s.Files[s.FilePosition]
			// Change the selection state
			fileKey := s.GetFileSelectionKey(selectedFile)
			s.Selection[fileKey] = !s.Selection[fileKey]
		}
	case "a":
		if s.ActivePanel == 1 {
			// Check if all directories are selected (excluding '..' and '.')
			allSelected := true
			for _, item := range items {
				if item != ".." && item != "." && !s.IsSelected(item) {
					allSelected = false
					break
				}
			}

			// If all are selected, deselect all (excluding '..' and '.')
			// If not all are selected, select all (excluding '..' and '.')
			for _, item := range items {
				if item != ".." && item != "." {
					// Process the directory and its subdirectories if the include mode is active
					dirPath := filepath.Join(s.Directory, item)
					processDirectoryRecursive(s, dirPath, item, !allSelected)
				}
			}
		} else if s.ActivePanel == 2 {
			// Check if all files are selected
			allSelected := true
			for _, file := range s.Files {
				fileKey := s.GetFileSelectionKey(file)
				if !s.Selection[fileKey] {
					allSelected = false
					break
				}
			}

			// If all are selected, deselect all
			// If not all are selected, select all
			for _, file := range s.Files {
				fileKey := s.GetFileSelectionKey(file)
				s.Selection[fileKey] = !allSelected
			}
		}
	}

	// Update the position in the selector
	s.Position = position

	// Update the files when navigating
	s.UpdateFilesForCurrentDirectory()

	return position
}

// UpdateFileList updates the file list for a directory
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

// Function to get a list of selected paths
func getSelectedPaths(selection map[string]bool) []string {
	var paths []string
	for path, selected := range selection {
		if selected {
			paths = append(paths, path)
		}
	}
	return paths
}

// Count selected files
func countSelectedFiles(s *Selector) int {
	count := 0
	processedDirs := make(map[string]bool)

	for path, selected := range s.Selection {
		if !selected {
			continue
		}

		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		if !info.IsDir() {
			count++
		} else if s.IncludeMode {
			// If the directory has already been processed, skip it
			if processedDirs[path] {
				continue
			}
			processedDirs[path] = true

			// Count files in the directory and subdirectories
			filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
				if err == nil && !info.IsDir() {
					count++
				}
				return nil
			})
		} else {
			// Count only files in the top level of the directory
			files, err := os.ReadDir(path)
			if err == nil {
				for _, file := range files {
					fileInfo, err := file.Info()
					if err == nil && !fileInfo.IsDir() {
						count++
					}
				}
			}
		}
	}
	return count
}

// Nueva función para filtrar items
func filterItems(items []string, query string) []string {
	if query == "" {
		return items
	}

	query = strings.ToLower(query)
	var filtered []string

	for _, item := range items {
		if strings.Contains(strings.ToLower(item), query) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}
