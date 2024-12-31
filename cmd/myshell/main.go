package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func processArguments(input string, includeWhiteSpace bool) []string {
	res := []string{}

	isSingleQuoted := false
	isDoubleQuoted := false
	curr := ""
	index := 0
	for index < len(input) {
		r := input[index]
		if r == '\'' && !isDoubleQuoted {
			if isSingleQuoted {
				res = append(res, curr)
				curr = ""
				oldIndex := index
				for index+1 < len(input) && input[index+1] == ' ' {
					index += 1
				}
				if oldIndex < index && includeWhiteSpace {
					res = append(res, " ")
				}
			}
			isSingleQuoted = !isSingleQuoted
		} else if r == '"' && !isSingleQuoted {
			if isDoubleQuoted {
				curr = strings.ReplaceAll(curr, `\\`, `\`)
				curr = strings.ReplaceAll(curr, `\$`, `$`)
				curr = strings.ReplaceAll(curr, `\"`, `"`)
				curr = strings.ReplaceAll(curr, `\\n`, `\n`)
				res = append(res, curr)
				curr = ""
				oldIndex := index
				for index+1 < len(input) && input[index+1] == ' ' {
					index += 1
				}
				if oldIndex < index && includeWhiteSpace {
					res = append(res, " ")
				}
			}
			isDoubleQuoted = !isDoubleQuoted
		} else if r == ' ' && !isSingleQuoted && !isDoubleQuoted {
			curr = strings.Join(strings.Fields(strings.TrimSpace(curr)), " ")
			curr = strings.ReplaceAll(curr, `\`, "")
			res = append(res, curr)
			for index+1 < len(input) && input[index+1] == ' ' {
				index += 1
			}
			if includeWhiteSpace {
				res = append(res, " ")
			}
			curr = ""
		} else if r == '\\' {
			curr += string(input[index : index+2])
			index += 1
		} else {
			curr += string(r)
		}
		index += 1
	}

	if len(curr) != 0 {
		curr = strings.Join(strings.Fields(strings.TrimSpace(curr)), " ")
		curr = strings.ReplaceAll(curr, `\`, "")
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
	message := strings.Join(processArguments(input, true), "")
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
				cmd := exec.Command(path+"/"+command, processArguments(input, false)...)
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout
				err := cmd.Run()
				if err != nil {
				}
			}
		}
	}
}

func main() {
	initPathCommands()
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Fatalf("error reading stdin")
		}

		input = input[:len(input)-1]
		command := processArguments(input, false)[0]

		// Hack for passing the tests since the test executable is made after program start
		initPathCommands()

		if _, ok := COMMAND_DESCRIPTIONS[command]; ok {
			if len(input) > len(command) {
				input = string(input[len(command)+1:])
			} else {
				input = ""
			}

			var file *os.File
			var err error
			if strings.Contains(input, "2>>") {
				index := strings.Index(input, "2>>")
				filename := input[index+4:]
				if index > 0 {
					input = input[:index-1]
				} else {
					input = ""
				}

				file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Fatalf("error opening file")
				}

				os.Stderr = file
			} else if strings.Contains(input, ">>") || strings.Contains(input, "1>>") {
				var filename string
				if strings.Contains(input, "1>>") {
					index := strings.Index(input, "1>>")
					filename = input[index+4:]
					if index > 0 {
						input = input[:index-1]
					} else {
						input = ""
					}
				} else if strings.Contains(input, ">>") {
					index := strings.Index(input, ">>")
					filename = input[index+3:]
					if index > 0 {
						input = input[:index-1]
					} else {
						input = ""
					}
				}

				file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Fatalf("error opening file")
				}

				os.Stdout = file
			} else if strings.Contains(input, "2>") {
				index := strings.Index(input, "2>")
				filename := input[index+3:]
				if index > 0 {
					input = input[:index-1]
				} else {
					input = ""
				}

				file, err = os.Create(filename)
				if err != nil {
					log.Fatalf("error opening file")
				}

				os.Stderr = file
			} else if strings.Contains(input, ">") || strings.Contains(input, "1>") {
				var filename string
				if strings.Contains(input, "1>") {
					index := strings.Index(input, "1>")
					filename = input[index+3:]
					if index > 0 {
						input = input[:index-1]
					} else {
						input = ""
					}
				} else if strings.Contains(input, ">") {
					index := strings.Index(input, ">")
					filename = input[index+2:]
					if index > 0 {
						input = input[:index-1]
					} else {
						input = ""
					}
				}

				file, err = os.Create(filename)
				if err != nil {
					log.Fatalf("error opening file")
				}

				os.Stdout = file
			}

			handler := COMMAND_FUNCTIONS[command]
			handler(input)

			file.Close()
			os.Stdout = originalStdout
			os.Stderr = originalStderr
		} else {
			fmt.Printf("%s: command not found\n", input)
		}
	}
}
