# 🐱 Cat Selector

> The fastest way to explore, select and concatenate files from your terminal.

![Demo](demo.gif)


Cat Selector is an interactive terminal tool for browsing, selecting, concatenating, and exporting text content in files and directories, combining multiple actions into a single seamless flow.

## Why Cat Selector?
Unlike traditional file managers (ranger, lf) or basic commands (find, cat, less), Cat Selector provides:
- True multi-selection across directories and files.
- Instant generation and single concatenated text output.

A direct, visual, and straight-forward content selector and concatenated text file exporter.

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


## Installation

### Requirements

- [Go](https://go.dev/dl/) (>=1.20) — needed for manual builds
- [Nerd Fonts](https://www.nerdfonts.com/) — required for correct icon rendering

### macOS (via Homebrew)

```bash
brew tap alexaldearroyo/catselector
brew install catsel
```

### Manual Installation (macOS / Linux)

First, clone the repository:

```bash
git clone https://github.com/alexaldearroyo/catselector.git
cd catselector
```

Then build and install:

```bash
make build
sudo make install
```
This will compile the project and install the catsel binary into `/usr/local/bin`.

## Usage

```bash
catsel            # Start Cat Selector
catsel --help     # Show help message
catsel --version  # Show version information
```

## Contributing

Contributions are welcome. Please open an issue to discuss major changes before submitting a pull request.

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.


© Alex Arroyo 2025
