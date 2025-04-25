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


func HandleKeyPress(key string, position, itemCount int, selected map[string]bool, items []string) int {
	switch key {
	case "down":
		position++
		if position >= itemCount { // Si llegamos al final, volvemos al primero
			position = 0
		}
	case "up":
		position--
		if position < 0 { // Si estamos al principio, saltamos al último
			position = itemCount - 1
		}
	case "j":
		position++
		if position >= itemCount {
			position = 0
		}
	case "k":
		position--
		if position < 0 {
			position = itemCount - 1
		}
	}

	// Marcar el item seleccionado
	selectedItem := items[position]
	selected[selectedItem] = !selected[selectedItem]

	return position
}
