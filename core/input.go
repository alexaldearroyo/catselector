package core

func CaptureInput(key string) string {
	// Aquí manejas las teclas que recibes
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
	case "down", "j":
		position++
		if position >= itemCount {
			position = 0
		}
	case "up", "k":
		position--
		if position < 0 {
			position = itemCount - 1
		}
	}

	// Actualizar la posición en el selector
	s.Position = position
	s.Filtered = items

	// Actualizar los archivos cuando se navega
	s.UpdateFilesForCurrentDirectory()

	// Actualizar la selección
	if position >= 0 && position < len(items) {
		selectedItem := items[position]
		selected[selectedItem] = !selected[selectedItem]
	}

	return position
}
