package core

import (
	"os"
	"path/filepath"
	"strings"
)

// Recursive function to process directories and files
func processDirectoryRecursive(selector *Selector, dirPath string, item string, selectState bool) {
	// Update the selection state of the current directory
	selectionKey := selector.GetSelectionKey(item)
	selector.Selection[selectionKey] = selectState

	// If we are deselecting, clear all files and subdirectories
	// of the current directory from the selection map
	if !selectState {
		// Delete selection entries related to this directory
		prefix := dirPath + string(os.PathSeparator)
		for path := range selector.Selection {
			if strings.HasPrefix(path, prefix) {
				selector.Selection[path] = false
			}
		}
		return
	}

	// If we are selecting and the include mode is active,
	// process recursively the subdirectories
	if selector.IncludeMode {
		// Read the content of the directory
		entries, err := os.ReadDir(dirPath)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					// Process recursively the subdirectory
					subDirPath := filepath.Join(dirPath, entry.Name())
					subItem := filepath.Join(item, entry.Name())
					processDirectoryRecursive(selector, subDirPath, subItem, selectState)
				}
				// We don't mark the files individually because the parent directory
				// is already selected and the IsFileSelected function will check this
			}
		}
	}
}
