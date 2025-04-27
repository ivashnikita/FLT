package main

import (
	"fmt"
	"io"
	"os"

	"analyzer/rxlex"
)

func main() {
	var input []byte
	var err error

	if len(os.Args) > 1 {
		input, err = os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
	} else {
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error reading stdin:", err)
			return
		}
	}

	tokens, err := rxlex.Lex(string(input))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, token := range tokens {
		fmt.Printf("%s: %s\n", token.Type, token.Value)
	}
}
