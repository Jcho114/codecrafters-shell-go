package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Fatalf("error reading stdin")
		}
		fmt.Printf("%s: command not found\n", input[:len(input)-1])
	}
}
