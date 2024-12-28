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

var COMMANDS = map[string]string{
	"exit": "a shell builtin",
	"echo": "a shell builtin",
	"type": "a shell builtin",
}

func executeType(input string) {
	command := input
	if description, ok := COMMANDS[command]; ok {
		fmt.Printf("%s is %s\n", command, description)
	} else {
		fmt.Printf("%s: not found\n", command)
	}
}

var COMMAND_MAPPINGS = map[string]func(string){
	"exit": executeExit,
	"echo": executeEcho,
	"type": executeType,
}

func initPathCommands() {
	pathString := os.Getenv("PATH")
	paths := strings.Split(pathString, ":")
	for _, path := range paths {
		split := strings.Split(path, "/")
		command := split[len(split)-1]
		COMMANDS[command] = path
	}
}

func main() {
	initPathCommands()

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
