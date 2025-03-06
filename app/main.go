package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	// Uncomment this block to pass the first stage
	// Wait for user input
	for {
		fmt.Fprint(os.Stdout, "$ ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Printf("%s: invalid input\n", input)
		}
		command := strings.TrimSpace(input)
		if command == "exit 0" {
			break
		} else if strings.Contains(command, "echo") {
			words := strings.Split(command, " ")
			fmt.Println(strings.Join(words[1:], " "))
		} else {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}
