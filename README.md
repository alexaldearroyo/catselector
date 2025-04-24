# CopyCat

Una herramienta para seleccionar, previsualizar y exportar archivos y directorios en formato de texto.

## Descripción

CopyCat es una herramienta de línea de comandos que permite:

- Navegar por directorios y archivos con una interfaz de tres paneles
- Seleccionar archivos y directorios individuales o en masa
- Ver previsualizaciones de archivos de texto
- Exportar los archivos seleccionados a un único archivo de texto
- Copiar el contenido al portapapeles

Esta es una versión en Go de la herramienta original escrita en Python.

## Instalación

```bash
go install github.com/alexadler/copycat@latest
```

O clonar el repositorio y compilar:

```bash
git clone https://github.com/alexadler/copycat.git
cd copycat
go build
```

## Uso

```bash
copycat [--dir DIRECTORIO_INICIAL]
```

Si no se especifica un directorio, se usará el directorio actual.

## Atajos de teclado

- **k/j**: Navegar arriba/abajo
- **h/l**: Navegar atrás/adelante en directorios
- **d**: Cambiar a panel de directorios
- **f**: Cambiar a panel de archivos
- **p**: Cambiar a panel de vista previa
- **Tab**: Alternar entre los paneles
- **s**: Seleccionar elemento actual
- **a**: Seleccionar/deseleccionar todos los elementos
- **i**: Alternar inclusión de subdirectorios
- **o**: Exportar archivos seleccionados y abrir
- **c**: Exportar y copiar al portapapeles
- **Esc/h**: Volver al directorio padre
- **q**: Salir
- **/**: Buscar

## Licencia

MIT 
