package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
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
		file, err := os.Open(path)
		if err != nil {
			continue
		}

		commands, err := file.Readdirnames(0)
		if err != nil {
			log.Fatalf("error listing directory")
		}

		fmt.Println("[DEBUG] - Commands:", commands)

		for _, command := range commands {
			if _, ok := COMMANDS[command]; ok {
				continue
			}

			COMMANDS[command] = path + "/" + command
			COMMAND_MAPPINGS[command] = func(input string) {
				cmd := exec.Command(path+"/"+command, strings.Split(input, " ")...)
				out, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Println(string(out))
					log.Fatalf("error running command %v", err)
				}
				fmt.Println(string(out))
			}
		}
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
				if len(input) > len(command) {
					input = string(input[len(command)+1:])
				} else {
					input = ""
				}
				handler(input)
				isValidCommand = true
			}
		}

		if !isValidCommand {
			fmt.Printf("%s: command not found\n", input)
		}
	}
}
