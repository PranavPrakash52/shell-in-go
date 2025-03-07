package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
// var _ = fmt.Fprint

func main() {
	// Uncomment this block to pass the first stage
	// Wait for user input
	PATH_ := strings.Split(os.Getenv("PATH"), ":")
	map_ := make(map[string]string)
	builtin_list := []string{
		"echo", "cd", "pwd", "exit", "type", "alias", "bg", "bind", "break", 
		"builtin", "caller", "command", "compgen", "complete", "compopt", 
		"continue", "declare", "dirs", "disown", "enable", "eval", "exec", 
		"export", "false", "fc", "fg", "getopts", "hash", "help", "history", 
		"jobs", "kill", "let", "local", "logout", "mapfile", "popd", "printf", 
		"pushd", "read", "readarray", "readonly", "return", "set", "shift", 
		"shopt", "source", "suspend", "test", "times", "trap", "true", "type", 
		"typeset", "ulimit", "umask", "unalias", "unset", "wait",
	}
	builtin_map_ := make(map[string]bool)
	for _, builtin := range builtin_list {
		builtin_map_[builtin] = true
	}
	for _, path := range PATH_ {
		entries, err := os.ReadDir(path)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if _, ok := map_[entry.Name()]; ok {
				continue
			}
			map_[entry.Name()] = path
		}
	}
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
			if _, ok := builtin_map_[words[1]]; ok {
				fmt.Println(words[1], "is", "builtin")
			} else if _, ok := map_[words[1]]; ok {
				fmt.Println(words[1], "is", map_[words[1]]+"/"+words[1])
			} else {
				fmt.Printf("%s: not found\n", words[1])
			}
		} else {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}
