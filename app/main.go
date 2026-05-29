package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

var builtInCommands = []string{"echo", "exit", "type"}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("$ ")

	for true {
		if scanner.Scan() {
			input := scanner.Text()

			if input == "exit" {
				handleExit()
			}

			wasHandle := handleEcho(input)

			if wasHandle {
				fmt.Print("$ ")
				continue
			}

			wasHandle = handleType(input)

			if wasHandle {
				fmt.Print("$ ")
				continue
			}

			fmt.Printf("%s: command not found\n", input)
		}

		fmt.Print("$ ")
	}
}

func handleExit() {
	os.Exit(0)
}

func handleEcho(input string) bool {
	if !strings.HasPrefix(input, "echo") {
		return false
	}

	fmt.Printf("%s\n", strings.Replace(input, "echo ", "", 1))
	return true
}

func handleType(input string) bool {
	if !strings.HasPrefix(input, "type") {
		return false
	}

	command := strings.Replace(input, "type ", "", 1)

	for _, c := range builtInCommands {
		if c == command {
			fmt.Printf("%s is a shell builtin\n", command)
			return true
		}
	}

	pathsToSearch := strings.Split(os.Getenv("PATH"), ":")

	for _, path := range pathsToSearch {
		fullPath := path + "/" + command

		commandInfo, err := os.Stat(fullPath)

		if err != nil {
			continue
		}

		hasExecPermissions := commandInfo.Mode().Perm()&0111 != 0

		if !hasExecPermissions {
			continue
		}

		fmt.Printf("%s is %s\n", commandInfo.Name(), fullPath)
		return true
	}

	fmt.Printf("%s: not found\n", command)
	return true
}
