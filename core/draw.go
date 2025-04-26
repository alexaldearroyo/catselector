package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func DrawLayout(position int, items []string, currentDir string, files []string, activePanel int, filePosition int) string {
	width, height := getTerminalSize()
	dirPrefix := "Directory: "
	titleText := "Cat Explorer"
	minSpacing := 2

	available := width - len(dirPrefix) - len(titleText) - minSpacing
	narrow := len(currentDir) > available

	// 1st line: prefix + spaces + title
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

	// 2nd line: parentDir + "/" + lastPart
	if narrow {
		// Split the directory into parts
		parts := strings.Split(currentDir, "/")
		lastPart := parts[len(parts)-1]
		parentDir := strings.Join(parts[:len(parts)-1], "/")
		if parentDir != "" {
			parentDir += "/"
		}
		header += "\n" + DirectoryText.Render(parentDir) + DirectoryDir.Render(lastPart)

		// Add new_text in a separate line after the pwd
		// Replace with subdirectory text and selection
		selector := GetCurrentSelector()
		var subdirText, selectedText string

		// Determine the state of the subdirectories
		if selector.IncludeMode {
			subdirText = White.Render("Subdirectories: ") + Magenta.Render("Included")
		} else {
			subdirText = White.Render("Subdirectories: ") + Magenta.Render("Not included")
		}

		// Count selected files and directories
		selectedFiles, selectedDirs := countSelected(selector)
		selectedText = White.Render("Selected: ") +
			Magenta.Render(fmt.Sprintf("%d", selectedFiles)) +
			White.Render(" Files") +
			White.Render(", ") +
			Magenta.Render(fmt.Sprintf("%d", selectedDirs)) +
			White.Render(" Directories")

		// Full text with Selected after Included/Not included
		infoText := subdirText + "   " + selectedText
		header += "\n" + infoText
	} else {
		// Split the directory into parts
		parts := strings.Split(currentDir, "/")
		lastPart := parts[len(parts)-1]
		parentDir := strings.Join(parts[:len(parts)-1], "/")
		if parentDir != "" {
			parentDir += "/"
		}
		inLineSpaces := width - len(dirPrefix) - len(currentDir) - len(titleText)
		header = fmt.Sprintf(
			"%s%s%s%s%s",
			DirectoryText.Render(dirPrefix),
			DirectoryText.Render(parentDir),
			DirectoryDir.Render(lastPart),
			strings.Repeat(" ", inLineSpaces),
			HeaderTitle.Render(titleText),
		)

		// Add a new line and then the information text
		// Replace with subdirectory text and selection
		selector := GetCurrentSelector()
		var subdirText, selectedText string

		// Determinar el estado de los subdirectorios
		if selector.IncludeMode {
			subdirText = White.Render("Subdirectories: ") + Magenta.Render("Included")
		} else {
			subdirText = White.Render("Subdirectories: ") + White.Render("Not included")
		}

		// Count selected files and directories
		selectedFiles, selectedDirs := countSelected(selector)
		selectedText = White.Render("Selected: ") +
			Magenta.Render(fmt.Sprintf("%d", selectedFiles)) +
			White.Render(" Files") +
			White.Render(", ") +
			Magenta.Render(fmt.Sprintf("%d", selectedDirs)) +
			White.Render(" Directories")

		// Full text with Selected after Included/Not included
		infoText := subdirText + "   " + selectedText
		header += "\n" + infoText
	}

	// Add a final new line before the panels
	header += "\n"

	// Get the current selector
	selector := GetCurrentSelector()

	// Panel layout
	panelWidth := width / 3

	renderLeft := func(text string, isActive bool, isCounter bool) string {
		return RenderLeft(text, isActive, isCounter, panelWidth)
	}

	// Add a white divider line above the panel headers
	header += White.Render(strings.Repeat("─", width)) + "\n"

	// Count elements for each panel
	var totalItems, totalFiles, totalSubdirs int
	var err error

	// Determine the directory for the count
	var countDir string
	if position >= 0 && position < len(items) {
		// If there is a selected directory, use that for the count
		item := items[position]
		if item == ".." {
			countDir = filepath.Dir(currentDir)
		} else if item == "." {
			countDir = currentDir
		} else {
			countDir = filepath.Join(currentDir, item)
		}
	} else {
		// If there is no selected directory, use the current directory
		countDir = currentDir
	}

	// Count elements in the determined directory
	totalItems, err = countItems(countDir)
	if err == nil {
		totalFiles, _ = countFiles(countDir)
		totalSubdirs, _ = countSubdirs(countDir)
	}

	// Add counters to the headers
	left := renderLeft("Directories", activePanel == 1, false)
	middle := renderLeft("Files", activePanel == 2, false)
	right := renderLeft("Preview", activePanel == 3, false)

	// Add counters to a new line
	leftCounter := renderLeft(fmt.Sprintf("%d items", totalItems), false, true)
	middleCounter := renderLeft(fmt.Sprintf("%d files", totalFiles), false, true)

	// Determinar el texto para el panel Preview
	var rightCounter string
	if activePanel == 2 && filePosition >= 0 && filePosition < len(files) {
		// Si estamos en el panel Files, mostrar el nombre del archivo
		rightCounter = renderLeft(files[filePosition], false, true)
	} else {
		// Si no, mostrar el contador de subdirectorios
		rightCounter = renderLeft(fmt.Sprintf("%d subdirs", totalSubdirs), false, true)
	}

	// Combine headers and counters
	header += left + White.Render("│") + middle + White.Render("│") + right + "\n"
	header += leftCounter + White.Render("│") + middleCounter + White.Render("│") + rightCounter + "\n"

	// Left panel (Directories)
	selected := map[string]bool{}
	start := 0
	panelHeight := height - 10
	active := activePanel == 1
	includeSubdirs := false

	leftPanel := renderLeftPanel(items, selected, currentDir, position, start, panelHeight, panelWidth, active, includeSubdirs)

	// Files panel
	filePanel := renderFilePanel(files, position, panelWidth, height, panelHeight, activePanel, filePosition)

	// Right panel (Preview)
	rightPanel := renderPreviewPanel(currentDir, panelWidth, panelHeight, files, filePosition, activePanel, items, position)

	// Combine the panels horizontally
	var result strings.Builder
	result.WriteString(header)

	// Split the panels into lines
	leftLines := strings.Split(leftPanel, "\n")
	fileLines := strings.Split(filePanel, "\n")
	rightLines := strings.Split(rightPanel, "\n")

	// Find the maximum number of lines
	maxLines := len(leftLines)
	if len(fileLines) > maxLines {
		maxLines = len(fileLines)
	}
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}

	// Ensure we have enough lines
	for maxLines < panelHeight {
		leftLines = append(leftLines, "")
		fileLines = append(fileLines, "")
		rightLines = append(rightLines, "")
		maxLines++
	}

	// Combine the lines horizontally
	for i := 0; i < maxLines; i++ {
		leftLine := ""
		if i < len(leftLines) {
			leftLine = leftLines[i]
		}
		fileLine := ""
		if i < len(fileLines) {
			fileLine = fileLines[i]
		}
		rightLine := ""
		if i < len(rightLines) {
			rightLine = rightLines[i]
		}

		// Ensure each line has the correct width
		leftLine = lipgloss.PlaceHorizontal(panelWidth, lipgloss.Left, leftLine)
		fileLine = lipgloss.PlaceHorizontal(panelWidth, lipgloss.Left, fileLine)
		rightLine = lipgloss.PlaceHorizontal(panelWidth, lipgloss.Left, rightLine)

		// Add white vertical lines between the panels
		result.WriteString(leftLine + White.Render("│") + fileLine + White.Render("│") + rightLine + "\n")
	}

	// Add the status bar at the bottom
	statusBar := strings.Repeat("─", width)
	if selector != nil && selector.StatusMessage != "" && time.Now().Unix()-selector.StatusTime < 3 {
		// Show the message for 3 seconds
		statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
		statusBar = statusStyle.Render(selector.StatusMessage)
	} else {
		// Show an empty status bar
		statusBar = strings.Repeat(" ", width)
	}

	// Key bindings menu (two lines, evenly distributed and colored)
	keyBindings := []struct {
		Key, Desc string
	}{
		{"k/j", "Up or Down"},
		{"Enter/l", "Enter"},
		{"Esc/h", "Back"},
		{"o/c", "Open or Copy"},
		{"s/a", "Select or All"},
		{"i", "Include"},
		{"Tab", "Change Panel"},
		{"q", "Quit"},
	}

	// Calculate the available width and the number of shortcuts per line
	numPerLine := 4 // 4 elementos por fila
	width, _ = getTerminalSize()

	var line1, line2 string
	for i, kb := range keyBindings {
		// Format the text with the requested colors
		keyText := Blue.Render(kb.Key + ":")
		descText := White.Render(" " + kb.Desc)
		combo := keyText + descText

		// Calculate the space to distribute evenly
		space := (width / numPerLine) - lipgloss.Width(combo)
		padding := strings.Repeat(" ", space)

		if i < numPerLine {
			line1 += combo + padding
		} else {
			line2 += combo + padding
		}
	}

	// Ensure both lines have the same width
	if lipgloss.Width(line1) < width {
		line1 += strings.Repeat(" ", width-lipgloss.Width(line1))
	}
	if lipgloss.Width(line2) < width {
		line2 += strings.Repeat(" ", width-lipgloss.Width(line2))
	}

	// Calcula la línea divisoria antes de usarla
	divider := White.Render(strings.Repeat("─", width))
	result.WriteString(divider + "\n")
	result.WriteString(line1 + "\n")
	result.WriteString(line2 + "\n")
	// if !strings.HasSuffix(statusBar, "\n") {
	// 	statusBar += "\n"
	// }
	result.WriteString(statusBar)

	return result.String()
}

