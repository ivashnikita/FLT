package models

type TokenType string

const (
	Keyword        TokenType = "Keyword"
	Identifier     TokenType = "Identifier"
	IntLiteral     TokenType = "Int"
	FloatLiteral   TokenType = "Float"
	StringLiteral  TokenType = "String"
	RuneLiteral    TokenType = "Rune"
	BooleanLiteral TokenType = "Boolean"
	Operator       TokenType = "Operator"
	Separator      TokenType = "Separator"
	Error          TokenType = "ERROR"
)

type Token struct {
	Type  TokenType
	Value string
}
