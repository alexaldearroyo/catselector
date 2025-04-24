package core

import (
	"os"
	"path/filepath"
	"strings"
)

// GetFileIcon devuelve un icono apropiado para el tipo de archivo
func GetFileIcon(filePath string) string {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "📄" // Archivo genérico si hay error
	}

	if fileInfo.IsDir() {
		return "📁" // Directorio
	}

	ext := strings.ToLower(filepath.Ext(filePath))

	// Archivos de código
	if contains([]string{".py", ".js", ".java", ".c", ".cpp", ".h", ".cs", ".php", ".rb", ".go", ".swift", ".kt", ".ts"}, ext) {
		return "📄" // Archivo de código
	}
	// Archivos de script
	if contains([]string{".sh", ".bat", ".ps1", ".cmd"}, ext) {
		return "📜" // Script
	}
	// Archivos de texto
	if contains([]string{".txt", ".md", ".rst", ".log"}, ext) {
		return "📝" // Texto
	}
	// Archivos de configuración
	if contains([]string{".json", ".yml", ".yaml", ".xml", ".ini", ".conf", ".cfg", ".toml"}, ext) {
		return "⚙️" // Configuración
	}
	// Archivos de imagen
	if contains([]string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".ico", ".tiff"}, ext) {
		return "🖼️" // Imagen
	}
	// Archivos de audio
	if contains([]string{".mp3", ".wav", ".ogg", ".flac", ".aac"}, ext) {
		return "🎵" // Audio
	}
	// Archivos de video
	if contains([]string{".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv"}, ext) {
		return "🎬" // Video
	}
	// Archivos comprimidos
	if contains([]string{".zip", ".rar", ".7z", ".tar", ".gz", ".bz2"}, ext) {
		return "📦" // Archivo comprimido
	}
	// Archivos PDF
	if contains([]string{".pdf"}, ext) {
		return "📕" // PDF
	}
	// Archivos de Word
	if contains([]string{".doc", ".docx"}, ext) {
		return "📘" // Word
	}
	// Archivos de Excel
	if contains([]string{".xls", ".xlsx"}, ext) {
		return "📗" // Excel
	}
	// Archivos de PowerPoint
	if contains([]string{".ppt", ".pptx"}, ext) {
		return "📙" // PowerPoint
	}
	// Archivos ejecutables
	if contains([]string{".exe", ".app", ".dmg", ".msi"}, ext) {
		return "🚀" // Ejecutable
	}
	// Archivos web
	if contains([]string{".html", ".htm", ".css"}, ext) {
		return "🌐" // Web
	}
	// Archivos Git
	if contains([]string{".git", ".gitignore"}, ext) {
		return "🐙" // Git
	}

	// Por defecto
	return "📄" // Archivo genérico
}

// contains comprueba si un string está en un slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
