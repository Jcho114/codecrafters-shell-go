package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func executeExit(input string) {
	os.Exit(0)
}

func executeEcho(input string) {
	message := input
	fmt.Println(message)
}

var COMMANDS = map[string]bool{
	"exit": true,
	"echo": true,
	"type": true,
}

func executeType(input string) {
	command := input
	if _, ok := COMMANDS[command]; ok {
		fmt.Println(command, "is a shell builtin")
	} else {
		fmt.Printf("%s: not found\n", command)
	}
}

var COMMAND_MAPPINGS = map[string]func(string){
	"exit 1": executeExit,
	"echo":   executeEcho,
	"type":   executeType,
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Fatalf("error reading stdin")
		}

		input = input[:len(input)-1]

		isValidCommand := false
		for command, handler := range COMMAND_MAPPINGS {
			if strings.HasPrefix(input, command) {
				handler(input[len(command)+1:])
				isValidCommand = true
			}
		}

		if !isValidCommand {
			fmt.Printf("%s: command not found\n", input)
		}
	}
}
