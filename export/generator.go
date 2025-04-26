package export

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// GenerateTextFile genera un archivo de texto con el contenido de los archivos seleccionados
func GenerateTextFile(selected []string, excluded []string, includeSubdirs bool, initialDir string, currentDir string) string {
	// Crear un hash único basado en la hora actual y los archivos seleccionados
	selectedCopy := make([]string, len(selected))
	copy(selectedCopy, selected)
	sort.Strings(selectedCopy)

	hashInput := fmt.Sprintf("%d%v", time.Now().Unix(), selectedCopy)
	hasher := md5.New()
	hasher.Write([]byte(hashInput))
	hashValue := hex.EncodeToString(hasher.Sum(nil))[:8]

	filesToProcess := []string{}

	// Crear map para búsqueda rápida de exclusiones
	excludedMap := make(map[string]bool)
	for _, path := range excluded {
		excludedMap[path] = true
	}

	// Recopilar archivos a procesar
	for _, path := range selected {
		if excludedMap[path] {
			continue
		}

		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		if info.IsDir() {
			if includeSubdirs {
				// Recorrer el directorio recursivamente
				filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
					if err != nil {
						return nil // Continuar con el siguiente archivo
					}

					if !fileInfo.IsDir() && !excludedMap[filePath] {
						filesToProcess = append(filesToProcess, filePath)
					}
					return nil
				})
			} else {
				// Recorrer solo el nivel superior del directorio
				files, err := os.ReadDir(path)
				if err == nil {
					for _, file := range files {
						filePath := filepath.Join(path, file.Name())

						fileInfo, err := file.Info()
						if err != nil {
							continue
						}

						if !fileInfo.IsDir() && !excludedMap[filePath] {
							filesToProcess = append(filesToProcess, filePath)
						}
					}
				}
			}
		} else if !info.IsDir() {
			filesToProcess = append(filesToProcess, path)
		}
	}

	if len(filesToProcess) == 0 {
		return ""
	}

	// Crear el archivo de salida
	outputFile := filepath.Join(initialDir, fmt.Sprintf("cs_%s.txt", hashValue))
	file, err := os.Create(outputFile)
	if err != nil {
		return ""
	}
	defer file.Close()

	// Escribir el contenido de cada archivo
	for _, filePath := range filesToProcess {
		relPath, err := filepath.Rel(currentDir, filePath)
		if err != nil {
			relPath = filePath
		}
		relativeName := filepath.ToSlash(relPath)

		file.WriteString("---------------------------------------------\n")
		file.WriteString(fmt.Sprintf("// File %s\n", relativeName))

		content, err := os.ReadFile(filePath)
		if err == nil {
			file.Write(content)
			// Asegurar que el contenido termina con una nueva línea
			if len(content) > 0 && content[len(content)-1] != '\n' {
				file.WriteString("\n")
			}
		} else if strings.Contains(err.Error(), "invalid UTF-8") {
			file.WriteString("[Binary file or incompatible encoding]\n")
		} else {
			file.WriteString(fmt.Sprintf("[Error reading file: %s]\n", err.Error()))
		}

		file.WriteString(fmt.Sprintf("// End of file %s\n\n", relativeName))
	}

	return outputFile
}

// GenerateCombinedFile genera un archivo combinado a partir de una lista de archivos
func GenerateCombinedFile(fileList []string, baseDir string) string {
	if len(fileList) == 0 {
		return ""
	}

	// Crear un hash único basado en la hora actual y los archivos seleccionados
	fileListCopy := make([]string, len(fileList))
	copy(fileListCopy, fileList)
	sort.Strings(fileListCopy)

	hashInput := fmt.Sprintf("%d%v", time.Now().Unix(), fileListCopy)
	hasher := md5.New()
	hasher.Write([]byte(hashInput))
	hashValue := hex.EncodeToString(hasher.Sum(nil))[:8]

	filesToProcess := []string{}

	// Recopilar archivos a procesar
	for _, path := range fileList {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		if info.IsDir() {
			// Recorrer el directorio recursivamente
			filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
				if err != nil {
					return nil // Continuar con el siguiente archivo
				}

				if !fileInfo.IsDir() {
					filesToProcess = append(filesToProcess, filePath)
				}
				return nil
			})
		} else {
			filesToProcess = append(filesToProcess, path)
		}
	}

	if len(filesToProcess) == 0 {
		return ""
	}

	// Crear el archivo de salida
	outputFile := filepath.Join(baseDir, fmt.Sprintf("cs_%s.txt", hashValue))
	file, err := os.Create(outputFile)
	if err != nil {
		return ""
	}
	defer file.Close()

	// Escribir el contenido de cada archivo
	for _, filePath := range filesToProcess {
		relPath, err := filepath.Rel(baseDir, filePath)
		if err != nil {
			relPath = filePath
		}
		relativeName := filepath.ToSlash(relPath)

		file.WriteString("---------------------------------------------\n")
		file.WriteString(fmt.Sprintf("// File %s\n", relativeName))

		content, err := os.ReadFile(filePath)
		if err == nil {
			file.Write(content)
			// Asegurar que el contenido termina con una nueva línea
			if len(content) > 0 && content[len(content)-1] != '\n' {
				file.WriteString("\n")
			}
		} else if strings.Contains(err.Error(), "invalid UTF-8") {
			file.WriteString("[Binary file or incompatible encoding]\n")
		} else {
			file.WriteString(fmt.Sprintf("[Error reading file: %s]\n", err.Error()))
		}

		file.WriteString(fmt.Sprintf("// End of file %s\n\n", relativeName))
	}

	return outputFile
}
