package main

import (
	"os"

	"vmtranslator/pkg/parser"
)

func main() {
	filePath := os.Args[1]
	p := parser.New(filePath)

	for p.HasMoreLines() {
		p.Advance()

		println(p.CommandType(), p.Arg1())
	}
}
