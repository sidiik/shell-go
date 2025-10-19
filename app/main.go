package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		args := parseUserInput(strings.Split(input, "\n")[0])

		if len(args) < 1 {
			fmt.Println("Format should be like <command args>")
			continue
		}

		switch args[0] {
		case "exit":
			var exitCode int
			if len(args) < 2 {
				exitCode = 0
			} else {
				exitCode, err = strconv.Atoi(args[1])
				if err != nil {
					fmt.Println("The exit code should be 1 or 0")
					continue
				}
			}

			os.Exit(exitCode)

		case "echo":
			// GOOD CODE
			fmt.Println(strings.Join(args[1:], " "))

		case "type":
			if len(args) < 2 {
				fmt.Println("Invalid argument, Please use this format type <command>")
				continue
			}

			arg := args[1]

			if isBuiltinCommand(arg) {
				fmt.Printf("%s is a shell builtin\n", arg)
			} else {
				path, err := findExecutablePath(arg)
				if err != nil {
					fmt.Printf("%s: not found\n", arg)
					continue
				}

				fmt.Printf("%s is %s\n", arg, path)
			}

		case "pwd":
			printWorkingDirectory()
		case "cd":
			if !isArgsValid(args, 2) {
				fmt.Printf("Usage: cd <dir>\n")
				continue
			}

			err := changeWorkingDirectory(args[1])

			if err != nil {
				fmt.Printf("cd: %s: No such file or directory\n", args[1])
				continue
			}

			// printWorkingDirectory()

		default:
			_, err := findExecutablePath(args[0])

			if err != nil {
				fmt.Printf("%s: command not found\n", args[0])
			}

			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			continue
		}

	}

}

func findExecutablePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty path passed")
	}

	pathStr, err := exec.LookPath(path)
	if err != nil {
		return "", fmt.Errorf("file not found. %w", err)
	}

	return pathStr, nil
}

func isBuiltinCommand(cmd string) bool {
	builtIns := []string{"exit", "echo", "type", "pwd"}

	for _, builtinCommand := range builtIns {
		if strings.EqualFold(cmd, builtinCommand) {
			return true
		}
	}

	return false
}

func isArgsValid(args []string, requiredArgsCount int) bool {
	return len(args) == requiredArgsCount
}

func printWorkingDirectory() {
	wd, _ := os.Getwd()
	fmt.Println(wd)
}

func changeWorkingDirectory(dest string) error {
	if dest == "~" {
		home := os.Getenv("HOME")
		return os.Chdir(home)
	}
	return os.Chdir(dest)
}

func parseUserInput(s string) []string {
	var result []string
	var currentToken string
	var isInQoutes, isInDoubleQoutes, escaped bool

	for _, str := range s {

		if escaped {
			currentToken += string(str)
			escaped = false
			continue
		}

		switch str {
		case '"':
			if !isInQoutes {
				isInDoubleQoutes = !isInDoubleQoutes
				continue
			}

			currentToken += string(str)
			continue

		case '\'':
			if !isInDoubleQoutes {
				isInQoutes = !isInQoutes
				continue
			} else {
				currentToken += string(str)
				continue
			}

		case '\\':
			if isInDoubleQoutes || isInQoutes {
				currentToken += string(str)
				continue
			} else {
				escaped = true
				continue
			}

		case ' ':
			if isInQoutes || isInDoubleQoutes {
				currentToken += string(str)
				continue
			} else {
				if currentToken != "" {
					result = append(result, currentToken)
					currentToken = ""
					continue
				}
			}

		default:
			currentToken += string(str)
		}
	}

	if currentToken != "" {
		result = append(result, currentToken)
	}

	return result

}

// 1. bus