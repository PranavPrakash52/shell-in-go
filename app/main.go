package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"golang.org/x/term"
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

func create_or_append_file(file_name string, appending bool) *os.File {
	if appending {
		file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Printf("%s: cannot redirect\n", file_name)
		}
		return file
	}
	file, err := os.Create(file_name)
	if err != nil {
		fmt.Printf("%s: cannot redirect\n", file_name)
	}
	return file
}

func execute_command(command string, args []string, out *os.File, std_out bool, std_err bool) {
	cmd_ := exec.Command(command, args...)
	cmd_.Stdin = os.Stdin
	if std_out {
		cmd_.Stdout = out
	} else {
		cmd_.Stdout = os.Stdout
	}
	if std_err {
		cmd_.Stderr = out
	} else {
		cmd_.Stderr = os.Stderr
	}
	cmd_.Run()
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
		separator := ""
		processedArgs = cleanup_args(args, special_char_double_quote)
		appending_seperator_list := []string{"1>>", ">>", "2>>", ">", "1>", "2>"}
		for _, seperator := range appending_seperator_list {
			if slices.Contains(processedArgs, seperator) {
				separator = seperator
				break
			}
		}
		if separator == "" {
			execute_command(processedArgs[0], processedArgs[1:], nil, false, false)
		} else {
			new_args = processedArgs[:slices.Index(processedArgs, separator)]
			file_name := processedArgs[slices.Index(processedArgs, separator)+1]
			var file *os.File
			switch separator {
			case ">", "1>":
				file = create_or_append_file(file_name, false)
				execute_command(new_args[0], new_args[1:], file, true, false)
			case "2>":
				file = create_or_append_file(file_name, false)
				execute_command(new_args[0], new_args[1:], file, false, true)
			case "1>>", ">>":
				file = create_or_append_file(file_name, true)
				execute_command(new_args[0], new_args[1:], file, true, false)
			case "2>>":
				file = create_or_append_file(file_name, true)
				execute_command(new_args[0], new_args[1:], file, false, true)
			}
			defer file.Close()
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

func autocomplete(input string, map_ map[string]string, builtin_map_ map[string]bool) string {
	if strings.Contains(input, " ") {
		return ""
	}
	for key := range builtin_map_ {
		if strings.HasPrefix(key, input) {
			return key[len(input):]
		}
	}
	for key := range map_ {
		if strings.HasPrefix(key, input) {
			return key[len(input):]
		}
	}
	fmt.Print("\x07")
	return ""
}

func read_input(map_ map[string]string, builtin_map_ map[string]bool) (input string) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	r := bufio.NewReader(os.Stdin)

loop:
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch c {
		case '\x03': // Ctrl+C
			os.Exit(0)
		case '\r', '\n': // Enter
			fmt.Fprint(os.Stdout, "\r\n")
			break loop
		case '\t': // Tab
			suffix := autocomplete(input, map_, builtin_map_)
			if suffix != "" {
				input += suffix + " "
				fmt.Fprint(os.Stdout, suffix+" ")
			}
		case '\x7F': // Backspace
			if length := len(input); length > 0 {
				input = input[:length-1]
				fmt.Fprint(os.Stdout, "\b \b")
			}
		default:
			input += string(c)
			fmt.Fprint(os.Stdout, string(c))
		}
	}
	return input
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
		input := read_input(map_, builtin_map_)
		if input == "" {
			continue
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
