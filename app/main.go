package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		// Display prompt
		fmt.Print("$ ")

		// Read input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			continue
		}

		// Remove newline character and trim spaces
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Split input into command and arguments
		parts := strings.Fields(input)
		command := parts[0]
		args := parts[1:]

		// Handle built-in commands
		switch command {
		case "exit":
			os.Exit(0)
		case "type":
			handleTypeCommand(args)
			continue
		}

		// Execute external program
		if err := executeExternalProgram(command, args); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", command, err)
		}
	}
}

func handleTypeCommand(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "type: missing argument\n")
		return
	}

	command := args[0]

	// Check if it's a built-in command
	switch command {
	case "exit", "type":
		fmt.Printf("%s is a shell builtin\n", command)
		return
	}

	// Search in PATH
	path := os.Getenv("PATH")
	dirs := strings.Split(path, string(os.PathListSeparator))

	for _, dir := range dirs {
		fullPath := filepath.Join(dir, command)
		if isExecutable(fullPath) {
			fmt.Printf("%s is %s\n", command, fullPath)
			return
		}
	}

	fmt.Printf("%s: not found\n", command)
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// Check if it's a regular file and executable
	if !info.Mode().IsRegular() {
		return false
	}

	// Check execute permission
	mode := info.Mode()
	return mode&0111 != 0 // Check if any execute bit is set
}

func executeExternalProgram(command string, args []string) error {
	// Search for the executable in PATH
	path := os.Getenv("PATH")
	dirs := strings.Split(path, string(os.PathListSeparator))

	var executablePath string
	for _, dir := range dirs {
		fullPath := filepath.Join(dir, command)
		if isExecutable(fullPath) {
			executablePath = fullPath
			break
		}
	}

	if executablePath == "" {
		return fmt.Errorf("command not found")
	}

	// Prepare the command with arguments
	cmd := exec.Command(executablePath, args...)

	// Set up standard I/O
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	return cmd.Run()
}
