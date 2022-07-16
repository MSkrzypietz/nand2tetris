package main

import (
	"os"
	"syntaxlyzer/pkg/tokenizer"
)

func main() {
	filename := os.Args[1]
	tokenizer := tokenizer.New(filename)
	for tokenizer.HasMoreTokens() {
		tokenizer.Advance()
		println(tokenizer.GetToken())
	}
}
