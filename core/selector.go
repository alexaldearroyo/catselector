package core

import (
	"os"
	"path/filepath"
)

// Structure to maintain the navigation history
type NavigationHistory struct {
	Directory string
	Position  int
}

// Structure Selector with the necessary fields
type Selector struct {
	Directory    string            // Current directory
	ActivePanel  int               // Active panel: 1 - Directories, 2 - Files, 3 - Preview
	Position     int               // Current position in the directory panel
	FilePosition int               // Current position in the files panel
	Selection    map[string]bool   // Selected items (key: relative path to the current directory)
	Filtered     []string          // Items filtered to display
	Files        []string          // Files in the current directory
	History      []NavigationHistory // Navigation history
	IncludeMode  bool              // Include mode for subdirectories
	StatusMessage string           // Status message to display to the user
	StatusTime   int64             // Time when the status message was set
	DirScroll    int               // Scroll position for directories panel
	FileScroll   int               // Scroll position for files panel
	// Nuevos campos para la búsqueda
	SearchMode   bool              // Indica si estamos en modo búsqueda
	SearchQuery  string            // La consulta de búsqueda actual
	OriginalItems []string         // Guarda los items originales antes de la búsqueda
	IsSearching  bool              // Nuevo campo para indicar si estamos en modo búsqueda global
}

// Method to update the files of the selected directory
func (s *Selector) UpdateFilesForCurrentDirectory() {
	// If we are in the directory panel, update the files of the selected directory
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

		// Update the list of files for the selected directory
		files, err := os.ReadDir(dir)
		if err == nil {
			var fileList []string
			for _, file := range files {
				if !file.IsDir() { // Only files
					fileList = append(fileList, file.Name())
				}
			}
			s.Files = fileList // Update the files
		} else {
			s.Files = []string{} // If there is an error, clear the list of files
		}
	}
	// We don't update the files when we are in the files panel
}

// Get the selection key for an item, combining the current directory with the name of the item
func (s *Selector) GetSelectionKey(item string) string {
	if item == "." || item == ".." {
		return item
	}
	return filepath.Join(s.Directory, item)
}

// Get the selection key for a file, taking into account the active directory
func (s *Selector) GetFileSelectionKey(file string) string {
	// If we are in the files panel, the file is in the selected directory
	if s.ActivePanel == 2 && s.Position < len(s.Filtered) {
		// Determine the current directory for the files
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

	// By default, we assume that the file is in the current directory
	return filepath.Join(s.Directory, file)
}

// Check if an item is selected
func (s *Selector) IsSelected(item string) bool {
	key := s.GetSelectionKey(item)
	// If the item is a directory and is selected, all its files are also selected
	if s.Selection[key] {
		return true
	}

	// If the item is a file, check if its parent directory is selected
	parentDir := filepath.Dir(key)
	return s.Selection[parentDir]
}

// Check if a file is selected
func (s *Selector) IsFileSelected(file string) bool {
	key := s.GetFileSelectionKey(file)
	// Check if the file is selected directly
	if s.Selection[key] {
		return true
	}

	// Check if the parent directory is selected
	parentDir := filepath.Dir(key)
	return s.Selection[parentDir]
}

// Recursive function to process directories and files
func processDirectory(selector *Selector, dirPath string, item string, selectState bool) {
	// Update the selection state of the current directory
	selectionKey := selector.GetSelectionKey(item)
	selector.Selection[selectionKey] = selectState

	// If the include mode is active, process recursively the subdirectories
	if selector.IncludeMode {
		// Read the content of the directory
		entries, err := os.ReadDir(dirPath)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					// Process recursively the subdirectory
					subDirPath := filepath.Join(dirPath, entry.Name())
					subItem := filepath.Join(item, entry.Name())
					processDirectory(selector, subDirPath, subItem, selectState)
				} else {
					// Select files in the current directory
					fileKey := filepath.Join(dirPath, entry.Name())
					selector.Selection[fileKey] = selectState
				}
			}
		}
	} else {
		// Original behavior: only select files in the current directory
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
