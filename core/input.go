package core

import (
	"catexplorer/export"
	"fmt"
	"os"
	"path/filepath"
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

func HandleKeyPress(key string, position, itemCount int, selected map[string]bool, items []string, s *Selector) int {
	switch key {
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
			} else if position >= s.DirScroll + 10 { // 10 es el número de líneas visibles
				// Cuando llegamos al último elemento visible, movemos el scroll una línea
				s.DirScroll++
			}
		} else if s.ActivePanel == 2 {
			s.FilePosition++
			if s.FilePosition >= len(s.Files) {
				s.FilePosition = 0
				s.FileScroll = 0
			} else if s.FilePosition >= s.FileScroll + 10 {
				// Cuando llegamos al último elemento visible, movemos el scroll una línea
				s.FileScroll++
			}
		}
	case "up", "k":
		if s.ActivePanel == 1 {
			position--
			if position < 0 {
				position = itemCount - 1
				// Ajustar el scroll para que el último elemento esté visible
				s.DirScroll = max(0, position - 9)
			} else if position < s.DirScroll {
				// Cuando subimos, movemos el scroll una línea hacia arriba
				s.DirScroll--
			}
		} else if s.ActivePanel == 2 {
			s.FilePosition--
			if s.FilePosition < 0 {
				s.FilePosition = len(s.Files) - 1
				// Ajustar el scroll para que el último elemento esté visible
				s.FileScroll = max(0, s.FilePosition - 9)
			} else if s.FilePosition < s.FileScroll {
				// Cuando subimos, movemos el scroll una línea hacia arriba
				s.FileScroll--
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
	case "esc", "h":
		// Navigate back if we are not in the root directory of the application
		rootDir := GetRootDirectory()
		if s.Directory != rootDir && len(s.History) > 0 {
			// Get the last state of the history
			lastState := s.History[len(s.History)-1]

			// Check if the directory of the history exists and is accessible
			if info, err := os.Stat(lastState.Directory); err == nil && info.IsDir() {
				s.Directory = lastState.Directory
				s.Filtered = PrepareDirItems(lastState.Directory)
				position = lastState.Position
				s.Position = lastState.Position
				items = s.Filtered

				// Delete the last state of the history
				s.History = s.History[:len(s.History)-1]
				s.DirScroll = 0 // Reset the scroll position
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