// This function should handle the rendering of the files
func renderFilePanel(files []string, position, panelWidth, height, panelHeight int, activePanel int, filePosition int) string {
	var b strings.Builder

	// Get the current selector
	selector := GetCurrentSelector()

	// Check if the current directory or the selected directory is selected
	currentDirSelected := selector.IsSelected(".")
	selectedDirSelected := false
	if selector.ActivePanel == 1 && selector.Position < len(selector.Filtered) {
		selectedItem := selector.Filtered[selector.Position]
		if selectedItem != ".." {
			selectedDirSelected = selector.IsSelected(selectedItem)
		}
	}

	for i := 0; i < panelHeight && i < len(files); i++ {
		file := files[i]
		icon := GetFileIcon(file)

		// Check if the file is selected or if any parent directory is selected
		isSelected := selector.IsFileSelected(file) || currentDirSelected || selectedDirSelected

		// Add asterisk if it is selected
		marker := "  "
		if isSelected {
			marker = " *"
		}

		line := icon + marker + file

		// Truncate the file name if it is too long
		maxWidth := panelWidth - 2 // Leave space for the scrollbar
		if lipgloss.Width(line) > maxWidth {
			// Calculate how much space we have for the file name
			iconWidth := lipgloss.Width(icon + marker)
			availableWidth := maxWidth - iconWidth - 3 // 3 for "..."

			// Truncate the file name
			if availableWidth > 0 {
				truncatedName := file
				if len(truncatedName) > availableWidth {
					truncatedName = truncatedName[:availableWidth] + "..."
				}
				line = icon + marker + truncatedName
			}
		}

		// Apply the Focus style if the panel is active and this is the selected file
		if activePanel == 2 && i == filePosition {
			// Pad the line to the panel width
			padding := panelWidth - lipgloss.Width(line)
			if padding > 0 {
				line += strings.Repeat(" ", padding)
			}
			b.WriteString(Focus.Render(line) + "\n")
		} else if isSelected {
			// Apply yellow style for selected files
			b.WriteString(Yellow.Render(line) + "\n")
		} else {
			b.WriteString(White.Render(line) + "\n")
		}
	}

	return b.String()
}

