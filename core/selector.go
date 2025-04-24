package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Selector es la estructura principal que gestiona la selección de archivos
type Selector struct {
	directory            string
	initialDirectory     string
	files                []string
	directories          []string
	selected             map[string]bool
	manuallyDeselected   map[string]bool
	position             int
	windowStart          int
	searchMode           bool
	searchTerm           string
	includeSubdirectories bool
	filtered             []string
	activePanel          int  // 1 = directories, 2 = files, 3 = preview
	previewContent       []string
	previewPosition      int
	previewWindowStart   int
	currentPreviewDirectory string
	lastDirectoryName    string
	statusMessage        string
	statusTime           float64

	// Componentes de la interfaz
	app                  *tview.Application
	mainGrid             *tview.Grid
	directoriesPanel     *tview.List
	filesPanel           *tview.List
	previewPanel         *tview.TextView
	statusBar            *tview.TextView
	headerPanel          *tview.TextView
	keybindingsPanel     *tview.TextView
}

// NewSelector crea una nueva instancia de Selector
func NewSelector(directory string) *Selector {
	s := &Selector{
		directory:            directory,
		initialDirectory:     directory,
		selected:             make(map[string]bool),
		manuallyDeselected:   make(map[string]bool),
		activePanel:          1, // Empezar en el panel de directorios
		includeSubdirectories: false,
	}
	return s
}

// Run inicia la interfaz del selector
func (s *Selector) Run() error {
	// Inicializar la interfaz gráfica con tview
	s.app = tview.NewApplication()

	// Configurar los paneles
	s.setupUI()

	// Cargar archivos iniciales
	s.loadFiles()

	// Iniciar la aplicación
	if err := s.app.Run(); err != nil {
		return err
	}

	return nil
}

// LoadFiles carga los archivos y directorios en el directorio actual
func (s *Selector) loadFiles() {
	s.files = []string{}
	s.directories = []string{}

	// Obtener elementos del directorio actual
	files, err := os.ReadDir(s.directory)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			s.directories = append(s.directories, file.Name())
		} else {
			s.files = append(s.files, file.Name())
		}
	}

	// Ordenar alfabéticamente
	sort.Strings(s.directories)
	sort.Strings(s.files)

	// Preparar lista filtrada inicial
	s.updateFiltered()
}

// UpdateFiltered actualiza la lista filtrada según el término de búsqueda y panel activo
func (s *Selector) updateFiltered() {
	if s.activePanel == 1 { // Panel de directorios
		var baseItems []string
		if s.directory == s.initialDirectory {
			baseItems = append([]string{"."}, s.directories...)
		} else {
			baseItems = append([]string{"..", "."}, s.directories...)
		}

		// Actualizar directorio actual para el panel de archivos
		if len(s.filtered) > 0 && s.position < len(s.filtered) {
			currentItem := s.filtered[s.position]
			if currentItem != ".." && currentItem != "." {
				s.currentPreviewDirectory = filepath.Join(s.directory, currentItem)
			} else if currentItem == ".." {
				s.currentPreviewDirectory = filepath.Dir(s.directory)
			} else { // Es "."
				s.currentPreviewDirectory = s.directory
			}
		}

		s.filtered = baseItems
	} else { // Panel de archivos
		if s.currentPreviewDirectory != "" && isDir(s.currentPreviewDirectory) {
			try := func() []string {
				files, err := os.ReadDir(s.currentPreviewDirectory)
				if err != nil {
					return []string{}
				}

				var result []string
				for _, fileInfo := range files {
					if !fileInfo.IsDir() {
						result = append(result, fileInfo.Name())
					}
				}
				sort.Strings(result)
				return result
			}

			s.files = try()
		}
		s.filtered = s.files
	}

	// Filtrar por término de búsqueda si está en modo búsqueda
	if s.searchMode && s.searchTerm != "" {
		var result []string
		for _, item := range s.filtered {
			if strings.Contains(strings.ToLower(item), strings.ToLower(s.searchTerm)) {
				result = append(result, item)
			}
		}
		s.filtered = result
	}

	// Resetear posición si está fuera de rango
	if len(s.filtered) == 0 {
		s.position = 0
	} else if s.position >= len(s.filtered) {
		s.position = len(s.filtered) - 1
	}

	// Actualizar archivos seleccionados en el panel de archivos
	if s.activePanel == 2 && s.currentPreviewDirectory != "" {
		files, err := os.ReadDir(s.currentPreviewDirectory)
		if err == nil {
			s.files = []string{}
			for _, file := range files {
				if !file.IsDir() {
					s.files = append(s.files, file.Name())
				}
			}
			sort.Strings(s.files)
			s.filtered = s.files
		}
	}
}

