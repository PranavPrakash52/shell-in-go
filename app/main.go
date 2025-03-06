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
		input_string := strings.TrimSpace(input)
		command := strings.Split(input_string, " ")[0]
		if command == "exit" {
			break
		} else if command == "echo" {
			words := strings.Split(input_string, " ")
			fmt.Println(strings.Join(words[1:], " "))
		} else if command == "type" {
			words := strings.Split(input_string, " ")
			if words[1] == "echo" {
				fmt.Println("echo is a shell builtin")
			} else if words[1] == "exit" {
				fmt.Println("exit is a shell builtin")
			} else if words[1] == "type" {
				fmt.Println("type is a shell builtin")
			} else {
				fmt.Printf("%s: command not found\n", words[1])
			}
		} else {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}
