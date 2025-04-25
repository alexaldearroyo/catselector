package core

import (
	"os"
	"path/filepath"
	"strings"
)

// GetFileIcon devuelve un icono apropiado según el tipo de archivo usando Nerd Fonts
func GetFileIcon(filePath string) string {
	// Verificar si es un directorio
	if info, err := os.Stat(filePath); err == nil && info.IsDir() {
		return "\uf07b" // 󰉋
	}

	// Obtener la extensión del archivo
	ext := strings.ToLower(filepath.Ext(filePath))

	// Mapeo de extensiones a iconos de Nerd Fonts
	switch ext {
	// Archivos de código
	case ".py":
		return "\ue235" // 󰆧
	case ".js":
		return "\ue74e" // 󰝎
	case ".java":
		return "\ue738" // 󰜸
	case ".c", ".cpp", ".h":
		return "\ue61d" // 󰘝
	case ".cs":
		return "\uf81a" // 󰠚
	case ".php":
		return "\ue73d" // 󰜽
	case ".rb":
		return "\ue21e" // 󰈞
	case ".go":
		return "\ue626" // 󰘦
	case ".swift":
		return "\ue755" // 󰝕
	case ".kt":
		return "\ue634" // 󰘴
	case ".ts":
		return "\ue628" // 󰘨
	// Scripts
	case ".sh", ".bat", ".ps1", ".cmd":
		return "\uf489" // 󰒉
	// Archivos de texto
	case ".txt", ".md", ".rst", ".log":
		return "\uf15c" // 󰅜
	// Archivos de configuración
	case ".json":
		return "\ue60b" // 󰘋
	case ".yml", ".yaml":
		return "\uf481" // 󰒁
	case ".xml":
		return "\uf72f" // 󰜯
	case ".ini", ".conf", ".cfg", ".toml":
		return "\uf013" // 󰀓
	// Imágenes
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".ico", ".tiff":
		return "\uf1c5" // 󰇅
	// Audio
	case ".mp3", ".wav", ".ogg", ".flac", ".aac":
		return "\uf001" // 󰀁
	// Video
	case ".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv":
		return "\uf03d" // 󰀽
	// Archivos comprimidos
	case ".zip", ".rar", ".7z", ".tar", ".gz", ".bz2":
		return "\uf1c6" // 󰇆
	// Documentos
	case ".pdf":
		return "\uf1c1" // 󰇁
	case ".doc", ".docx":
		return "\uf1c2" // 󰇂
	case ".xls", ".xlsx":
		return "\uf1c3" // 󰇃
	case ".ppt", ".pptx":
		return "\uf1c4" // 󰇄
	// Ejecutables
	case ".exe", ".app", ".dmg", ".msi":
		return "\uf2e0" // 󰋠
	// Archivos web
	case ".html", ".htm":
		return "\uf13b" // 󰄻
	case ".css":
		return "\ue42b" // 󰐫
	// Archivos de git
	case ".git", ".gitignore":
		return "\ue702" // 󰜂
	// Por defecto
	default:
		return "\uf15b" // 󰅛
	}
}
