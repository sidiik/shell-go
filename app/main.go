package main

import (
	"bufio"
	"fmt"
	"os"
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

		default:
			fmt.Printf("%s: command not found\n", args[0])
			continue
		}

	}

}
