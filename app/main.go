package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/helpers"
	"github.com/google/shlex"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

type Handler func(value string) bool

func init() {
	builtInCommands = map[string]Handler{
		"exit": handleExit,
		"echo": handleEcho,
		"type": handleType,
		"pwd":  handlePwd,
		"cd":   handleCd,
	}
}

var builtInCommands map[string]Handler

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("$ ")

	for {
		if scanner.Scan() {
			input := scanner.Text()

			input, inputTokens := sanitizeInput(input)

			handleExit(input)

			var wasHandle bool

			for _, handler := range builtInCommands {
				wasHandle = handler(input)

				if wasHandle {
					break
				}
			}

			if wasHandle {
				fmt.Print("$ ")
				continue
			}

			handleExternalExec(inputTokens)
		}

		fmt.Print("$ ")
	}
}

func handleExit(input string) bool {
	if input != "exit" {
		return false
	}

	os.Exit(0)

	return true
}

func handleEcho(input string) bool {
	if !strings.HasPrefix(input, "echo") {
		return false
	}

	fmt.Printf("%s\n", strings.Replace(input, "echo ", "", 1))
	return true
}

func handlePwd(input string) bool {
	if !strings.HasPrefix(input, "pwd") {
		return false
	}

	path, _ := os.Getwd()

	fmt.Printf("%s\n", path)

	return true
}

func handleType(input string) bool {
	if !strings.HasPrefix(input, "type") {
		return false
	}

	command := strings.Replace(input, "type ", "", 1)

	for c := range maps.Keys(builtInCommands) {
		if c == command {
			fmt.Printf("%s is a shell builtin\n", command)
			return true
		}
	}

	pathsToSearch := strings.SplitSeq(os.Getenv("PATH"), ":")

	for path := range pathsToSearch {
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

func handleCd(input string) bool {
	if !strings.HasPrefix(input, "cd") {
		return false
	}

	command := strings.Replace(input, "cd ", "", 1)

	if command == "~" {
		command, _ = os.UserHomeDir()
	}

	err := os.Chdir(command)

	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", command)
	}

	return true
}

func handleExternalExec(inputTokens []string) {
	pathsToSearch := strings.Split(os.Getenv("PATH"), ":")

	for _, path := range pathsToSearch {
		fullPath := path + "/" + inputTokens[0]

		commandInfo, err := os.Stat(fullPath)

		if err != nil {
			continue
		}

		hasExecPermissions := helpers.FileHasExecPermissions(commandInfo)

		if !hasExecPermissions {
			continue
		}

		cmd := exec.Command(inputTokens[0], inputTokens[1:]...)

		out, _ := cmd.CombinedOutput()

		fmt.Printf("%s", string(out))
		return
	}

	fmt.Printf("%s: not found\n", inputTokens[0])
}

func sanitizeInput(input string) (string, []string) {
	inputTokens, _ := shlex.Split(input)

	return strings.Join(inputTokens, " "), inputTokens
}
