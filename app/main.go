package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
// var _ = fmt.Fprint
func run_pwd() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("%s: command not found\n", "pwd")
	}
	fmt.Println(path)
}

func run_cd(path string) {
	if path == "~" {
		path = os.Getenv("HOME")
	}
	err := os.Chdir(path)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", path)
	}
}

func cleanup_args(args []string, special_char_double_quote map[string]bool) []string {
	inputString := strings.Join(args, " ")
	var currentArg strings.Builder
	inQuotes := false
	inDoubleQuotes := false
	processedArgs := []string{}
	for i := 0; i < len(inputString); i++ {
		char := inputString[i]
		if char == '"' && !inQuotes {
			inDoubleQuotes = !inDoubleQuotes
			continue
		}
		if char == '\'' && !inDoubleQuotes {
			inQuotes = !inQuotes
			continue // Skip the quote character
		}

		if char == ' ' && !inDoubleQuotes && !inQuotes {
			// Space outside quotes means end of current argument
			if currentArg.Len() > 0 {
				processedArgs = append(processedArgs, currentArg.String())
				currentArg.Reset()
			}
		} else {
			if char == '\\' && !inQuotes {
				if inDoubleQuotes && special_char_double_quote[string(inputString[i+1])] {
					char = inputString[i+1]
					i += 1
				} else if char == '\\' && !inDoubleQuotes && !inQuotes {
					char = inputString[i+1]
					i += 1
				}
			}
			currentArg.WriteByte(char)
		}
	}
	// Add the last argument if there is one
	if currentArg.Len() > 0 {
		processedArgs = append(processedArgs, currentArg.String())
	}
	return processedArgs
}

func run_command(command string, args []string) {
	var cmd_ *exec.Cmd
	special_char_double_quote := map[string]bool{
		"\\": true,
		"$":  true,
		"\"": true,
	}
	if len(args) == 1 {
		cmd_ = exec.Command(command)
		cmd_.Stdin = os.Stdin
		cmd_.Stdout = os.Stdout
		cmd_.Stderr = os.Stderr
		err := cmd_.Run()
		if err != nil {
			fmt.Printf("%s: command not found\n", command)
		}
	} else {
		processedArgs := []string{}
		new_args := []string{}
		separator := ">"
		processedArgs = cleanup_args(args, special_char_double_quote)
		if slices.Contains(processedArgs, ">") || slices.Contains(processedArgs, "1>") {
			if slices.Contains(processedArgs, "1>") {
				separator = "1>"
			}
			new_args = processedArgs[:slices.Index(processedArgs, separator)]
			file_name := processedArgs[slices.Index(processedArgs, separator)+1]
			file, err := os.Create(file_name)
			if err != nil {
				fmt.Printf("%s: cannot redirect\n", command)
			}
			cmd_ = exec.Command(new_args[0], new_args[1:]...)
			cmd_.Stdin = os.Stdin
			cmd_.Stdout = file
			cmd_.Stderr = os.Stderr
			err = cmd_.Run()
			// if err != nil {
			// 	fmt.Printf("%s: command not found\n", command)
			// }
		} else {
			cmd_ = exec.Command(processedArgs[0], processedArgs[1:]...)
			cmd_.Stdin = os.Stdin
			cmd_.Stdout = os.Stdout
			cmd_.Stderr = os.Stderr
			err := cmd_.Run()
			if err != nil {
				fmt.Printf("%s: command not found\n", command)
			}
		}
	}
}

func type_command(command string, words []string, map_ map[string]string, builtin_map_ map[string]bool) {
	if _, ok := builtin_map_[words[1]]; ok {
		fmt.Println(words[1], "is a shell builtin")
	} else if _, ok := map_[words[1]]; ok {
		fmt.Println(words[1], "is", map_[words[1]]+"/"+words[1])
	} else {
		fmt.Printf("%s: not found\n", words[1])
	}
}

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
		args := strings.Split(input_string, " ")
		command := args[0]
		if command == "exit" {
			break
		} else if command == "echo" {
			run_command(command, args)
		} else if command == "type" && len(args) > 1 {
			type_command(command, args, map_, builtin_map_)
		} else if command == "pwd" {
			run_pwd()
		} else if command == "cd" {
			run_cd(args[1])
		} else {
			run_command(command, args)
		}

	}
}
