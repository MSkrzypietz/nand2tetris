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
	file *os.File
}

func New(filePath string) *CodeWriter {
	dir := filepath.Dir(filePath)
	fileParts := strings.Split(filepath.Base(filePath), ".")
	f, err := os.Create(filepath.Join(dir, fileParts[0]+".asm"))
	if err != nil {
		log.Fatal(err)
	}

	return &CodeWriter{
		file: f,
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

// "sub" "neg" "eq" "gt" "lt" "and" "or" "not"
func (cw *CodeWriter) WriteArithmetic(command string) {
	ab := newAsmBuilder()

	if command == "add" {
		ab.Add(popStack()...)
		ab.Add("D=M")
		ab.Add(popStack()...)
		ab.Add("D=D+M")
		ab.Add(pushDRegToStack()...)
	}

	cw.writeToFile(ab.Instructions()...)
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
