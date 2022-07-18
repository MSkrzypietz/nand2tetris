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
		currTokenIndex: -1,
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
			println(i, closingIndex)
			println(string(fileData[i+1 : i+20]))
			tokens = append(tokens, fileData[i+1:i+1+closingIndex])
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
