package main

import (
	"assembler/pkg/code"
	"assembler/pkg/parser"
	"assembler/pkg/symtable"
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
	st := symtable.New()

	currInstructionIndex := 0
	for p.HasMoreLines() {
		p.Advance()

		if p.InstructionType() == parser.LInstruction {
			st.AddEntry(p.Symbol(), currInstructionIndex)
		} else {
			currInstructionIndex++
		}
	}

	nextVariableIndex := 16
	p = parser.New(filePath)
	for p.HasMoreLines() {
		p.Advance()

		var binary string
		switch p.InstructionType() {
		case parser.AInstruction:
			symbol := p.Symbol()
			num, err := strconv.Atoi(symbol)
			if err != nil {
				if !st.Contains(symbol) {
					st.AddEntry(symbol, nextVariableIndex)
					nextVariableIndex++
				}
				num = st.GetAddress(symbol)
			} 
			binary = fmt.Sprintf("%016v", strconv.FormatInt(int64(num), 2))
		case parser.CInstruction:
			binary = "111" + code.Comp(p.Comp()) + code.Dest(p.Dest()) + code.Jump(p.Jump())
		case parser.LInstruction:
			continue
		}

		f.WriteString(binary + "\n")
	}

	f.Sync()
}
