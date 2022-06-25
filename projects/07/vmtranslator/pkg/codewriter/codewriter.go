package codewriter

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"vmtranslator/pkg/parser"
)

type CodeWriter struct {
	file             *os.File
	uniqueLabelIndex int
}

func New(filePath string) *CodeWriter {
	dir := filepath.Dir(filePath)
	fileParts := strings.Split(filepath.Base(filePath), ".")
	f, err := os.Create(filepath.Join(dir, fileParts[0]+".asm"))
	if err != nil {
		log.Fatal(err)
	}

	return &CodeWriter{
		file:             f,
		uniqueLabelIndex: 0,
	}
}

func (cw *CodeWriter) Close() {
	cw.file.Sync()
	cw.file.Close()
}

func (cw *CodeWriter) WriteEnd() {
	cw.writeToFile(
		"(END)",
		"@END",
		"0;JMP",
	)
}

func (cw *CodeWriter) WriteArithmetic(command string) {
	ab := newAsmBuilder()

	switch command {
	case "add":
		ab.Add(popStack()...)
		ab.Add("D=M")
		ab.Add(popStack()...)
		ab.Add("D=D+M")
		ab.Add(pushDRegToStack()...)
	case "sub":
		ab.Add(popStack()...)
		ab.Add("D=M")
		ab.Add(popStack()...)
		ab.Add("D=M-D")
		ab.Add(pushDRegToStack()...)
	case "neg":
		ab.Add(popStack()...)
		ab.Add("D=-M")
		ab.Add(pushDRegToStack()...)
	case "eq":
		ab.Add(cw.compare("JEQ")...)
	case "gt":
		ab.Add(cw.compare("JLT")...)
	case "lt":
		ab.Add(cw.compare("JGT")...)
	case "and":
		ab.Add(popStack()...)
		ab.Add("D=M")
		ab.Add(popStack()...)
		ab.Add("D=D&M")
		ab.Add(pushDRegToStack()...)
	case "or":
		ab.Add(popStack()...)
		ab.Add("D=M")
		ab.Add(popStack()...)
		ab.Add("D=D|M")
		ab.Add(pushDRegToStack()...)
	case "not":
		ab.Add(popStack()...)
		ab.Add("D=!M")
		ab.Add(pushDRegToStack()...)
	}

	cw.writeToFile(ab.Instructions()...)
}

func (cw *CodeWriter) compare(jumpInstruction string) []string {
	index := cw.getNextUniqueLabelIndex()
	onTrueLabel := "ON_TRUE_" + index
	endLabel := "END_" + index

	ab := newAsmBuilder()
	ab.Add(popStack()...)
	ab.Add("D=M")
	ab.Add(popStack()...)
	ab.Add("D=D-M")
	ab.Add("@" + onTrueLabel)
	ab.Add("D;" + jumpInstruction)
	ab.Add("D=0")
	ab.Add("@" + endLabel)
	ab.Add("0;JMP")
	ab.Add("(" + onTrueLabel + ")")
	ab.Add("D=-1")
	ab.Add("(" + endLabel + ")")
	ab.Add(pushDRegToStack()...)
	return ab.Instructions()
}

func (cw *CodeWriter) getNextUniqueLabelIndex() string {
	index := cw.uniqueLabelIndex
	cw.uniqueLabelIndex++
	return strconv.Itoa(index)
}

func (cw *CodeWriter) WritePushPop(cmdType parser.CmdType, segment string, index int) {
	ab := newAsmBuilder()

	if cmdType == parser.CmdPush {
		if segment == "constant" {
			ab.Add(
				"@"+strconv.Itoa(index),
				"D=A",
			)
			ab.Add(pushDRegToStack()...)
		}
	}

	if cmdType == parser.CmdPop {

	}

	cw.writeToFile(ab.Instructions()...)
}

func popStack() []string {
	return []string{
		"@SP",
		"M=M-1",
		"A=M",
	}
}

func pushDRegToStack() []string {
	return []string{
		"@SP",
		"A=M",
		"M=D",
		"@SP",
		"M=M+1",
	}
}

func (cw *CodeWriter) writeToFile(assembly ...string) {
	for _, asmInstruction := range assembly {
		cw.file.WriteString(asmInstruction + "\n")
	}
}