// renderPreviewPanel shows the content of the selected file or the subdirectories
func renderPreviewPanel(dir string, width, height int, files []string, filePosition int, activePanel int, items []string, position int) string {
	var b strings.Builder

	// Get the current selector
	selector := GetCurrentSelector()

	// If the active panel is the directories panel, show subdirectories
	if activePanel == 1 {
		// Determine the selected directory
		var selectedDir string
		if position >= 0 && position < len(items) {
			item := items[position]
			if item == ".." {
				selectedDir = filepath.Dir(dir)
			} else if item == "." {
				selectedDir = dir
			} else {
				selectedDir = filepath.Join(dir, item)
			}
		} else {
			selectedDir = dir
		}

		// Check if the directory exists and is accessible
		info, err := os.Stat(selectedDir)
		if err != nil || !info.IsDir() {
			return strings.Repeat(" ", width) + "\n"
		}

		// Read the subdirectories
		entries, err := os.ReadDir(selectedDir)
		if err != nil {
			return strings.Repeat(" ", width) + "\n"
		}

		// Filter only directories and sort them
		var subdirs []string
		for _, entry := range entries {
			if entry.IsDir() {
				subdirs = append(subdirs, entry.Name())
			}
		}
		sort.Strings(subdirs)

		// Show the subdirectories
		for i := 0; i < height && i < len(subdirs); i++ {
			subdir := subdirs[i]
			icon := GetFileIcon(filepath.Join(selectedDir, subdir))

			// Check if the subdirectory is selected
			isSelected := selector.Selection[filepath.Join(selectedDir, subdir)]

			// Add the corresponding marker
			marker := "  "
			if isSelected {
				marker = " •"
			}

			line := icon + marker + subdir

			// Pad the line to the panel width
			padding := width - lipgloss.Width(line)
			if padding > 0 {
				line += strings.Repeat(" ", padding)
			}

			// Apply the corresponding style
			if isSelected {
				b.WriteString(Yellow.Render(line) + "\n")
			} else {
				b.WriteString(Green.Render(line) + "\n")
			}
		}

		// Pad with empty lines if necessary
		for i := len(subdirs); i < height; i++ {
			b.WriteString(strings.Repeat(" ", width) + "\n")
		}
	} else {
		// If we are in the files panel, show the content of the selected file
		if len(files) > 0 && filePosition >= 0 && filePosition < len(files) {
			// Determine the current directory for the files
			var currentDir string
			if position >= 0 && position < len(items) {
				item := items[position]
				if item == ".." {
					currentDir = filepath.Dir(dir)
				} else if item == "." {
					currentDir = dir
				} else {
					currentDir = filepath.Join(dir, item)
				}
			} else {
				currentDir = dir
			}

			filePath := filepath.Join(currentDir, files[filePosition])

			// Check if the file is binary or too large
			info, err := os.Stat(filePath)
			if err != nil {
				// Error accessing the file
				ShowErrorMessage(&b, "Cannot access the file", filePath, width, height)
				return b.String()
			}

			// Check if the file is binary or too large
			if IsBinaryFile(filePath) || info.Size() > 1024*1024 { // More than 1MB
				ShowBinaryFileMessage(&b, filePath, info.Size(), width, height)
				return b.String()
			}

			// Read the content of the file
			content, err := os.ReadFile(filePath)
			if err == nil {
				// Convert the content to string and limit to the first lines that fit
				lines := strings.Split(string(content), "\n")

				// Limit the number of lines to avoid overflows
				maxLines := height
				if len(lines) > maxLines {
					lines = lines[:maxLines]
				}

				for i := 0; i < len(lines); i++ {
					line := lines[i]

					// Sanitize the line to avoid problematic characters
					line = SanitizeLine(line)

					// Truncate the line if it is too long
					if lipgloss.Width(line) > width {
						line = line[:width-3] + "..."
					}

					// Pad the line to the panel width
					padding := width - lipgloss.Width(line)
					if padding > 0 {
						line += strings.Repeat(" ", padding)
					}

					b.WriteString(White.Render(line) + "\n")
				}

				// Pad with empty lines if necessary
				for i := len(lines); i < height; i++ {
					b.WriteString(strings.Repeat(" ", width) + "\n")
				}
			} else {
				// Error reading the file
				ShowErrorMessage(&b, "Cannot read the file", filePath, width, height)
			}
		} else {
			// Message when no file is selected
			msg := "No file selected"
			padding := width - lipgloss.Width(msg)
			if padding > 0 {
				msg += strings.Repeat(" ", padding)
			}
			b.WriteString(White.Render(msg) + "\n")

			// Pad with empty lines
			for i := 1; i < height; i++ {
				b.WriteString(strings.Repeat(" ", width) + "\n")
			}
		}
	}

	return b.String()
}

