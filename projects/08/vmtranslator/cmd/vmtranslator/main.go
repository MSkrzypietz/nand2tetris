package main

import (
	"os"
	"strconv"

	"vmtranslator/pkg/codewriter"
	"vmtranslator/pkg/parser"
)

func main() {
	filePath := os.Args[1]

	p := parser.New(filePath)
	cw := codewriter.New(filePath)

	for p.HasMoreLines() {
		p.Advance()

		cmdType := p.CommandType()

		if cmdType == parser.CmdArithmetic {
			cw.WriteArithmetic(p.Arg1())
		}

		if cmdType == parser.CmdPush || cmdType == parser.CmdPop {
			index, _ := strconv.Atoi(p.Arg2())
			cw.WritePushPop(cmdType, p.Arg1(), index)
		}
	}

	cw.WriteEnd()
	cw.Close()
}
