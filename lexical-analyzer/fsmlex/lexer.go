package fsmlex

import (
	"strings"
	"unicode"

	"analyzer/models"
)

type State int

const (
	Start State = iota
	InIdentifier
	InNumber
	InHexNumber
	InOctalNumber
	InBinaryNumber
	InFloat
	InExponent
	InExponentDigits
	InString
	InRawString
	InRune
	InOperator
	InLineComment
	InBlockComment
)

var keywords = map[string]bool{
	"break": true, "case": true, "chan": true, "const": true, "continue": true,
	"default": true, "defer": true, "else": true, "fallthrough": true,
	"for": true, "func": true, "go": true, "goto": true, "if": true,
	"import": true, "interface": true, "map": true, "package": true,
	"range": true, "return": true, "select": true, "struct": true,
	"switch": true, "type": true, "var": true,
}

var operators = map[string]bool{
	"+": true, "-": true, "*": true, "/": true, "%": true,
	"&": true, "|": true, "^": true, "<<": true, ">>": true,
	"&^": true, "+=": true, "-=": true, "*=": true, "/=": true,
	"%=": true, "&=": true, "|=": true, "^=": true, "<<=": true,
	">>=": true, "&^=": true, "&&": true, "||": true, "<-": true,
	"++": true, "--": true, "==": true, "<": true, ">": true,
	"=": true, "!": true, "!=": true, "<=": true, ">=": true,
	":=": true, "...": true,
}

var separators = map[rune]bool{
	'(': true, ')': true, '[': true, ']': true, '{': true,
	'}': true, ',': true, ';': true, ':': true, '.': true,
}

