package core

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// Function that generates the layout as a string instead of printing it
func DrawLayout() string {
	width, height := getTerminalSize()
	dir           := getCurrentDirectory()

	dirPrefix  := "Directory: "
	titleText  := "Cat Explorer"
	minSpacing := 2

	available := width - len(dirPrefix) - len(titleText) - minSpacing
	narrow    := len(dir) > available

	// 1ª línea: prefijo + espacios + título
	spaces := width - len(dirPrefix) - len(titleText)
	if spaces < minSpacing {
			spaces = minSpacing
	}
	header := fmt.Sprintf(
			"%s%s%s",
			DirectoryText.Render(dirPrefix),
			strings.Repeat(" ", spaces),
			HeaderTitle.Render(titleText),
	)

	// 2ª y 3ª líneas en pantalla estrecha
	if narrow {
			// 2ª línea: Directory: + título alineado a la derecha
			spaceCount := width - len(dirPrefix) - len(titleText)
			header += "\n" +
					DirectoryText.Render(dirPrefix) +
					strings.Repeat(" ", spaceCount) +
					HeaderTitle.Render(titleText)

			// 3ª línea: directorio en sí
			header += "\n" + DirectoryDir.Render(dir)
	} else {
			// todo en una línea si cabe
			inLineSpaces := width - len(dirPrefix) - len(dir) - len(titleText)
			header = fmt.Sprintf(
					"%s%s%s%s",
					DirectoryText.Render(dirPrefix),
					DirectoryDir.Render(dir),
					strings.Repeat(" ", inLineSpaces),
					HeaderTitle.Render(titleText),
			)
	}

	// relleno hasta height
	header += "\n"
	for i := 0; i < height-3; i++ {
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
