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
		} else if cmdType == parser.CmdPush || cmdType == parser.CmdPop {
			index, _ := strconv.Atoi(p.Arg2())
			cw.WritePushPop(cmdType, p.Arg1(), index)
		} else if cmdType == parser.CmdLabel {
			cw.WriteLabel(p.Arg1())
		} else if cmdType == parser.CmdIf {
			cw.WriteIf(p.Arg1())
		} else if cmdType == parser.CmdGoto {
			cw.WriteGoto(p.Arg1())
		}
	}

	cw.WriteEnd()
	cw.Close()
}
