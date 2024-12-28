package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func exit() {
	os.Exit(0)
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Fatalf("error reading stdin")
		}

		command := input[:len(input)-1]

		if command == "exit" {
			exit()
		}

		fmt.Printf("%s: command not found\n", command)
	}
}
