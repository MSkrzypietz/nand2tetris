package tokenizer

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type TokenType int

const (
	Keyword TokenType = iota
	Symbol
	Identifier
	IntConst
	StringConst
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
	fileData := readFileDatatWithoutSingleLineComments(filename)
	fileData = removeMultiLineComments(fileData)

	var tokens []string
	var token string
	for i := 0; i < len(fileData); i++ {
		char := string(fileData[i])
		if char == " " {
			tokens = addNonEmptyToken(tokens, token)
			token = ""
		} else if isSymbol(char) {
			tokens = addNonEmptyToken(tokens, token)
			tokens = append(tokens, char)
			token = ""
		} else if char == "\"" {
			tokens = addNonEmptyToken(tokens, token)
			closingIndex := strings.Index(string(fileData[i+1:]), "\"")
			tokens = append(tokens, fileData[i:i+2+closingIndex])
			i += closingIndex + 1
			token = ""
		} else {
			token += char
		}
	}
	tokens = addNonEmptyToken(tokens, token)
	return tokens
}

func readFileDatatWithoutSingleLineComments(filename string) string {
	file, _ := os.Open(filename)
	defer file.Close()

	fileData := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line, _, _ := strings.Cut(scanner.Text(), "//")
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			fileData += trimmedLine
		}
	}
	return fileData
}

func removeMultiLineComments(fileData string) string {
	beforeCommentText, _, foundComment := strings.Cut(fileData, "/*")
	if foundComment {
		_, afterCommentText, _ := strings.Cut(fileData, "*/")
		return beforeCommentText + removeMultiLineComments(afterCommentText)
	}
	return beforeCommentText
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

func (t *Tokenizer) TokenType() TokenType {
	currToken := t.tokens[t.currTokenIndex]
	if isSymbol(currToken) {
		return Symbol
	} else if isKeyword(currToken) {
		return Keyword
	} else if isIntConst(currToken) {
		return IntConst
	} else if isStringConst(currToken) {
		return StringConst
	}
	return Identifier
}

func isIntConst(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func isStringConst(s string) bool {
	if len(s) < 2 {
		return false
	}
	return string(s[0]) == "\"" && string(s[len(s)-1]) == "\""
}

func isKeyword(s string) bool {
	for _, symbol := range getKeywords() {
		if s == symbol {
			return true
		}
	}
	return false
}

func (t *Tokenizer) getToken() string {
	return t.tokens[t.currTokenIndex]
}

func (t *Tokenizer) KeyWord() string {
	return t.getToken()
}

func (t *Tokenizer) Symbol() string {
	return t.getToken()
}

func (t *Tokenizer) Identifier() string {
	return t.getToken()
}

func (t *Tokenizer) IntVal() int {
	val, _ := strconv.Atoi(t.getToken())
	return val
}

func (t *Tokenizer) StringVal() string {
	return t.getToken()[1 : len(t.getToken())-1]
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
