package main

import (
	"fmt"
	"testing"
)

func TestFSM(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expected   bool
		finalState int
	}{
		{
			name:       "Valid string literal",
			input:      "\"hello world\"",
			expected:   true,
			finalState: StateS2,
		},
		{
			name:       "Valid empty string",
			input:      "\"\"",
			expected:   true,
			finalState: StateS2,
		},
		{
			name:       "Valid identifier expression",
			input:      "abc+def",
			expected:   true,
			finalState: StateA2,
		},
		{
			name:       "Valid identifier with numbers",
			input:      "abc123+def456",
			expected:   true,
			finalState: StateA2,
		},
		{
			name:       "Valid integer expression",
			input:      "123+456",
			expected:   true,
			finalState: StateB2,
		},
		{
			name:       "Valid float expression",
			input:      "123.45+67.89",
			expected:   true,
			finalState: StateE2,
		},
		{
			name:       "Valid mixed expression",
			input:      "abc+123.45",
			expected:   true,
			finalState: StateE2,
		},
		{
			name:       "Invalid: just identifier",
			input:      "abc",
			expected:   false,
			finalState: StateA1,
		},
		{
			name:       "Invalid: just number",
			input:      "123",
			expected:   false,
			finalState: StateB1,
		},
		{
			name:       "Invalid: just operator",
			input:      "+",
			expected:   false,
			finalState: -1,
		},
		{
			name:       "Invalid: missing second operand",
			input:      "abc+",
			expected:   false,
			finalState: StateOp,
		},
		{
			name:       "Invalid: missing first operand",
			input:      "+abc",
			expected:   false,
			finalState: -1,
		},
		{
			name:       "Invalid: invalid character",
			input:      "abc@def",
			expected:   false,
			finalState: -1,
		},
		{
			name:       "Invalid: unclosed string",
			input:      "\"hello",
			expected:   false,
			finalState: StateS1,
		},
		{
			name:       "Invalid: invalid float",
			input:      "123.+456",
			expected:   false,
			finalState: -1,
		},
		{
			name:       "Valid: multiple operators",
			input:      "a+b+c",
			expected:   false,
			finalState: -1,
		},
		{
			name:       "Valid: float with leading zero",
			input:      "0.123+456",
			expected:   true,
			finalState: StateB2,
		},
		{
			name:       "Invalid: float with trailing dot",
			input:      "123.+456",
			expected:   false,
			finalState: -1,
		},
		{
			name:       "Invalid: float with multiple dots",
			input:      "123.45.67",
			expected:   false,
			finalState: -1,
		},
		{
			name:       "Valid: identifier with underscore",
			input:      "abc_123+def",
			expected:   false,
			finalState: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := StateStart
			errorFlag := false

			fmt.Printf("Testing input: %q\n", tt.input)
			for _, ch := range tt.input {
				state, errorFlag = FSM(ch, state)
				fmt.Printf("char: %q, state: %d, error: %v\n", ch, state, errorFlag)
				if errorFlag {
					break
				}
			}

			finalStates := map[int]bool{
				StateA2: true,
				StateB2: true,
				StateE2: true,
				StateS2: true,
			}

			accepting := finalStates[state]
			fmt.Printf("Final state: %d, accepting: %v\n", state, accepting)

			if accepting != tt.expected {
				t.Errorf("FSM(%q) = %v; want %v", tt.input, accepting, tt.expected)
			}
			if state != tt.finalState {
				t.Errorf("Final state = %d; want %d", state, tt.finalState)
			}
		})
	}
}

// Test individual state transitions
func TestStateTransitions(t *testing.T) {
	tests := []struct {
		name     string
		state    int
		input    rune
		expected int
		error    bool
	}{
		{
			name:     "Start to String",
			state:    StateStart,
			input:    '"',
			expected: StateS1,
			error:    false,
		},
		{
			name:     "Start to Identifier",
			state:    StateStart,
			input:    'a',
			expected: StateA1,
			error:    false,
		},
		{
			name:     "Start to Number",
			state:    StateStart,
			input:    '1',
			expected: StateB1,
			error:    false,
		},
		{
			name:     "Start to Invalid",
			state:    StateStart,
			input:    '@',
			expected: -1,
			error:    true,
		},
		{
			name:     "Identifier to Identifier",
			state:    StateA1,
			input:    'b',
			expected: StateA1,
			error:    false,
		},
		{
			name:     "Identifier to Operator",
			state:    StateA1,
			input:    '+',
			expected: StateOp,
			error:    false,
		},
		{
			name:     "Number to Number",
			state:    StateB1,
			input:    '2',
			expected: StateB1,
			error:    false,
		},
		{
			name:     "Number to Decimal",
			state:    StateB1,
			input:    '.',
			expected: StateD1,
			error:    false,
		},
		{
			name:     "String to String",
			state:    StateS1,
			input:    'a',
			expected: StateS1,
			error:    false,
		},
		{
			name:     "String to End",
			state:    StateS1,
			input:    '"',
			expected: StateS2,
			error:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newState, errorFlag := FSM(tt.input, tt.state)
			if newState != tt.expected || errorFlag != tt.error {
				t.Errorf("FSM(%d, %q) = (%d, %v); want (%d, %v)",
					tt.state, tt.input, newState, errorFlag, tt.expected, tt.error)
			}
		})
	}
}
