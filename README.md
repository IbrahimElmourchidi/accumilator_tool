# Accumilator

Accumilator is a powerful command-line tool that helps you accumulate files with specific extensions from selected directories into a single text file. It's designed to make it easy to consolidate code files for review, analysis, or sharing.

## Features

- Interactive directory selection from the current working directory
- Support for selecting multiple directories using an intuitive menu (with visual checkboxes)
- Option to select all directories with the "all" keyword
- Flexible file extension filtering (e.g., ts, dart, json, js, py, etc.)
- **NEW**: Wildcard support (`*`) to include ALL file types without filtering
- **NEW**: File size filtering to exclude large files and avoid bloated output
- Clean output format with file headers and separators
- Cross-platform compatibility (Windows, Linux, macOS)
- Modern I/O operations using os.ReadFile/os.ReadDir (replacing deprecated ioutil)
- Proper error handling with graceful skipping of inaccessible files
- Intuitive multi-select interface with visual checkboxes ([x] / [ ])
- Smart defaults and clear visual feedback (📁 processing, ✓ success, ⚠ skipped)
- Ctrl+C support for graceful interruption of the process
- Detailed processing feedback with visual indicators
- Comprehensive summary report after processing

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
4. Ask for maximum file size limit (optional)
5. Combine all matching files into a single `accumulated_files.txt` file

### Example Usage

```bash
# Run in current directory
accumilator

# The tool will then guide you through the process:
# 1. Select directories using an interactive menu (with current directory option)
# 2. Specify file extensions (comma-separated) or use '*' for all files
# 3. Set maximum file size limit in KB (0 = no limit)
# 4. Wait for processing to complete
```

### Interactive Prompts

- **Directory Selection**: Use UP/DOWN arrows to navigate, ENTER to select/deselect directories from an interactive list (including the current directory option), then select "FINISH SELECTION" to confirm. Visual checkboxes ([x]/[ ]) show your selections.
- **Extension Selection**: Enter file extensions separated by commas (e.g., "ts,dart,json"), or use `*` to include all files regardless of extension
- **Size Limit**: Enter maximum file size in KB (0 = no limit). Files exceeding this limit will be skipped with a clear notification.

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
// File: src/main.ts (2456 bytes)
console.log("Hello, world!");

// File: src/utils.ts (1203 bytes)
export function helper() { ... }

// File: tests/main.test.ts (3567 bytes)
describe("Test suite", () => { ... });

//------------------------------------------------------------------------------
```

### Processing Feedback

During execution, you'll see visual feedback like:

```
📁 Processing directory: src
  ✓ src/main.ts (2 KB)
  ✓ src/utils.ts (1 KB)
  
📁 Processing directory: tests
  ✓ tests/main.test.ts (3 KB)
  ⚠ Skipped (size: 150 KB > limit 100 KB): large_file.bin
  
=======================================================
✅ ACCUMULATION COMPLETE
   Processed files:  3
   Skipped (size):   1
   Skipped (ext):    5
   Total size:       0.01 MB
=======================================================
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