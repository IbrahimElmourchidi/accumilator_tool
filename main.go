package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

func main() {
	fmt.Println("Accumilator - File Accumulation Tool")
	fmt.Println("=====================================")

	// Get current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Current directory: %s\n", currentDir)

	// Get all directories in the current directory
	dirs, err := ioutil.ReadDir(currentDir)
	if err != nil {
		fmt.Printf("Error reading current directory: %v\n", err)
		os.Exit(1)
	}

	var availableDirs []string
	for _, dir := range dirs {
		if dir.IsDir() {
			availableDirs = append(availableDirs, dir.Name())
		}
	}

	if len(availableDirs) == 0 {
		fmt.Println("No subdirectories found in current directory.")
		return
	}

	// Ask user to select directories using interactive prompt
	selectedDirs := selectDirectoriesInteractive(availableDirs)
	if len(selectedDirs) == 0 {
		fmt.Println("No directories selected. Exiting.")
		return
	}

	// Ask for file extensions
	fmt.Println() // Add a newline before asking for extensions
	extensions := selectExtensions()
	if len(extensions) == 0 {
		fmt.Println("No extensions selected. Exiting.")
		return
	}

	// Process selected directories
	outputFile := filepath.Join(currentDir, "accumulated_files.txt")
	err = processDirectories(selectedDirs, extensions, outputFile)
	if err != nil {
		fmt.Printf("Error processing directories: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nDone! Combined output saved to: %s\n", outputFile)
}

// BellSkipper implements an io.WriteCloser that skips bell characters
type BellSkipper struct{}

func (bs *BellSkipper) Write(b []byte) (int, error) {
	// Skip the bell character (ASCII 7)
	const bellChar = 7
	filtered := make([]byte, 0, len(b))
	for _, c := range b {
		if c != bellChar {
			filtered = append(filtered, c)
		}
	}
	// Write the filtered bytes to os.Stdout
	return os.Stdout.Write(filtered)
}

func (bs *BellSkipper) Close() error {
	return nil
}

func selectDirectoriesInteractive(dirs []string) []string {
	// Include current directory as an option
	allDirs := append([]string{"."}, dirs...)

	// Create a multi-select prompt
	selectedDirs := []string{}

	// Create a label to show current selections
	label := "Select directories to process (Use UP/DOWN to navigate, ENTER to select/deselect, 'f' to finish)"

	// Use a loop to allow multiple selections
	for {
		// Create a copy of the directory list with selection indicators
		displayItems := make([]string, len(allDirs))
		for i, dir := range allDirs {
			if dir == "." {
				dir = "."
			}

			selected := false
			for _, sel := range selectedDirs {
				if sel == dir {
					selected = true
					break
				}
			}
			if selected {
				displayItems[i] = "[x] " + dir
			} else {
				displayItems[i] = "[ ] " + dir
			}
		}

		// Add option to finish selection
		displayItems = append(displayItems, ">>> FINISH SELECTION <<<")

		prompt := promptui.Select{
			Label: label,
			Items: displayItems,
			Stdout: &BellSkipper{}, // Disable sound
			Templates: &promptui.SelectTemplates{
				Active:   `{{ ">" | green }} {{ . }}`,
				Inactive: `  {{ . }}`,
				Selected: `{{ . }}`,
			},
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return nil
		}

		// Check if user wants to finish selection
		if result == ">>> FINISH SELECTION <<<" {
			break
		}

		// Extract the actual directory name (remove the checkbox)
		dirName := strings.TrimPrefix(result, "[x] ")
		dirName = strings.TrimPrefix(dirName, "[ ] ")

		// Convert back to actual directory name
		if dirName == "." {
			dirName = "."
		}

		// Toggle selection
		foundIndex := -1
		for i, sel := range selectedDirs {
			if sel == dirName {
				foundIndex = i
				break
			}
		}

		if foundIndex != -1 {
			// Remove from selection
			selectedDirs = append(selectedDirs[:foundIndex], selectedDirs[foundIndex+1:]...)
		} else {
			// Add to selection
			selectedDirs = append(selectedDirs, dirName)
		}

		// Update label to show current count
		if len(selectedDirs) > 0 {
			label = fmt.Sprintf("Selected %d directories - Use UP/DOWN to navigate, ENTER to select/deselect, 'f' to finish", len(selectedDirs))
		} else {
			label = "Select directories to process (Use UP/DOWN to navigate, ENTER to select/deselect, 'f' to finish)"
		}
	}

	return selectedDirs
}

func selectExtensions() []string {
	// Use a prompt to get file extensions
	prompt := promptui.Prompt{
		Label: "Enter file extensions to accumulate (comma-separated, e.g., ts,dart,json), or press Enter for all files",
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil
	}

	if result == "" {
		// Return empty slice to indicate all files
		return []string{}
	}

	// Split and clean extensions
	extList := strings.Split(result, ",")
	var extensions []string

	for _, ext := range extList {
		ext = strings.TrimSpace(ext)
		if ext != "" {
			// Ensure extension starts with a dot
			if !strings.HasPrefix(ext, ".") {
				ext = "." + ext
			}
			extensions = append(extensions, strings.ToLower(ext))
		}
	}

	// Remove duplicates while preserving order
	seen := make(map[string]bool)
	var resultSlice []string
	for _, ext := range extensions {
		if !seen[ext] {
			seen[ext] = true
			resultSlice = append(resultSlice, ext)
		}
	}

	return resultSlice
}

func processDirectories(dirs []string, extensions []string, outputFile string) error {
	// Create or clear the output file
	outputFileHandle, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFileHandle.Close()

	fileCount := 0

	// Process each selected directory
	for _, dir := range dirs {
		fmt.Printf("\nProcessing directory: %s\n", dir)

		// Walk through the directory tree
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				// Skip inaccessible files/directories
				return nil
			}

			if !info.IsDir() {
				// Check if file has one of the selected extensions (or if extensions is empty - meaning all files)
				shouldProcess := len(extensions) == 0 // If no extensions specified, process all files
				if !shouldProcess {
					ext := strings.ToLower(filepath.Ext(path))
					for _, allowedExt := range extensions {
						if ext == allowedExt {
							shouldProcess = true
							break
						}
					}
				}

				if shouldProcess {
					// Read file content
					content, err := ioutil.ReadFile(path)
					if err != nil {
						fmt.Printf("Error reading file %s: %v\n", path, err)
						return nil
					}

					// Write header and content to output file
					_, err = outputFileHandle.WriteString(fmt.Sprintf("// File: %s\n", path))
					if err != nil {
						return err
					}

					_, err = outputFileHandle.WriteString(string(content))
					if err != nil {
						return err
					}

					_, err = outputFileHandle.WriteString("\n\n//------------------------------------------------------------------------------\n\n")
					if err != nil {
						return err
					}

					fileCount++
					fmt.Printf("Processed: %s\n", path)
				}
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("error walking directory %s: %v", dir, err)
		}
	}

	fmt.Printf("\nTotal processed files: %d\n", fileCount)
	return nil
}