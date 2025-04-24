package main

import (
	"os"
	"path/filepath"
	"strings"
)

// GetFileIcon devuelve el icono correspondiente según el tipo de archivo
func GetFileIcon(path string) string {
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		return "📁" // Directorio
	}

	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".py", ".js", ".java", ".c", ".cpp", ".h", ".cs", ".php", ".rb", ".go", ".swift", ".kt", ".ts":
		return "📄" // Código
	case ".sh", ".bat", ".ps1", ".cmd":
		return "📜" // Script
	case ".txt", ".md", ".rst", ".log":
		return "📝" // Texto
	case ".json", ".yml", ".yaml", ".xml", ".ini", ".conf", ".cfg", ".toml":
		return "⚙️" // Config
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".ico", ".tiff":
		return "🖼️" // Imagen
	case ".mp3", ".wav", ".ogg", ".flac", ".aac":
		return "🎵" // Audio
	case ".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv":
		return "🎬" // Video
	case ".zip", ".rar", ".7z", ".tar", ".gz", ".bz2":
		return "📦" // Comprimido
	case ".pdf":
		return "📕" // PDF
	case ".doc", ".docx":
		return "📘" // Word
	case ".xls", ".xlsx":
		return "📗" // Excel
	case ".ppt", ".pptx":
		return "📙" // PowerPoint
	case ".exe", ".app", ".dmg", ".msi":
		return "🚀" // Ejecutable
	case ".html", ".htm", ".css":
		return "🌐" // Web
	case ".git", ".gitignore":
		return "🐙" // Git
	default:
		return "📄" // Genérico
	}
}
