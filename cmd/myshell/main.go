package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func processArguments(input string) []string {
	res := []string{}

	isSingleQuoted := false
	isDoubleQuoted := false
	curr := ""
	for _, r := range input {
		if r == '\'' && !isDoubleQuoted {
			if isSingleQuoted {
				res = append(res, curr)
			}
			isSingleQuoted = !isSingleQuoted
			curr = ""
		} else if r == '"' {
			if isDoubleQuoted {
				res = append(res, curr)
			}
			isDoubleQuoted = !isDoubleQuoted
			curr = ""
		} else if r == ' ' && !isSingleQuoted && !isDoubleQuoted && curr != "" {
			var err error
			curr = strings.ReplaceAll(curr, `\ `, " ")
			curr, err = strconv.Unquote("\"" + curr + "\"")
			if err != nil {
				log.Fatalf("error unquoting string")
			}
			res = append(res, curr)
			curr = ""
		} else {
			curr += string(r)
		}
	}

	if len(curr) != 0 {
		res = append(res, curr)
	}

	return res
}

var COMMAND_DESCRIPTIONS = map[string]string{
	"exit": "a shell builtin",
	"echo": "a shell builtin",
	"type": "a shell builtin",
	"pwd":  "a shell builtin",
	"cd":   "a shell builtin",
}

func executeExit(input string) {
	os.Exit(0)
}

func executeEcho(input string) {
	var message string
	var err error
	if strings.Contains(input, "'") || strings.Contains(input, "\"") {
		message = strings.Join(processArguments(input), "")
	} else {
		message = strings.ReplaceAll(input, `\ `, " ")
		message, err = strconv.Unquote("\"" + message + "\"")
		if err != nil {
			log.Fatalf("error unquoting string %v", err)
		}
	}
	fmt.Println(message)
}

func executeType(input string) {
	command := input
	if description, ok := COMMAND_DESCRIPTIONS[command]; ok {
		fmt.Printf("%s is %s\n", command, description)
	} else {
		fmt.Printf("%s: not found\n", command)
	}
}

func executePwd(input string) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("error in getting current directory")
	}
	fmt.Println(currentDirectory)
}

func executeCd(input string) {
	var directory string
	var err error

	if input == "~" {
		directory, err = os.UserHomeDir()
		if err != nil {
			log.Fatalf("error obtaining home directory")
		}
	} else {
		directory = input
	}

	_, err = os.Stat(directory)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", directory)
		return
	}

	err = os.Chdir(directory)
	if err != nil {
		log.Fatalf("error changing directories")
	}
}

var COMMAND_FUNCTIONS = map[string]func(string){
	"exit": executeExit,
	"echo": executeEcho,
	"type": executeType,
	"pwd":  executePwd,
	"cd":   executeCd,
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

		for _, command := range commands {
			if _, ok := COMMAND_DESCRIPTIONS[command]; ok {
				continue
			}

			COMMAND_DESCRIPTIONS[command] = path + "/" + command
			COMMAND_FUNCTIONS[command] = func(input string) {
				cmd := exec.Command(path+"/"+command, processArguments(input)...)
				out, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Println(string(out))
					log.Fatalf("error running command %v", err)
				}
				if len(out) == 0 {
					fmt.Println(string(out))
				} else {
					// removed \n at the end for tests to pass
					fmt.Println(string(out[:len(out)-1]))
				}
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
		command := strings.Split(input, " ")[0]

		// Hack for passing the tests since the test executable is made after program start
		initPathCommands()

		if _, ok := COMMAND_DESCRIPTIONS[command]; ok {
			if len(input) > len(command) {
				input = string(input[len(command)+1:])
			} else {
				input = ""
			}
			handler := COMMAND_FUNCTIONS[command]
			handler(input)
		} else {
			fmt.Printf("%s: command not found\n", input)
		}
	}
}