func Lex(input string) []models.Token {
	var tokens []models.Token
	state := Start
	var buffer strings.Builder
	var strDelim rune
	pos := 0

	for pos < len(input) {
		ch := rune(input[pos])

		switch state {
		case Start:
			if unicode.IsSpace(ch) {
				pos++
				continue
			}

			if ch == '/' && pos+1 < len(input) {
				next := input[pos+1]
				if next == '/' {
					state = InLineComment
					pos += 2
					continue
				}
				if next == '*' {
					state = InBlockComment
					pos += 2
					continue
				}
			}

			if isOperatorStart(ch) {
				state = InOperator
				buffer.WriteRune(ch)
				pos++
				continue
			}

			if separators[ch] {
				tokens = append(tokens, models.Token{Type: models.Separator, Value: string(ch)})
				pos++
				continue
			}

			if ch == '`' {
				state = InRawString
				pos++
				continue
			}

			if ch == '"' || ch == '\'' {
				state = InString
				strDelim = ch
				if ch == '\'' {
					state = InRune
				}
				pos++
				continue
			}

			if unicode.IsDigit(ch) {
				state = InNumber
				buffer.WriteRune(ch)
				pos++
				continue
			}

			if unicode.IsLetter(ch) || ch == '_' {
				state = InIdentifier
				buffer.WriteRune(ch)
				pos++
				continue
			}

			tokens = append(tokens, models.Token{Type: models.Error, Value: string(ch)})
			pos++

		case InIdentifier:
			if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
				buffer.WriteRune(ch)
				pos++
			} else {
				value := buffer.String()
				if keywords[value] {
					tokens = append(tokens, models.Token{Type: models.Keyword, Value: value})
				} else if value == "true" || value == "false" {
					tokens = append(tokens, models.Token{Type: models.BooleanLiteral, Value: value})
				} else {
					tokens = append(tokens, models.Token{Type: models.Identifier, Value: value})
				}
				buffer.Reset()
				state = Start
			}

		case InNumber:
			switch {
			case ch == '_':
				buffer.WriteRune(ch)
				pos++
			case ch == '.':
				buffer.WriteRune(ch)
				state = InFloat
				pos++
			case ch == 'e' || ch == 'E':
				buffer.WriteRune(ch)
				state = InExponent
				pos++
			case ch == 'x' || ch == 'X':
				buffer.WriteRune(ch)
				state = InHexNumber
				pos++
			case ch == 'o' || ch == 'O':
				buffer.WriteRune(ch)
				state = InOctalNumber
				pos++
			case ch == 'b' || ch == 'B':
				buffer.WriteRune(ch)
				state = InBinaryNumber
				pos++
			case unicode.IsDigit(ch):
				buffer.WriteRune(ch)
				pos++
			default:
				tokens = append(tokens, models.Token{Type: models.IntLiteral, Value: buffer.String()})
				buffer.Reset()
				state = Start
			}

		case InHexNumber:
			if strings.ContainsRune("0123456789abcdefABCDEF_", ch) {
				buffer.WriteRune(ch)
				pos++
			} else {
				tokens = append(tokens, models.Token{Type: models.IntLiteral, Value: buffer.String()})
				buffer.Reset()
				state = Start
			}

		case InOctalNumber:
			if strings.ContainsRune("01234567_", ch) {
				buffer.WriteRune(ch)
				pos++
			} else {
				tokens = append(tokens, models.Token{Type: models.IntLiteral, Value: buffer.String()})
				buffer.Reset()
				state = Start
			}

		case InBinaryNumber:
			if strings.ContainsRune("01_", ch) {
				buffer.WriteRune(ch)
				pos++
			} else {
				tokens = append(tokens, models.Token{Type: models.IntLiteral, Value: buffer.String()})
				buffer.Reset()
				state = Start
			}

		case InFloat:
			if ch == '_' || unicode.IsDigit(ch) {
				buffer.WriteRune(ch)
				pos++
			} else if ch == 'e' || ch == 'E' {
				buffer.WriteRune(ch)
				state = InExponent
				pos++
			} else {
				tokens = append(tokens, models.Token{Type: models.FloatLiteral, Value: buffer.String()})
				buffer.Reset()
				state = Start
			}

		case InExponent:
			if ch == '+' || ch == '-' {
				buffer.WriteRune(ch)
				state = InExponentDigits
				pos++
			} else if unicode.IsDigit(ch) {
				buffer.WriteRune(ch)
				state = InExponentDigits
				pos++
			} else {
				tokens = append(tokens, models.Token{Type: models.Error, Value: buffer.String() + string(ch)})
				buffer.Reset()
				state = Start
				pos++
			}

		case InExponentDigits:
			if unicode.IsDigit(ch) || ch == '_' {
				buffer.WriteRune(ch)
				pos++
			} else {
				tokens = append(tokens, models.Token{Type: models.FloatLiteral, Value: buffer.String()})
				buffer.Reset()
				state = Start
			}

		case InString:
			if ch == strDelim {
				tokens = append(tokens, models.Token{Type: models.StringLiteral, Value: buffer.String()})
				buffer.Reset()
				state = Start
				pos++
			} else if ch == '\\' {
				if pos+1 < len(input) {
					buffer.WriteRune(ch)
					buffer.WriteRune(rune(input[pos+1]))
					pos += 2
				} else {
					tokens = append(tokens, models.Token{Type: models.Error, Value: "unterminated escape"})
					pos++
				}
			} else {
				buffer.WriteRune(ch)
				pos++
			}

		case InRawString:
			if ch == '`' {
				tokens = append(tokens, models.Token{Type: models.StringLiteral, Value: buffer.String()})
				buffer.Reset()
				state = Start
				pos++
			} else {
				buffer.WriteRune(ch)
				pos++
			}

		case InRune:
			if ch == strDelim {
				if buffer.Len() > 0 {
					tokens = append(tokens, models.Token{Type: models.RuneLiteral, Value: buffer.String()})
				} else {
					tokens = append(tokens, models.Token{Type: models.Error, Value: "empty rune"})
				}
				buffer.Reset()
				state = Start
				pos++
			} else if ch == '\\' {
				if pos+1 < len(input) {
					buffer.WriteRune(ch)
					buffer.WriteRune(rune(input[pos+1]))
					pos += 2
				} else {
					tokens = append(tokens, models.Token{Type: models.Error, Value: "unterminated escape"})
					pos++
				}
			} else {
				buffer.WriteRune(ch)
				pos++
			}

		case InOperator:
			possibleOp := buffer.String() + string(ch)
			if operators[possibleOp] {
				buffer.WriteRune(ch)
				pos++
			} else {
				if operators[buffer.String()] {
					tokens = append(tokens, models.Token{Type: models.Operator, Value: buffer.String()})
					buffer.Reset()
					state = Start
				} else {
					tokens = append(tokens, models.Token{Type: models.Error, Value: buffer.String()})
					buffer.Reset()
					state = Start
					pos++
				}
			}

		case InLineComment:
			if ch == '\n' {
				state = Start
			}
			pos++

		case InBlockComment:
			if ch == '*' && pos+1 < len(input) && input[pos+1] == '/' {
				pos += 2
				state = Start
			} else {
				pos++
			}

		default:
			pos++
		}
	}

	return tokens
}

func isOperatorStart(ch rune) bool {
	switch ch {
	case '+', '-', '*', '/', '%', '&', '|', '^', '<', '>', '!', '=', ':', '~':
		return true
	}
	return false
}
