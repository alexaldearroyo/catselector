package core

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// Function that generates the layout as a string instead of printing it
func DrawLayout() string {
	// Get the terminal width and height
	width, height := getTerminalSize()

	// Get the current directory
	dir := getCurrentDirectory()

	// Create the first line with proper spacing
	dirPrefix := "Directory: "
	titleText := "Cat Explorer"

	// Format the directory display
	var directoryDisplay string
	if len(dir) > width-len(dirPrefix)-len(titleText)-2 {
		// Split the directory into multiple lines if needed
		parts := splitDirectory(dir, width-len(dirPrefix)-2) // Leave space for prefix
		directoryDisplay = strings.Join(parts, "\n")
	} else {
		directoryDisplay = dir
	}

	// Create the header with proper alignment
	var header string

	// First line: Directory: prefix, directory, and title
	spacing := width - len(dirPrefix) - len(dir) - len(titleText)
	if spacing < 0 {
		spacing = 0
	}
	header = fmt.Sprintf("%s%s%s%s",
		DirectoryText.Render(dirPrefix),
		DirectoryDir.Render(directoryDisplay),
		strings.Repeat(" ", spacing),
		HeaderTitle.Render(titleText),
	)

	// Add a second empty line
	header += "\n"

	// Fill the rest of the screen with empty lines
	remainingHeight := height - 2 // -2 for the two header lines
	for i := 0; i < remainingHeight; i++ {
		header += "\n"
	}

	return header
}

// Helper function to split directory path into multiple lines
func splitDirectory(dir string, maxWidth int) []string {
	var parts []string
	current := dir

	for len(current) > maxWidth {
		// Find the last separator before maxWidth
		splitIndex := strings.LastIndex(current[:maxWidth], "/")
		if splitIndex == -1 {
			splitIndex = maxWidth - 1
		}
		parts = append(parts, current[:splitIndex+1])
		current = current[splitIndex+1:]
	}
	if len(current) > 0 {
		parts = append(parts, current)
	}

	return parts
}

// Get the terminal size
func getTerminalSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width, height = 80, 24 // Default terminal size if error occurs
	}
	return width, height
}

// Get the current directory
func getCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		return "/" // Return root if there's an error getting the current directory
	}
	return dir
}