// GetPreviewContent obtiene el contenido de vista previa para un archivo o directorio
func (s *Selector) getPreviewContent(item string) []string {
	preview := []string{}

	if s.activePanel == 1 { // Panel de directorios
		// Obtener contenido del directorio seleccionado
		var dirPath string
		if item == ".." {
			dirPath = filepath.Dir(s.directory)
		} else if item == "." {
			dirPath = s.directory
		} else {
			dirPath = filepath.Join(s.directory, item)
		}

		files, err := os.ReadDir(dirPath)
		if err != nil {
			return []string{fmt.Sprintf("Error accediendo al directorio: %s", err.Error())}
		}

		// Separar directorios y archivos
		var dirs []string
		var fileList []string
		for _, fileInfo := range files {
			if fileInfo.IsDir() {
				dirs = append(dirs, fileInfo.Name())
			} else {
				fileList = append(fileList, fileInfo.Name())
			}
		}

		// Ordenar ambas listas
		sort.Strings(dirs)
		sort.Strings(fileList)

		// Mostrar solo subdirectorios cuando el foco está en Directorios
		for _, dirItem := range dirs {
			icon := "📁"
			fullItemPath := filepath.Join(dirPath, dirItem)
			// Verificar si la ruta está seleccionada
			isSelected := s.selected[fullItemPath]

			marker := " *"
			if isSelected {
				marker = " *"
			} else {
				marker = "  "
			}
			preview = append(preview, fmt.Sprintf("%s %s%s", icon, marker, dirItem))
		}
	} else { // Panel de archivos
		if item == "" {
			return []string{"No hay archivo seleccionado"}
		}

		// Corregir la ruta del archivo en el panel Files
		var filePath string
		if s.activePanel == 2 && s.currentPreviewDirectory != "" {
			filePath = filepath.Join(s.currentPreviewDirectory, item)
		} else {
			filePath = filepath.Join(s.directory, item)
		}

		// Intentar mostrar vista previa del contenido para archivos de texto
		textExtensions := []string{".txt", ".py", ".js", ".java", ".c", ".cpp", ".h", ".html", ".css",
			".json", ".xml", ".md", ".sh", ".bat", ".ps1", ".yaml", ".yml",
			".ini", ".cfg", ".conf", ".rb", ".go"}

		isText := false
		for _, ext := range textExtensions {
			if strings.HasSuffix(strings.ToLower(filePath), ext) {
				isText = true
				break
			}
		}

		if isText {
			content, err := os.ReadFile(filePath)
			if err != nil {
				return []string{fmt.Sprintf("Error leyendo archivo: %s", err.Error())}
			}

			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if len(line) > 80 {
					line = line[:77] + "..."
				}
				preview = append(preview, line)
			}
		} else {
			preview = append(preview, "Archivo binario")
			preview = append(preview, "Vista previa no disponible")
		}
	}

	return preview
}

// GetItemIcon devuelve el icono apropiado para un elemento
func (s *Selector) getItemIcon(item string, isInDirectoryPanel bool) string {
	if item == ".." && isInDirectoryPanel {
		return "←" // Flecha izquierda para '..'
	}

	path := filepath.Join(s.directory, item)
	info, err := os.Stat(path)
	if err != nil {
		return "📄" // Por defecto si hay error
	}

	if info.IsDir() {
		return "📁"
	}

	return GetFileIcon(path)
}

// SetupUI configura la interfaz de usuario
func (s *Selector) setupUI() {
	// Crear componentes principales
	s.mainGrid = tview.NewGrid()
	s.directoriesPanel = tview.NewList().ShowSecondaryText(false)
	s.filesPanel = tview.NewList().ShowSecondaryText(false)
	s.previewPanel = tview.NewTextView().SetDynamicColors(true)
	s.statusBar = tview.NewTextView().SetDynamicColors(true)
	s.headerPanel = tview.NewTextView().SetDynamicColors(true)
	s.keybindingsPanel = tview.NewTextView().SetDynamicColors(true)

	// Configurar el panel de directorios
	s.directoriesPanel.SetBorder(true).SetTitle("Directorios [d]")
	s.filesPanel.SetBorder(true).SetTitle("Archivos [f]")
	s.previewPanel.SetBorder(true).SetTitle("Vista previa")

	// Configurar colores
	s.directoriesPanel.SetTitleColor(tcell.ColorBlue)
	s.filesPanel.SetTitleColor(tcell.ColorBlue)
	s.previewPanel.SetTitleColor(tcell.ColorBlue)

	// Configurar el grid principal
	s.mainGrid.SetRows(3, -1, 3) // Header, contenido, footer
	s.mainGrid.SetColumns(-1, -1, -1) // Tres columnas iguales

	// Asignar componentes al grid
	s.mainGrid.AddItem(s.headerPanel, 0, 0, 1, 3, 0, 0, false) // Header
	s.mainGrid.AddItem(s.directoriesPanel, 1, 0, 1, 1, 0, 0, true) // Directorios
	s.mainGrid.AddItem(s.filesPanel, 1, 1, 1, 1, 0, 0, false) // Archivos
	s.mainGrid.AddItem(s.previewPanel, 1, 2, 1, 1, 0, 0, false) // Vista previa
	s.mainGrid.AddItem(s.keybindingsPanel, 2, 0, 1, 2, 0, 0, false) // Teclas
	s.mainGrid.AddItem(s.statusBar, 2, 2, 1, 1, 0, 0, false) // Barra de estado

	// Configurar manejo de eventos
	s.app.SetInputCapture(s.handleInput)

	// Establecer componente root
	s.app.SetRoot(s.mainGrid, true)
}

// FocusDot establece el foco en el directorio actual
func (s *Selector) focusDot() {
	if s.directory != s.initialDirectory {
		s.position = 1
	} else {
		s.position = 0
	}
}

// isDir comprueba si una ruta es un directorio
func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
