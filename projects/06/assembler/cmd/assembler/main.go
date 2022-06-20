package main

import (
	"assembler/pkg/parser"
	"fmt"
	"os"
)

func main() {
	filePath := os.Args[1]

	parser := parser.New(filePath)

	for parser.HasMoreLines() {
		parser.Advance()

		fmt.Printf("InstructionType %v\n", parser.InstructionType())
		fmt.Printf("Symbol %v\n", parser.Symbol())
		fmt.Printf("dest %v\n", parser.Dest())
		fmt.Printf("comp %v\n", parser.Comp())
		fmt.Printf("jump %v\n", parser.Jump())
	}
}
