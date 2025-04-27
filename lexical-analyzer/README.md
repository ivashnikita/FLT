## How run lexer

You can write your own code on golang and add it in examples folder.

For getting lexemes you need:
```
go build -o lexer main.go
./lexer ./examples/example.go 
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