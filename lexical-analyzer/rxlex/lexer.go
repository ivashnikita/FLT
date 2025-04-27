package rxlex

import (
	"fmt"
	"regexp"
	"strings"

	"analyzer/models"
)

func Lex(input string) ([]models.Token, error) {
	var tokens []models.Token

	keywords := []string{
		"break", "default", "func", "interface", "select",
		"case", "defer", "go", "map", "struct",
		"chan", "else", "goto", "package", "switch",
		"const", "fallthrough", "if", "range", "type",
		"continue", "for", "import", "return", "var",
	}

	operators := [][]string{
		{"<<=", ">>=", "&^="},
		{":=", "...", "++", "--", "+=", "-=", "*=", "/=", "%=", "&=", "|=", "^=", "<<", ">>", "&&", "||", "==", "!=", "<=", ">=", "<-", "&^"},
		{"+", "-", "*", "/", "%", "&", "|", "^", "!", "=", "<", ">", "~"},
	}

	separators := map[rune]bool{
		'(': true,
		')': true,
		'[': true,
		']': true,
		'{': true,
		'}': true,
		',': true,
		';': true,
		':': true,
		'.': true,
	}

	for len(input) > 0 {
		// Delete empty lines
		if whitespace := regexp.MustCompile(`^\s+`).FindString(input); whitespace != "" {
			input = input[len(whitespace):]
			continue
		}

		// Delete comments
		if strings.HasPrefix(input, "//") {
			end := strings.Index(input, "\n")
			if end == -1 {
				input = ""
			} else {
				input = input[end:]
			}
			continue
		}

		if strings.HasPrefix(input, "/*") {
			end := strings.Index(input, "*/")
			if end == -1 {
				return nil, fmt.Errorf("unterminated block comment")
			}
			input = input[end+2:]
			continue
		}

		// Check for keywords
		keywordFound := false
		for _, kw := range keywords {
			if strings.HasPrefix(input, kw) {
				remaining := input[len(kw):]
				if len(remaining) == 0 || !isIdentifierPart(rune(remaining[0])) {
					tokens = append(tokens, models.Token{Type: models.Keyword, Value: kw})
					input = remaining
					keywordFound = true
					break
				}
			}
		}
		if keywordFound {
			continue
		}

		// Check for literals

		// `` string
		if rawStr := regexp.MustCompile("^`[^`]*`").FindString(input); rawStr != "" {
			tokens = append(tokens, models.Token{Type: models.StringLiteral, Value: rawStr})
			input = input[len(rawStr):]
			continue
		}

		// "" string
		if interpretedStr := regexp.MustCompile(`^"(?:\\.|[^"\\])*"`).FindString(input); interpretedStr != "" {
			tokens = append(tokens, models.Token{Type: models.StringLiteral, Value: interpretedStr})
			input = input[len(interpretedStr):]
			continue
		}

		// Rune
		if runeLit := regexp.MustCompile(`^'(?:\\.|.)'`).FindString(input); runeLit != "" {
			tokens = append(tokens, models.Token{Type: models.RuneLiteral, Value: runeLit})
			input = input[len(runeLit):]
			continue
		}

		// Boolean
		if strings.HasPrefix(input, "true") {
			remaining := input[4:]
			if len(remaining) == 0 || !isIdentifierPart(rune(remaining[0])) {
				tokens = append(tokens, models.Token{Type: models.BooleanLiteral, Value: "true"})
				input = remaining
				continue
			}
		}

		if strings.HasPrefix(input, "false") {
			remaining := input[5:]
			if len(remaining) == 0 || !isIdentifierPart(rune(remaining[0])) {
				tokens = append(tokens, models.Token{Type: models.BooleanLiteral, Value: "false"})
				input = remaining
				continue
			}
		}

		// Hex integer
		if hex := regexp.MustCompile(`^0[xX][0-9a-fA-F_]+`).FindString(input); hex != "" {
			tokens = append(tokens, models.Token{Type: models.IntLiteral, Value: hex})
			input = input[len(hex):]
			continue
		}

		// Binary integer
		if binary := regexp.MustCompile(`^0[bB][01_]+`).FindString(input); binary != "" {
			tokens = append(tokens, models.Token{Type: models.IntLiteral, Value: binary})
			input = input[len(binary):]
			continue
		}

		// Octal integer
		if octal := regexp.MustCompile(`^0[oO]?[0-7_]+`).FindString(input); octal != "" {
			tokens = append(tokens, models.Token{Type: models.IntLiteral, Value: octal})
			input = input[len(octal):]
			continue
		}

		// Numbers
		if num := regexp.MustCompile(`^([0-9][0-9_]*(\.[0-9_]*)?|\.[0-9_]+)([eE][+-]?[0-9_]+)?`).FindString(input); num != "" {
			if strings.ContainsAny(num, ".eE") {
				tokens = append(tokens, models.Token{Type: models.FloatLiteral, Value: num})
			} else {
				tokens = append(tokens, models.Token{Type: models.IntLiteral, Value: num})
			}
			input = input[len(num):]
			continue
		}

		// Operators
		opFound := false
		for _, group := range operators {
			for _, op := range group {
				if strings.HasPrefix(input, op) {
					tokens = append(tokens, models.Token{Type: models.Operator, Value: op})
					input = input[len(op):]
					opFound = true
					break
				}
			}

			if opFound {
				break
			}
		}

		if opFound {
			continue
		}

		// Separators
		if len(input) > 0 {
			r := rune(input[0])
			if separators[r] {
				tokens = append(tokens, models.Token{Type: models.Separator, Value: string(r)})
				input = input[1:]
				continue
			}
		}

		// Identifiers
		if id := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`).FindString(input); id != "" {
			tokens = append(tokens, models.Token{Type: models.Identifier, Value: id})
			input = input[len(id):]
			continue
		}

		// Invalid token case
		return nil, fmt.Errorf("invalid token at: %q", input[:1])
	}

	return tokens, nil
}

// It is for checking that some lexeme is standalone
func isIdentifierPart(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
}
