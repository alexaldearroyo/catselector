package export

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// GenerateTextFile genera un archivo de texto con el contenido de los archivos seleccionados
func GenerateTextFile(selected []string, excluded []string, includeSubdirs bool, initialDir, currentDir string) (string, error) {
	// Generar hash para el nombre del archivo
	now := time.Now().UnixNano()
	sort.Strings(selected)
	hashInput := fmt.Sprintf("%d%v", now, selected)
	hash := md5.Sum([]byte(hashInput))
	hashStr := fmt.Sprintf("%x", hash)[:8]

	var filesToProcess []string

	// Recopilar archivos a procesar
	for _, path := range selected {
		if !contains(excluded, path) {
			fileInfo, err := os.Stat(path)
			if err != nil {
				continue
			}

			if fileInfo.IsDir() {
				if includeSubdirs {
					// Recorrer recursivamente el directorio
					err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
						if err != nil {
							return err
						}
						if !info.IsDir() && !contains(excluded, filePath) {
							filesToProcess = append(filesToProcess, filePath)
						}
						return nil
					})
					if err != nil {
						continue
					}
				} else {
					// Solo incluir archivos inmediatos
					files, err := os.ReadDir(path)
					if err != nil {
						continue
					}
					for _, file := range files {
						if !file.IsDir() {
							filePath := filepath.Join(path, file.Name())
							if !contains(excluded, filePath) {
								filesToProcess = append(filesToProcess, filePath)
							}
						}
					}
				}
			} else {
				filesToProcess = append(filesToProcess, path)
			}
		}
	}

	if len(filesToProcess) == 0 {
		return "", fmt.Errorf("no hay archivos para procesar")
	}

	// Crear archivo de salida
	outputFile := filepath.Join(initialDir, fmt.Sprintf("cs_%s.txt", hashStr))
	file, err := os.Create(outputFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Escribir contenido
	for _, filePath := range filesToProcess {
		relativeName, err := filepath.Rel(currentDir, filePath)
		if err != nil {
			relativeName = filePath
		}

		// Escribir separador y nombre del archivo
		if _, err := file.WriteString(fmt.Sprintf("---------------------------------------------\n// File %s\n", relativeName)); err != nil {
			return "", err
		}

		// Leer y escribir el contenido del archivo
		content, err := os.ReadFile(filePath)
		if err != nil {
			if _, err := file.WriteString(fmt.Sprintf("[Error de lectura: %v]\n", err)); err != nil {
				return "", err
			}
		} else {
			if _, err := file.Write(content); err != nil {
				return "", err
			}
			// Asegurar que haya un salto de línea al final
			if len(content) > 0 && content[len(content)-1] != '\n' {
				if _, err := file.WriteString("\n"); err != nil {
					return "", err
				}
			}
		}

		// Escribir marca de fin de archivo
		if _, err := file.WriteString(fmt.Sprintf("// End of file %s\n\n", relativeName)); err != nil {
			return "", err
		}
	}

	return outputFile, nil
}

// GenerateCombinedFile genera un archivo combinado con todos los archivos listados
func GenerateCombinedFile(fileList []string, baseDir string) (string, error) {
	if len(fileList) == 0 {
		return "", fmt.Errorf("no hay archivos para procesar")
	}

	// Generar hash para el nombre del archivo
	now := time.Now().UnixNano()
	sort.Strings(fileList)
	hashInput := fmt.Sprintf("%d%v", now, fileList)
	hash := md5.Sum([]byte(hashInput))
	hashStr := fmt.Sprintf("%x", hash)[:8]

	var filesToProcess []string

	// Recopilar archivos a procesar
	for _, path := range fileList {
		fileInfo, err := os.Stat(path)
		if err != nil {
			continue
		}

		if fileInfo.IsDir() {
			// Recorrer recursivamente el directorio
			err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					filesToProcess = append(filesToProcess, filePath)
				}
				return nil
			})
			if err != nil {
				continue
			}
		} else {
			filesToProcess = append(filesToProcess, path)
		}
	}

	if len(filesToProcess) == 0 {
		return "", fmt.Errorf("no hay archivos para procesar")
	}

	// Crear archivo de salida
	outputFile := filepath.Join(baseDir, fmt.Sprintf("cs_%s.txt", hashStr))
	file, err := os.Create(outputFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Escribir contenido
	for _, filePath := range filesToProcess {
		relativeName, err := filepath.Rel(baseDir, filePath)
		if err != nil {
			relativeName = filePath
		}

		// Escribir separador y nombre del archivo
		if _, err := file.WriteString(fmt.Sprintf("---------------------------------------------\n// File %s\n", relativeName)); err != nil {
			return "", err
		}

		// Leer y escribir el contenido del archivo
		content, err := os.ReadFile(filePath)
		if err != nil {
			if _, err := file.WriteString(fmt.Sprintf("[Error de lectura: %v]\n", err)); err != nil {
				return "", err
			}
		} else {
			if _, err := file.Write(content); err != nil {
				return "", err
			}
			// Asegurar que haya un salto de línea al final
			if len(content) > 0 && content[len(content)-1] != '\n' {
				if _, err := file.WriteString("\n"); err != nil {
					return "", err
				}
			}
		}

		// Escribir marca de fin de archivo
		if _, err := file.WriteString(fmt.Sprintf("// End of file %s\n\n", relativeName)); err != nil {
			return "", err
		}
	}

	return outputFile, nil
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
