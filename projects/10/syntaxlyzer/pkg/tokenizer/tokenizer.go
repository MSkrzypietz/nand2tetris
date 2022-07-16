package tokenizer

import (
	"bufio"
	"os"
	"strings"
)

type Tokenizer struct {
	tokens         []string
	currTokenIndex int
}

func New(filename string) *Tokenizer {
	return &Tokenizer{
		tokens:         readTokens(filename),
		currTokenIndex: 0,
	}
}

func readTokens(filename string) []string {
	file, _ := os.Open(filename)
	defer file.Close()

	var tokens []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line, _, _ := strings.Cut(scanner.Text(), "//")
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}

		token := ""
		for i := 0; i < len(trimmedLine)-1; i++ {
			char := string(trimmedLine[i])
			if char == " " {
				tokens = addNonEmptyToken(tokens, token)
				token = ""
				continue
			} else if isSymbol(char) {
				tokens = addNonEmptyToken(tokens, token)
				tokens = append(tokens, char)
				token = ""
				continue
			}
			token += char
		}
		tokens = addNonEmptyToken(tokens, token)
	}
	return tokens
}

func addNonEmptyToken(tokens []string, token string) []string {
	trimmedToken := strings.TrimSpace(token)
	if trimmedToken != "" {
		return append(tokens, trimmedToken)
	}
	return tokens
}

func (t *Tokenizer) HasMoreTokens() bool {
	return t.currTokenIndex < len(t.tokens)-1
}

func (t *Tokenizer) Advance() {
	t.currTokenIndex++
}

func (t *Tokenizer) GetToken() string {
	return t.tokens[t.currTokenIndex]
}

func getKeywords() []string {
	return []string{
		"class",
		"constructor",
		"function",
		"method",
		"field",
		"static",
		"var",
		"int",
		"char",
		"boolean",
		"void",
		"true",
		"false",
		"null",
		"this",
		"let",
		"do",
		"if",
		"else",
		"while",
		"return",
	}
}

func isSymbol(s string) bool {
	for _, symbol := range getSymbols() {
		if s == symbol {
			return true
		}
	}
	return false
}

func getSymbols() []string {
	return []string{
		"{",
		"}",
		"(",
		")",
		"[",
		"]",
		",",
		".",
		";",
		"+",
		"-",
		"*",
		"/",
		"&",
		"|",
		"<",
		">",
		"=",
		"~",
	}
}
