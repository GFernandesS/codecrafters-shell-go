package main

import (
	"bufio"
	"bytes"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/helpers"
	"github.com/google/shlex"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

type Handler func(value string, valueTokens []string)

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

			if len(inputTokens) == 0 {
				fmt.Print("$ ")
				continue
			}

			commandName := inputTokens[0]

			if handler, foundCommand := builtInCommands[commandName]; foundCommand {
				handler(input, inputTokens)

				fmt.Print("$ ")
				continue
			}

			handleExternalExec(inputTokens)
		}

		fmt.Print("$ ")
	}
}

func handleExit(_ string, _ []string) {
	os.Exit(0)
}

func handleEcho(input string, inputTokens []string) {

	if slices.Contains(inputTokens, ">") || slices.Contains(inputTokens, "1>") {
		directionIndex := slices.Index(inputTokens, ">")

		if directionIndex == -1 {
			directionIndex = slices.Index(inputTokens, "1>")
		}

		contentToWrite := strings.Join(inputTokens[1:directionIndex], " ")

		fileToWrite := inputTokens[directionIndex+1]

		if err := os.MkdirAll(filepath.Dir(fileToWrite), 0755); err != nil {
			fmt.Println(err.Error())
			return
		}

		if err := os.WriteFile(fileToWrite, []byte(contentToWrite), 0644); err != nil {
			fmt.Printf(err.Error())
		}

		return
	}

	fmt.Printf("%s\n", strings.Replace(input, "echo ", "", 1))
}

func handlePwd(_ string, _ []string) {
	path, _ := os.Getwd()

	fmt.Printf("%s\n", path)
}

func handleType(input string, _ []string) {
	command := strings.Replace(input, "type ", "", 1)

	for c := range maps.Keys(builtInCommands) {
		if c == command {
			fmt.Printf("%s is a shell builtin\n", command)
			return
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
		return
	}

	fmt.Printf("%s: not found\n", command)
}

func handleCd(input string, _ []string) {
	command := strings.Replace(input, "cd ", "", 1)

	if command == "~" {
		command, _ = os.UserHomeDir()
	}

	err := os.Chdir(command)

	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", command)
	}

	return
}

func handleExternalExec(inputTokens []string) {
	pathsToSearch := strings.SplitSeq(os.Getenv("PATH"), ":")

	for path := range pathsToSearch {
		fullPath := path + "/" + inputTokens[0]

		commandInfo, err := os.Stat(fullPath)

		if err != nil {
			continue
		}

		hasExecPermissions := helpers.FileHasExecPermissions(commandInfo)

		if !hasExecPermissions {
			continue
		}

		var outputFile string

		if slices.Contains(inputTokens, ">") || slices.Contains(inputTokens, "1>") {
			directionIndex := slices.Index(inputTokens, ">")

			if directionIndex == -1 {
				directionIndex = slices.Index(inputTokens, "1>")
			}

			outputFile = inputTokens[directionIndex+1]

			inputTokens = append(inputTokens[:directionIndex])
		}

		cmd := exec.Command(inputTokens[0], inputTokens[1:]...)

		var outputBuffer bytes.Buffer

		var errBuffer bytes.Buffer

		cmd.Stdout = &outputBuffer

		cmd.Stderr = &errBuffer

		cmdErr := cmd.Run()

		if outputFile != "" {
			if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
				fmt.Printf(err.Error())
				return
			}

			if err := os.WriteFile(outputFile, outputBuffer.Bytes(), 0644); err != nil {
				fmt.Printf(err.Error())
			}

			if cmdErr != nil {
				fmt.Printf(errBuffer.String())
			}

			return
		}

		commandOutput := outputBuffer.String()

		if string(commandOutput[len(commandOutput)-1]) != "\n" {
			commandOutput = commandOutput + "\n"
		}

		fmt.Print(commandOutput)

		if cmdErr != nil {
			fmt.Printf(errBuffer.String())
		}

		return
	}

	fmt.Printf("%s: not found\n", inputTokens[0])
}

func sanitizeInput(input string) (string, []string) {
	inputTokens, _ := shlex.Split(input)

	return strings.Join(inputTokens, " "), inputTokens
}
