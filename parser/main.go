package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

// Lexeme types
const (
	ILLEGAL = iota
	EOF
	IDENTIFIER
	INTEGER
	FLOAT
	OPERATOR
	STRING
)

// FSM states
const (
	StateStart = iota
	StateA1    // First identifier
	StateB1    // First integer
	StateD1    // First decimal point
	StateE1    // First fractional part
	StateOp    // Operator
	StateA2    // Second identifier
	StateB2    // Second integer
	StateD2    // Second decimal point
	StateE2    // Second fractional part
	StateS1    // String start
	StateS2    // String end
)

func isAlpha(ch rune) bool {
	return unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func isOperator(ch rune) bool {
	return ch == '+' || ch == '-' || ch == '*' || ch == '/'
}

// FSM function that processes one character at a time
// Checking I=(((I | DF) ("+" | "-" | "/" | "*") (I | DF)) | S) form
func FSM(ch rune, state int) (int, bool) {
	switch state {
	case StateStart:
		if ch == '"' {
			return StateS1, false
		} else if isAlpha(ch) {
			return StateA1, false
		} else if isDigit(ch) {
			return StateB1, false
		}
		return -1, true

	case StateA1:
		if isAlpha(ch) || isDigit(ch) {
			return StateA1, false
		} else if isOperator(ch) {
			return StateOp, false
		}
		return -1, true

	case StateB1:
		if isDigit(ch) {
			return StateB1, false
		} else if ch == '.' {
			return StateD1, false
		} else if isOperator(ch) {
			return StateOp, false
		}
		return -1, true

	case StateD1:
		if isDigit(ch) {
			return StateE1, false
		}
		return -1, true

	case StateE1:
		if isDigit(ch) {
			return StateE1, false
		} else if isOperator(ch) {
			return StateOp, false
		}
		return -1, true

	case StateOp:
		if isAlpha(ch) {
			return StateA2, false
		} else if isDigit(ch) {
			return StateB2, false
		}
		return -1, true

	case StateA2:
		if isAlpha(ch) || isDigit(ch) {
			return StateA2, false
		}
		return -1, true

	case StateB2:
		if isDigit(ch) {
			return StateB2, false
		} else if ch == '.' {
			return StateD2, false
		}
		return -1, true

	case StateD2:
		if isDigit(ch) {
			return StateE2, false
		}
		return -1, true

	case StateE2:
		if isDigit(ch) {
			return StateE2, false
		}
		return -1, true

	case StateS1:
		if ch == '"' {
			return StateS2, false
		}
		return StateS1, false

	case StateS2:
		return -1, true

	default:
		return -1, true
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	state := StateStart
	finalStates := map[int]bool{
		StateA2: true,
		StateB2: true,
		StateE2: true,
		StateS2: true,
	}

	for scanner.Scan() {
		input := scanner.Text()
		if input == "" {
			break
		}

		state = StateStart
		for _, ch := range input {
			var errorFlag bool
			state, errorFlag = FSM(ch, state)
			if errorFlag {
				break
			}
		}

		if finalStates[state] {
			fmt.Println("Yes")
		} else {
			fmt.Println("No")
		}
	}
}
