package core

import (
	"os"
	"path/filepath"
	"strings"
)

// GetFileIcon devuelve un icono apropiado según el tipo de archivo
func GetFileIcon(filePath string) string {
	// Verificar si es un directorio
	if info, err := os.Stat(filePath); err == nil && info.IsDir() {
		return "📁"
	}

	// Obtener la extensión del archivo
	ext := strings.ToLower(filepath.Ext(filePath))

	// Mapeo de extensiones a iconos
	switch ext {
	// Archivos de código
	case ".py", ".js", ".java", ".c", ".cpp", ".h", ".cs", ".php", ".rb", ".go", ".swift", ".kt", ".ts":
		return "📄"
	// Scripts
	case ".sh", ".bat", ".ps1", ".cmd":
		return "📜"
	// Archivos de texto
	case ".txt", ".md", ".rst", ".log":
		return "📝"
	// Archivos de configuración
	case ".json", ".yml", ".yaml", ".xml", ".ini", ".conf", ".cfg", ".toml":
		return "⚙️"
	// Imágenes
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".ico", ".tiff":
		return "🖼️"
	// Audio
	case ".mp3", ".wav", ".ogg", ".flac", ".aac":
		return "🎵"
	// Video
	case ".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv":
		return "🎬"
	// Archivos comprimidos
	case ".zip", ".rar", ".7z", ".tar", ".gz", ".bz2":
		return "📦"
	// Documentos
	case ".pdf":
		return "📕"
	case ".doc", ".docx":
		return "📘"
	case ".xls", ".xlsx":
		return "📗"
	case ".ppt", ".pptx":
		return "📙"
	// Ejecutables
	case ".exe", ".app", ".dmg", ".msi":
		return "🚀"
	// Archivos web
	case ".html", ".htm", ".css":
		return "🌐"
	// Archivos de git
	case ".git", ".gitignore":
		return "🐙"
	// Por defecto
	default:
		return "📄"
	}
}
