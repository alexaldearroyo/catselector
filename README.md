# 🐱 Cat Selector

Cat Selector is an interactive terminal tool for browsing, selecting, concatenating, and exporting files and directories. It allows users to explore directories, select multiple files (or entire folders), and combine their text content for viewing, external editing, or clipboard copying.

## Key Features

- **Split Navigation**: Divided panels for directory, file navigation and file and subdirectories preview.
- **Multiple Selection**: Quick selection of multiple files and subdirectories
- **Concatenation**: Combine selected content into a single output.
- **Flexible Export**: 
  - Export to temporary file
  - Direct clipboard copying
- **Intuitive Navigation**: Keyboard keybindings optimized for productivity

## Keybindings

| Key | Action |
|-------|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `Enter` / `l` | Enter directory |
| `Esc` / `h` | Go to previous directory |
| `s` | Select/Deselect |
| `a` | Select/Deselect all |
| `i` | Toggle include subdirectories |
| `o` | Concatenate and open in external editor |
| `c` | Concatenate and copy to clipboard |
| `Tab` | Switch panel |
| `f` | Go to files panel |
| `d` | Go to directories panel |
| `q` | Quit |

## Technologies

Built in Go using:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for terminal interface
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) for visual styling

## Technical Features

- Efficient handling of plain text files
- Support for complex folder structures
- Optional subdirectory inclusion
- Cross-platform compatibility
- Intuitive and responsive user interface

## 🚀 Installation

_TO DO_

## Usage

```bash
catsel            # Start Cat Selector
catsel --help     # Show help message
catsel --version  # Show version information
```

## Contributing

Contributions are welcome. Please open an issue to discuss major changes before submitting a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


© Alex Arroyo 2025
