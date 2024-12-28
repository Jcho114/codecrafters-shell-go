package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Fprint(os.Stdout, "$ ")

	for {
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Fatalf("error reading stdin")
		}
		fmt.Printf("%s: command not found\n", input[:len(input)-1])
	}
}
