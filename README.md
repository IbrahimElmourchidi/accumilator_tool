# Accumilator

Accumilator is a powerful command-line tool that helps you accumulate files with specific extensions from selected directories into a single text file. It's designed to make it easy to consolidate code files for review, analysis, or sharing.

## Features

- Interactive directory selection from the current working directory
- Support for selecting multiple directories using comma-separated values or ranges (e.g., "1,3,5" or "2-4")
- Option to select all directories with the "all" keyword
- Flexible file extension filtering (e.g., ts, dart, json, js, py, etc.)
- Option to include all files regardless of extension
- Clean output format with file headers and separators
- Cross-platform compatibility (Windows, Linux, macOS)

## Installation

### Pre-built Binaries

Download the appropriate binary for your platform from the [releases page](https://github.com/yourusername/accumilator/releases):

- **Windows**: `accumilator-windows-amd64.exe`
- **Linux**: `accumilator-linux-amd64`
- **macOS**: `accumilator-darwin-amd64`

For ARM64 systems:
- **Linux ARM64**: `accumilator-linux-arm64`
- **macOS ARM64**: `accumilator-darwin-arm64`

After downloading, rename the binary to `accumilator` (or `accumilator.exe` on Windows) and add it to your system's PATH.

### Building from Source

#### Prerequisites
- Install Go (version 1.19 or later)

#### Build Instructions by Platform

**Build for Current Platform:**
```bash
go build -o accumilator main.go
```

**Build for Windows (from any platform):**
```bash
GOOS=windows GOARCH=amd64 go build -o accumilator-windows-amd64.exe main.go
```

**Build for Linux:**
```bash
GOOS=linux GOARCH=amd64 go build -o accumilator-linux-amd64 main.go
```

**Build for macOS:**
```bash
GOOS=darwin GOARCH=amd64 go build -o accumilator-darwin-amd64 main.go
```

**Build for Linux ARM64:**
```bash
GOOS=linux GOARCH=arm64 go build -o accumilator-linux-arm64 main.go
```

**Build for macOS ARM64:**
```bash
GOOS=darwin GOARCH=arm64 go build -o accumilator-darwin-arm64 main.go
```

**Build all platforms at once (using build script):**
```bash
chmod +x build.sh  # On Unix-like systems
./build.sh
```

#### Cross-compilation Notes
- Set the `GOOS` environment variable to your target OS (`windows`, `linux`, `darwin`)
- Set the `GOARCH` environment variable to your target architecture (`amd64`, `arm64`)
- The output binary will match the target platform regardless of your build machine

## Usage

Simply run the `accumilator` command from any directory. The tool will:

1. Show all subdirectories in the current directory
2. Prompt you to select which directories to process
3. Ask for the file extensions you want to accumulate
4. Combine all matching files into a single `accumulated_files.txt` file

### Example Usage

```bash
# Run in current directory
accumilator

# The tool will then guide you through the process:
# 1. Select directories using an interactive menu (with current directory option)
# 2. Specify file extensions (comma-separated)
# 3. Wait for processing to complete
```

### Interactive Prompts

- **Directory Selection**: Use UP/DOWN arrows to navigate, ENTER to select/deselect directories from an interactive list (including the current directory option), then select "FINISH SELECTION" to confirm
- **Extension Selection**: Enter file extensions separated by commas (e.g., "ts,dart,json"), or press Enter to include all files

## Output Format

The output file (`accumulated_files.txt`) contains:

- A comment header with the file path: `// File: path/to/file.ext`
- The full content of each file
- A separator line between files: `//------------------------------------------------------------------------------`

## Example

If you have a project structure like:

```
my-project/
├── src/
│   ├── main.ts
│   └── utils.ts
├── tests/
│   └── main.test.ts
└── docs/
    └── README.md
```

Running `accumilator` and selecting `src` and `tests` with extension `.ts` will create:

```text
// File: src/main.ts
console.log("Hello, world!");

// File: src/utils.ts
export function helper() { ... }

// File: tests/main.test.ts
describe("Test suite", () => { ... });

//------------------------------------------------------------------------------
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have suggestions for improvements, please open an issue on the GitHub repository.