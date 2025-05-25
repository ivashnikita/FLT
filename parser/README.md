# FSM Parser

A Go implementation of a Finite State Machine (FSM) for parsing expressions according to the grammar:
```
I=(((I | DF) ("+" | "-" | "/" | "*") (I | DF)) | S)
```

## Running

If you want to put your input:
```bash
go run main.go
```
Print `Enter` to stop programm

To run all tests:
```bash
go test -v
```

This will run:
- TestFSM: Tests various valid and invalid expressions
- TestStateTransitions: Tests individual state transitions

Example test output:
```
=== RUN   TestFSM/Valid_identifier_expression
Testing input: "abc+def"
char: 'a', state: 1, error: false
char: 'b', state: 1, error: false
char: 'c', state: 1, error: false
char: '+', state: 5, error: false
char: 'd', state: 6, error: false
char: 'e', state: 6, error: false
char: 'f', state: 6, error: false
Final state: 6, accepting: true
```
