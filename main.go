// File: main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

	// Get all directories in the current directory (excluding hidden ones)
	entries, err := os.ReadDir(currentDir)
	if err != nil {
		fmt.Printf("Error reading current directory: %v\n", err)
		os.Exit(1)
	}

	var availableDirs []string
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			availableDirs = append(availableDirs, entry.Name())
		}
	}

	// ALWAYS include current directory as an option (even if no subdirs exist)
	// Proceed to selection regardless of subdir count
	selectedDirs := selectDirectoriesInteractive(availableDirs)
	if len(selectedDirs) == 0 {
		fmt.Println("No directories selected. Exiting.")
		return
	}

	// Ask for file extensions with wildcard support
	fmt.Println()
	extensions := selectExtensions()
	if extensions == nil {
		fmt.Println("Selection cancelled. Exiting.")
		return
	}

	// Ask for max file size limit (default 10 KB)
	fmt.Println()
	maxSizeKB := selectMaxFileSize()

	// Process selected directories
	outputFile := filepath.Join(currentDir, "accumulated_files.txt")
	err = processDirectories(selectedDirs, extensions, maxSizeKB, outputFile)
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

func selectDirectoriesInteractive(subdirs []string) []string {
	// ALWAYS include current directory as first option
	allDirs := []string{"."}
	allDirs = append(allDirs, subdirs...)

	selectedDirs := []string{}
	label := "Select directories to process (↑↓ to navigate, Enter to toggle selection)"

	for {
		// Create display items with checkboxes
		displayItems := make([]string, len(allDirs))
		for i, dir := range allDirs {
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

		// Add finish option
		displayItems = append(displayItems, ">>> FINISH SELECTION <<<")

		prompt := promptui.Select{
			Label:  label,
			Items:  displayItems,
			Stdout: &BellSkipper{},
			Templates: &promptui.SelectTemplates{
				Active:   `{{ ">" | green }} {{ . }}`,
				Inactive: `  {{ . }}`,
				Selected: `{{ . }}`,
			},
		}

		_, result, err := prompt.Run()
		if err != nil {
			if err == promptui.ErrInterrupt {
				fmt.Println("\nSelection cancelled by user.")
				os.Exit(0)
			}
			fmt.Printf("Prompt failed: %v\n", err)
			return nil
		}

		// Finish selection
		if result == ">>> FINISH SELECTION <<<" {
			if len(selectedDirs) == 0 {
				fmt.Println("⚠ No directories selected. Continuing anyway...")
			}
			break
		}

		// Extract actual directory name
		dirName := strings.TrimPrefix(strings.TrimPrefix(result, "[x] "), "[ ] ")

		// Toggle selection
		foundIndex := -1
		for i, sel := range selectedDirs {
			if sel == dirName {
				foundIndex = i
				break
			}
		}

		if foundIndex != -1 {
			selectedDirs = append(selectedDirs[:foundIndex], selectedDirs[foundIndex+1:]...)
		} else {
			selectedDirs = append(selectedDirs, dirName)
		}

		// Update label
		if len(selectedDirs) > 0 {
			label = fmt.Sprintf("✓ %d directories selected (↑↓ to navigate, Enter to toggle, choose 'FINISH' when done)", len(selectedDirs))
		} else {
			label = "Select directories to process (↑↓ to navigate, Enter to toggle selection)"
		}
	}

	return selectedDirs
}

func selectExtensions() []string {
	prompt := promptui.Prompt{
		Label: "Enter file extensions (comma-separated, e.g., ts,dart,json) or '*' for all files",
		Validate: func(input string) error {
			input = strings.TrimSpace(input)
			if input == "" {
				return fmt.Errorf("input cannot be empty - enter '*' for all files or specify extensions")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			return nil
		}
		fmt.Printf("Prompt failed: %v\n", err)
		return nil
	}

	result = strings.TrimSpace(result)

	// Handle wildcard
	if result == "*" {
		fmt.Println("✓ Wildcard '*' selected: ALL file types will be included")
		return []string{}
	}

	// Process extensions
	extList := strings.Split(result, ",")
	var extensions []string

	for _, ext := range extList {
		ext = strings.TrimSpace(ext)
		if ext != "" {
			if !strings.HasPrefix(ext, ".") {
				ext = "." + ext
			}
			extensions = append(extensions, strings.ToLower(ext))
		}
	}

	// Remove duplicates
	seen := make(map[string]bool)
	var resultSlice []string
	for _, ext := range extensions {
		if !seen[ext] {
			seen[ext] = true
			resultSlice = append(resultSlice, ext)
		}
	}

	if len(resultSlice) > 0 {
		fmt.Printf("✓ Selected extensions: %s\n", strings.Join(resultSlice, ", "))
	} else {
		fmt.Println("⚠ No valid extensions specified - will process all files")
	}
	return resultSlice
}

func selectMaxFileSize() int64 {
	prompt := promptui.Prompt{
		Label:   "Maximum file size to include (in KB, 0 = no limit)",
		Default: "10", // CHANGED DEFAULT TO 10 KB (safer default)
		Validate: func(input string) error {
			size, err := strconv.ParseInt(strings.TrimSpace(input), 10, 64)
			if err != nil {
				return fmt.Errorf("please enter a valid number")
			}
			if size < 0 {
				return fmt.Errorf("size cannot be negative")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			os.Exit(0)
		}
		fmt.Printf("Prompt failed: %v\n", err)
		return 10 // Safe default on error
	}

	sizeKB, _ := strconv.ParseInt(strings.TrimSpace(result), 10, 64)
	if sizeKB == 0 {
		fmt.Println("✓ No size limit (all files will be included)")
	} else {
		fmt.Printf("✓ Files larger than %d KB will be skipped\n", sizeKB)
	}
	return sizeKB
}

func processDirectories(dirs []string, extensions []string, maxSizeKB int64, outputFile string) error {
	// Resolve absolute path of output file for accurate comparison
	absOutputFile, err := filepath.Abs(outputFile)
	if err != nil {
		return fmt.Errorf("error resolving output file path: %v", err)
	}

	outputFileHandle, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFileHandle.Close()

	fileCount := 0
	skippedSize := 0
	skippedExt := 0
	skippedOutput := 0 // Track skipped output files
	totalSize := int64(0)

	for _, dir := range dirs {
		fmt.Printf("\n📁 Processing directory: %s\n", dir)

		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("  ⚠ Skipped (inaccessible): %s\n", path)
				return nil
			}

			if info.IsDir() {
				// Skip node_modules, .git, and other common large directories
				dirName := info.Name()
				if dirName == "node_modules" || dirName == ".git" || dirName == "vendor" || dirName == "__pycache__" {
					return filepath.SkipDir
				}
				return nil
			}

			// EXPLICITLY EXCLUDE accumulated_files.txt (case-insensitive)
			baseName := filepath.Base(path)
			if strings.EqualFold(baseName, "accumulated_files.txt") {
				skippedOutput++
				return nil
			}

			// Skip hidden files (starting with .) except .gitignore, .env, etc.
			if strings.HasPrefix(baseName, ".") && 
			   baseName != ".gitignore" && 
			   baseName != ".env" && 
			   baseName != ".env.example" {
				return nil
			}

			// Skip the output file itself (in case it exists before creation)
			absPath, _ := filepath.Abs(path)
			if absPath == absOutputFile {
				skippedOutput++
				return nil
			}

			// Size check
			if maxSizeKB > 0 {
				fileSizeKB := info.Size() / 1024
				if fileSizeKB > maxSizeKB {
					fmt.Printf("  ⚠ Skipped (size: %d KB > limit %d KB): %s\n", fileSizeKB, maxSizeKB, path)
					skippedSize++
					return nil
				}
			}

			// Extension check
			shouldProcess := len(extensions) == 0
			if !shouldProcess {
				ext := strings.ToLower(filepath.Ext(path))
				for _, allowedExt := range extensions {
					if ext == allowedExt {
						shouldProcess = true
						break
					}
				}
			}

			if !shouldProcess {
				skippedExt++
				return nil
			}

			// Read and write file
			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("  ⚠ Error reading file %s: %v\n", path, err)
				return nil
			}

			header := fmt.Sprintf("// File: %s (%d bytes)\n", path, info.Size())
			if _, err := outputFileHandle.WriteString(header); err != nil {
				return err
			}
			if _, err := outputFileHandle.Write(content); err != nil {
				return err
			}
			if _, err := outputFileHandle.WriteString("\n\n//------------------------------------------------------------------------------\n\n"); err != nil {
				return err
			}

			fileCount++
			totalSize += info.Size()
			fmt.Printf("  ✓ %s (%d KB)\n", path, info.Size()/1024)
			return nil
		})

		if err != nil {
			return fmt.Errorf("error walking directory %s: %v", dir, err)
		}
	}

	// Final summary
	fmt.Printf("\n" + strings.Repeat("=", 55))
	fmt.Println("\n✅ ACCUMULATION COMPLETE")
	fmt.Printf("   Processed files:    %d\n", fileCount)
	fmt.Printf("   Skipped (size):     %d\n", skippedSize)
	fmt.Printf("   Skipped (ext):      %d\n", skippedExt)
	fmt.Printf("   Skipped (output):   %d\n", skippedOutput) // Show excluded accumulation files
	fmt.Printf("   Total size:         %.2f MB\n", float64(totalSize)/1024/1024)
	fmt.Println(strings.Repeat("=", 55))

	return nil
}