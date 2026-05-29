package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

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
