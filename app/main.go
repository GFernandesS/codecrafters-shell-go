package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/helpers"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

var builtInCommands = []string{"echo", "exit", "type"}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("$ ")

	for {
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

			handleExternalExec(input)
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

		hasExecPermissions := helpers.FileHasExecPermissions(commandInfo)

		if !hasExecPermissions {
			continue
		}

		fmt.Printf("%s is %s\n", commandInfo.Name(), fullPath)
		return true
	}

	fmt.Printf("%s: not found\n", command)
	return true
}

func handleExternalExec(input string) {
	pathsToSearch := strings.Split(os.Getenv("PATH"), ":")

	commandParts := strings.Split(input, " ")

	for _, path := range pathsToSearch {
		fullPath := path + "/" + commandParts[0]

		commandInfo, err := os.Stat(fullPath)

		if err != nil {
			continue
		}

		hasExecPermissions := helpers.FileHasExecPermissions(commandInfo)

		if !hasExecPermissions {
			continue
		}

		cmd := exec.Command(fullPath, commandParts[1:]...)

		out, err := cmd.CombinedOutput()

		if err != nil {
			return
		}

		fmt.Printf("%s", string(out))
		return
	}

	fmt.Printf("%s: not found\n", commandParts[0])
}