// Get the terminal size
func getTerminalSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width, height = 80, 24 // Default terminal size if error occurs
	}
	return width, height
}

func renderLeftPanel(items []string, selected map[string]bool, directory string, position, start, height, width int, active bool, includeSubdirs bool) string {
	var b strings.Builder

	end := start + height
	if end > len(items) {
		end = len(items)
	}

	// Get the current selector
	selector := GetCurrentSelector()

	for i := start; i < end; i++ {
		item := items[i]
		fullPath := filepath.Join(directory, item)
		if item == ".." {
			fullPath = filepath.Dir(directory)
		}
		isSelected := false
		if item != ".." {
			itemKey := selector.GetSelectionKey(item)
			isSelected = selector.Selection[itemKey]
		}
		hasFocus := active && i == position

		marker := "  "
		if isSelected {
			marker = " •"
		}
		content := marker + item
		maxWidth := width - 3
		if lipgloss.Width(content) > maxWidth {
			content = content[:maxWidth-3] + "..."
		}

		icon := GetFileIcon(fullPath)
		line := icon + content

		// Pad the line to the panel width
		padding := width - lipgloss.Width(line)
		if padding > 0 {
			line += strings.Repeat(" ", padding)
		}

		// Styles
		if hasFocus {
			b.WriteString(Focus.Render(line) + "\n")
		} else if isSelected {
			b.WriteString(Yellow.Render(line) + "\n")
		} else {
			b.WriteString(Green.Render(line) + "\n")
		}
	}

	return b.String()
}

