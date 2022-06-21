package main

import (
	"assembler/pkg/code"
	"assembler/pkg/parser"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	filePath := os.Args[1]

	baseFilePath := filepath.Base(filePath)
	fileName, _, _ := strings.Cut(baseFilePath, ".")
	outputFile := fileName + ".hack"
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	p := parser.New(filePath)

	for p.HasMoreLines() {
		p.Advance()

		var binary string
		switch p.InstructionType() {
		case parser.AInstruction:
			n, _ := strconv.Atoi(p.Symbol())
			binary = fmt.Sprintf("%016v", strconv.FormatInt(int64(n), 2))
		case parser.CInstruction:
			binary = "111" + code.Comp(p.Comp()) + code.Dest(p.Dest()) + code.Jump(p.Jump())
		}

		f.WriteString(binary + "\n")
	}

	f.Sync()
}
