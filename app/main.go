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
	PATH_ := strings.Split(os.Getenv("PATH"), ":")
	for {
		fmt.Fprint(os.Stdout, "$ ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Printf("%s: invalid input\n", input)
		}
		input_string := strings.TrimSpace(input)
		command := strings.Split(input_string, " ")[0]
		words := strings.Split(input_string, " ")
		if command == "exit" {
			break
		} else if command == "echo" {
			fmt.Println(strings.Join(words[1:], " "))
		} else if command == "type" && len(words) > 1 {
			found := false
			for _, path := range PATH_ {
				if strings.Contains(path, words[1]) {
					fmt.Println(words[1], "is", path)
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("%s: not found\n", words[1])
			}
		} else {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}
