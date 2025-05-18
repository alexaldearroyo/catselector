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

// Structure to maintain separated search results
type SearchResults struct {
	Directories []string
	Files       []string
}

// Function to search recursively and separate results
func searchRecursively(rootDir string, query string) SearchResults {
	var results SearchResults
	query = strings.ToLower(query)

	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Get the relative name from the root directory
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return nil
		}

		// If the name contains the query, add it to the corresponding results
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
	// If we are in search mode
	if s.SearchMode {
		switch key {
		case "esc":
			// Exit search mode
			s.SearchMode = false
			s.SearchQuery = ""
			s.Filtered = s.OriginalItems
			s.Files = []string{} // Clear the files
			return position
		case "enter":
			// Exit search mode and move the cursor to the first result
			s.SearchMode = false
			s.SearchQuery = ""

			// If there are directories, move to the directory panel
			if len(s.Filtered) > 0 {
				s.ActivePanel = 1 // Directory panel
				s.Position = 0    // First position
				s.DirScroll = 0   // Reset scroll
				return 0
			}
			// If there are no directories but there are files, move to the file panel
			if len(s.Files) > 0 {
				s.ActivePanel = 2 // File panel
				s.FilePosition = 0 // First position
				s.FileScroll = 0   // Reset scroll
				return position
			}
			// If there are no results, return to normal view
			s.Filtered = s.OriginalItems
			s.Files = []string{}
			return position

		case "backspace":
			// Delete the last character of the search
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
			// Add a character to the search
			if len(key) == 1 {
				s.SearchQuery += key
				results := searchRecursively(GetRootDirectory(), s.SearchQuery)
				s.Filtered = results.Directories
				s.Files = results.Files
			}
			return position
		}
	}

	// Normal key handling
	switch key {
	case "/":
		// Enter search mode
		s.SearchMode = true
		s.SearchQuery = ""
		s.OriginalItems = items
		return position
	case "esc", "h":
		// Si estamos en una búsqueda, volver a la vista normal
		if s.SearchMode || (len(s.Filtered) != len(s.OriginalItems) && len(s.OriginalItems) > 0) {
			s.SearchMode = false
			s.Filtered = s.OriginalItems
			s.Files = []string{}
			s.DirScroll = 0
			s.FileScroll = 0
			return 0
		}

		// Comportamiento normal de ESC/h
		rootDir := GetRootDirectory()

		// Si no estamos en el directorio raíz, ir al directorio padre
		if s.Directory != rootDir {
			// Guardar el estado actual en el historial antes de cambiar
			if len(s.History) == 0 || s.History[len(s.History)-1].Directory != s.Directory {
				s.History = append(s.History, NavigationHistory{
					Directory: s.Directory,
					Position:  position,
				})
			}

			// Obtener el directorio padre
			parentDir := filepath.Dir(s.Directory)

			// Verificar que el directorio padre existe y es accesible
			if info, err := os.Stat(parentDir); err == nil && info.IsDir() {
				s.Directory = parentDir
				s.Filtered = PrepareDirItems(parentDir)

				// Buscar la posición del directorio actual en la nueva lista
				// currentDirName := filepath.Base(s.Directory)
				// parentDirName := filepath.Base(parentDir)
				targetName := filepath.Base(s.Directory)

				// Si estamos en el directorio raíz, usar "."
				if parentDir == rootDir {
					targetName = "."
				}

				foundPosition := false
				for i, item := range s.Filtered {
					if item == targetName {
						position = i
						s.Position = i
						foundPosition = true
						break
					}
				}

				// Si no encontramos la posición, usar la primera
				if !foundPosition {
					position = 0
					s.Position = 0
				}

				items = s.Filtered
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
				// Calculate the number of visible lines based on the terminal size
				_, height := getTerminalSize()
				visibleLines := height - 9 // 9 lines for headers and other elements

				// If the current position is outside the visible area, adjust the scroll
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
				// Calculate the number of visible lines based on the terminal size
				_, height := getTerminalSize()
				visibleLines := height - 9 // 9 lines for headers and other elements

				// If the current position is outside the visible area, adjust the scroll
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
				// Adjust the scroll to keep the last element visible
				_, height := getTerminalSize()
				visibleLines := height - 9
				s.DirScroll = max(0, position - visibleLines + 1)
			} else if position < s.DirScroll {
				// Adjust the scroll to keep the current element visible
				s.DirScroll = position
			}
		} else if s.ActivePanel == 2 {
			s.FilePosition--
			if s.FilePosition < 0 {
				s.FilePosition = len(s.Files) - 1
				// Adjust the scroll to keep the last element visible
				_, height := getTerminalSize()
				visibleLines := height - 9
				s.FileScroll = max(0, s.FilePosition - visibleLines + 1)
			} else if s.FilePosition < s.FileScroll {
				// Adjust the scroll to keep the current element visible
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
