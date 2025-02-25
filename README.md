# lazyjson

Interactive TUI JSON viewer with tree navigation, expand/collapse functionality, and real-time filtering.

![LazyJSON Demo](screenshots/demo.gif)

## Features

- 🌳 Tree-based JSON navigation
- ⌨️ Vim-style keyboard shortcuts
- 🔍 Real-time JSON filtering
- 📁 Support for large JSON files
- 🎨 Syntax highlighting
- 🔄 Expand/Collapse nodes

## Installation

```bash
go install github.com/cksidharthan/lazyjson@latest
```

Or build from source:

```bash
git clone https://github.com/cksidharthan/lazyjson.git
cd lazyjson
go build
```

## Usage

### Basic Usage

```bash
# View a JSON file
lazyjson data.json

# View JSON from a URL
curl https://api.example.com/data.json > data.json && lazyjson data.json
```

![Basic Usage](screenshots/basic-usage.png)

### Navigation

- `↑/k`: Move cursor up
- `↓/j`: Move cursor down
- `Enter/Space`: Expand/Collapse node
- `h`: Collapse node
- `l`: Expand node
- `q`: Quit

![Navigation Demo](screenshots/navigation.png)

### Filtering

Type to filter JSON nodes in real-time. The filter supports:
- Exact matches
- Partial matches
- Case-insensitive search

![Filtering Demo](screenshots/filtering.png)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) for details
