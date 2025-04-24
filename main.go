package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alexadler/copycat/core"
)

func main() {
	// Procesar directorio inicial
	var startDir string
	flag.StringVar(&startDir, "dir", ".", "Directorio inicial")
	flag.Parse()

	// Obtener ruta absoluta
	absPath, err := filepath.Abs(startDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error al resolver la ruta: %v\n", err)
		os.Exit(1)
	}

	// Crear y ejecutar el selector
	selector := core.NewSelector(absPath)
	if err := selector.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
