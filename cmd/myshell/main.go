package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func exit() {
	os.Exit(0)
}

func echo(message string) {
	fmt.Println(message)
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Fatalf("error reading stdin")
		}

		input = input[:len(input)-1]

		if input == "exit 0" {
			exit()
		} else if strings.HasPrefix(input, "echo") {
			echo(input[5:])
		} else {
			fmt.Printf("%s: command not found\n", input)
		}
	}
}
