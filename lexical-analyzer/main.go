package main

import (
	"fmt"
	"os"

	"analyzer/fsmlex"
	"analyzer/models"
	"analyzer/rxlex"
)

func main() {
	var input []byte
	var err error

	if len(os.Args) > 2 {
		input, err = os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		var tokens []models.Token
		if os.Args[2] == "fsm" {
			tokens = fsmlex.Lex(string(input))
		} else if os.Args[2] == "rx" {
			tokens, err = rxlex.Lex(string(input))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

		for _, token := range tokens {
			fmt.Printf("%s: %s\n", token.Type, token.Value)
		}
	} else {
		fmt.Println("Not enough params. Example: lexer ./examples/example.go rx")
	}

}
