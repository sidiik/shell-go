package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

const debug = false
const dryRun = false

var builtIns = []string{"exit", "echo", "type", "pwd"}

func main() {

	completer := readline.NewPrefixCompleter(
		readline.PcItem("exit"),
		readline.PcItem("echo"),
		readline.PcItem("cd"),
		readline.PcItem("type"),
		readline.PcItem("pwd"),
	)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "$ ",
		AutoComplete:    completer,
		HistoryFile:     "./history",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})

	if err != nil {
		fmt.Println("Error creating readline:", err)
		return
	}

	defer rl.Close()

	for {
		line, err := rl.Readline()

		if err == readline.ErrInterrupt {
			// Ctrl+C
			if len(line) == 0 {
				break
			}
			continue
		} else if err == io.EOF {
			// Ctrl+D (exit)
			break
		}

		args := parseUserInput(strings.TrimSpace(line))

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

		default:
			executeExternalCommand(args)
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
	var isInSingleQoutes, isInDoubleQoutes, escaped, escapeForRedirector bool
	var skipCount int

	for idx, str := range s {

		if skipCount > 0 {
			skipCount--
			continue
		}

		if escapeForRedirector {
			escapeForRedirector = false
			continue
		}

		if escaped {
			currentToken += string(str)
			escaped = false
			continue
		}

		switch str {
		case '"':
			if !isInSingleQoutes {
				isInDoubleQoutes = !isInDoubleQoutes
				continue
			}

			currentToken += string(str)
			continue

		case '\'':
			if !isInDoubleQoutes {
				isInSingleQoutes = !isInSingleQoutes
				continue
			} else {
				currentToken += string(str)
				continue
			}

		case '\\':
			if isInDoubleQoutes {
				nextRune := rune(s[idx+1])
				if nextRune == '"' || nextRune == '\\' {
					escaped = true
					continue
				}
				currentToken += string(str)
				continue
			} else {
				if isInSingleQoutes {
					currentToken += string(str)
					continue
				}
				escaped = true
				continue
			}

		case ' ':
			if isInSingleQoutes || isInDoubleQoutes {
				currentToken += string(str)
				continue
			} else {
				if currentToken != "" {
					result = append(result, currentToken)
					currentToken = ""
					continue
				}
			}

		case '>':
			if isInDoubleQoutes || isInSingleQoutes {
				currentToken += string(str)
				continue
			} else {
				if currentToken != "" {
					result = append(result, currentToken)
					currentToken = ""
					result = append(result, ">")
					continue
				} else {
					if s[idx+1] == '>' {
						result = append(result, ">>")
						escapeForRedirector = true
						continue
					}
					result = append(result, ">")
					continue
				}

			}
		case '1', '2':
			if len(s) == idx+1 {
				currentToken += string(str)
				continue
			}

			if s[idx+1] == '>' && s[idx+2] != '>' {
				if isInDoubleQoutes || isInSingleQoutes {
					currentToken += string(str)
					continue
				}

				escapeForRedirector = true
				currentToken += string(str) + ">"
				continue
			}

			if s[idx+1] == '>' && s[idx+2] == '>' {
				if isInDoubleQoutes || isInSingleQoutes {
					currentToken += string(str)
					continue
				}

				escapeForRedirector = true
				currentToken += string(str) + ">>"
				skipCount = 1
				continue
			}

			currentToken += string(str)
			continue

		default:
			currentToken += string(str)
		}
	}

	if currentToken != "" {
		result = append(result, currentToken)
	}

	if debug {
		fmt.Printf("TOKENS: %+#v\n", result)
	}

	return result

}

func checkForRedirector(tokens []string) (redirectorIdx int, appendMode bool) {
	for idx, token := range tokens {
		switch token {
		case ">>", "1>>", "2>>":
			redirectorIdx = idx
			appendMode = true
			return
		case ">", "1>", "2>":
			redirectorIdx = idx
			appendMode = false
			return
		}
	}

	return

}

func executeExternalCommand(args []string) {

	_, err := findExecutablePath(args[0])

	redirectorIdx, appendMode := checkForRedirector(args)
	if redirectorIdx != 0 {
		commandTokens := []string{}
		for idx, arg := range args {
			if idx == redirectorIdx {
				break
			}

			commandTokens = append(commandTokens, arg)
			continue

		}

		redirectorType := args[redirectorIdx]
		fileName := args[redirectorIdx+1]
		if err := os.MkdirAll(filepath.Dir(fileName), 0755); err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			return
		}

		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			file, err := os.Create(fileName)
			if err != nil {
				fmt.Printf("Unable to create file: %+#v\n", err)
			}
			defer file.Close()

		}

		fileInfo, _ := os.Stat(fileName)

		if fileInfo.IsDir() {
			fmt.Println("Can not write to dir")
			return
		}

		var f *os.File

		if appendMode {
			f, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				fmt.Println("Unable to open file")
				return
			}
		} else {
			f, err = os.OpenFile(fileName, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				fmt.Println("Unable to open file")
				return
			}
		}

		defer f.Close()

		cmd := exec.Command(commandTokens[0], commandTokens[1:]...)

		switch redirectorType {
		case "2>", "2>>":
			cmd.Stderr = f
			cmd.Stdout = os.Stdout
		default:
			cmd.Stderr = os.Stderr
			cmd.Stdout = f
		}

		run(cmd)

		return
	}

	if err != nil {
		fmt.Printf("%s: command not found\n", args[0])
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	run(cmd)

}

func run(cmd *exec.Cmd) error {
	if !dryRun {
		return cmd.Run()
	}

	fmt.Println("(dry-run) Command skipped.")

	return nil
}