// Global variable to maintain the current selector
var currentSelector *Selector

// Function to set the current selector
func SetCurrentSelector(s *Selector) {
	currentSelector = s
}

// Function to get the current selector
func GetCurrentSelector() *Selector {
	if currentSelector == nil {
		currentSelector = &Selector{
			Selection: make(map[string]bool),
		}
	}
	return currentSelector
}

// countItems counts the total number of elements (files + subdirectories) in a directory
func countItems(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	return len(entries), nil
}

// countFiles counts the number of files in a directory
func countFiles(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			count++
		}
	}
	return count, nil
}

// countSubdirs counts the number of subdirectories in a directory
func countSubdirs(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			count++
		}
	}
	return count, nil
}

// countSelected counts the number of selected files and directories
func countSelected(selector *Selector) (int, int) {
	selectedFiles := 0
	selectedDirs := 0
	processedDirs := make(map[string]bool)

	for key, selected := range selector.Selection {
		if !selected {
			continue
		}

		// Check if it is a directory or a file
		fileInfo, err := os.Stat(key)
		if err != nil {
			continue
		}

		if fileInfo.IsDir() {
			selectedDirs++
			// Si el directorio ya fue procesado, lo saltamos
			if processedDirs[key] {
				continue
			}
			processedDirs[key] = true

			// Contar archivos en el directorio y sus subdirectorios
			if selector.IncludeMode {
				// Contar recursivamente todos los archivos en subdirectorios
				filepath.Walk(key, func(p string, info os.FileInfo, err error) error {
					if err == nil && !info.IsDir() {
						selectedFiles++
					}
					return nil
				})
			} else {
				// Contar solo archivos en el nivel superior del directorio
				files, err := os.ReadDir(key)
				if err == nil {
					for _, file := range files {
						fileInfo, err := file.Info()
						if err == nil && !fileInfo.IsDir() {
							selectedFiles++
						}
					}
				}
			}
		} else {
			selectedFiles++
		}
	}

	return selectedFiles, selectedDirs
}
