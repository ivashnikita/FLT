## How run lexer

You can write your own code on golang and add it in examples folder.

There are two realization of lexer: 
- based on regular expression
- based on finite state machine

For getting lexemes you need build binary at first
```
go build -o lexer main.go
```

If you want to use regexp based lexer:
```
./lexer ./examples/example.go rx
```

If you want to use finite state machine based lexer:
```
./lexer ./examples/example.go fsm
```

You get something like:
```
Keyword: package
Identifier: examples
Keyword: import
String: "fmt"
Keyword: func
Identifier: main
Separator: (
Separator: )
Separator: {
Identifier: fmt
Separator: .
Identifier: Println
...
```