package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	// Uncomment this block to pass the first stage
	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		args := strings.Fields(strings.TrimSpace(input[:len(input)-1]))

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
			// BAD CODE

			// for idx, arg := range args[1:] {
			// 	if idx == 0 {
			// 		fmt.Print(arg)
			// 	} else {
			// 		fmt.Print(" " + arg)
			// 	}
			// }
			// fmt.Print("\n")

			// GOOD CODE
			fmt.Println(strings.Join(args[1:], " "))

		case "type":
			if len(args) < 2 {
				fmt.Println("Invalid argument, Please use this format type <command>")
				continue
			}

			arg := args[1]

			switch arg {
			case "exit":
				fmt.Println("exit is a shell builtin")
			case "echo":
				fmt.Println("echo is a shell builtin")
			case "type":
				fmt.Println("type is a shell builtin")
			default:
				path, err := exec.LookPath(arg)
				if err != nil {
					fmt.Printf("%s: not found\n", arg)
					continue
				}

				fmt.Printf("%s is %s\n", arg, path)
				continue
			}

		default:
			inp := args[0]
			path, err := exec.LookPath(inp)

			if err != nil {
				fmt.Printf("%s: invalid command\n", inp)
				continue
			}

			fmt.Printf("Program was passed %d args (including program name).\n", len(args))

			for idx, arg := range args {
				if idx == 0 {
					fmt.Printf("Arg #%d (program name): %s\n", idx, arg)
				} else {
					fmt.Printf("Arg #%d: %s\n", idx, arg)
				}
			}

			cmd := exec.Command(path, args[:1]...)
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}
			continue
		}

	}

}
