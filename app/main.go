package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	// Uncomment this block to pass the first stage
	// Wait for user input
	bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Fprint(os.Stdout, "invalid_command: command not found")
}
