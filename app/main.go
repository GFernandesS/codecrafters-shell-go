package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("$ ")

	for true {
		if scanner.Scan() {
			input := scanner.Text()

			switch input {
			case "exit":
				os.Exit(0)
			default:
				fmt.Printf("%s: command not found\n", input)
			}

		}

		fmt.Print("$ ")
	}

}
