package codewriter

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"vmtranslator/pkg/parser"
)

type CodeWriter struct {
	file              *os.File
	fileName          string
	uniqueLabelIndex  int
	functionCallIndex int
}

func New(filePath string) *CodeWriter {
	dir := filepath.Dir(filePath)
	dirName := filepath.Base(dir)
	f, err := os.Create(filepath.Join(dir, dirName+".asm"))
	if err != nil {
		log.Fatal(err)
	}

	cw := &CodeWriter{
		file:              f,
		fileName:          "",
		uniqueLabelIndex:  0,
		functionCallIndex: 0,
	}
	cw.writeBootstrap()
	return cw
}

func (cw *CodeWriter) SetFileName(fileName string) {
	cw.fileName = fileName
}

func (cw *CodeWriter) Close() {
	cw.file.Sync()
	cw.file.Close()
}

func (cw *CodeWriter) writeBootstrap() {
	ab := newAsmBuilder()
	ab.Add("@256")
	ab.Add("D=A")
	ab.Add("@SP")
	ab.Add("M=D")
	cw.writeToFile(ab.Instructions()...)

	cw.WriteCall("Sys.init", 0)
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
		ab.Add(cw.compare("JGT")...)
	case "lt":
		ab.Add(cw.compare("JLT")...)
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
	ab.Add("D=M-D")
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
	if cmdType == parser.CmdPush {
		cw.writePush(segment, index)
	} else if cmdType == parser.CmdPop {
		cw.writePop(segment, index)
	}
}

func (cw *CodeWriter) writePop(segment string, index int) {
	ab := newAsmBuilder()
	segmentAddress := cw.getSegmentAddress(segment, index)

	if segment == "temp" || segment == "pointer" {
		ab.Add(popStack()...)
		ab.Add("D=M")
		ab.Add("@" + segmentAddress)
		ab.Add("M=D")
	} else {
		if segment == "static" {
			ab.Add("@" + segmentAddress)
			ab.Add("D=A")
		} else {
			ab.Add("@" + segmentAddress)
			ab.Add("D=M")
		}
		ab.Add("@" + strconv.Itoa(index))
		ab.Add("D=D+A")
		ab.Add("@R13")
		ab.Add("M=D")
		ab.Add(popStack()...)
		ab.Add("D=M")
		ab.Add("@R13")
		ab.Add("A=M")
		ab.Add("M=D")
	}

	cw.writeToFile(ab.Instructions()...)
}

func (cw *CodeWriter) writePush(segment string, index int) {
	ab := newAsmBuilder()
	segmentAddress := cw.getSegmentAddress(segment, index)

	ab.Add("@" + segmentAddress)
	if segment == "constant" {
		ab.Add("D=A")
	} else if segment == "temp" || segment == "pointer" || segment == "static" {
		ab.Add("D=M")
	} else {
		ab.Add("D=M")
		ab.Add("@" + strconv.Itoa(index))
		ab.Add("A=D+A")
		ab.Add("D=M")
	}
	ab.Add(pushDRegToStack()...)

	cw.writeToFile(ab.Instructions()...)
}

func (cw *CodeWriter) getSegmentAddress(segment string, index int) string {
	switch segment {
	case "constant":
		return strconv.Itoa(index)
	case "local":
		return "LCL"
	case "argument":
		return "ARG"
	case "this":
		return "THIS"
	case "that":
		return "THAT"
	case "pointer":
		if index == 0 {
			return "THIS"
		} else {
			return "THAT"
		}
	case "static":
		return cw.fileName + "." + strconv.Itoa(index)
	case "temp":
		return "R" + strconv.Itoa(5+index)
	default:
		return ""
	}
}

func (cw *CodeWriter) WriteLabel(label string) {
	cw.writeToFile("(" + label + ")")
}

func (cw *CodeWriter) WriteIf(label string) {
	ab := newAsmBuilder()

	ab.Add(popStack()...)
	ab.Add("D=M")
	ab.Add("@" + label)
	ab.Add("D;JNE")

	cw.writeToFile(ab.Instructions()...)
}

func (cw *CodeWriter) WriteGoto(label string) {
	ab := newAsmBuilder()

	ab.Add("@" + label)
	ab.Add("0;JMP")

	cw.writeToFile(ab.Instructions()...)
}

func (cw *CodeWriter) WriteFunction(functionName string, nVars int) {
	ab := newAsmBuilder()

	ab.Add("(" + functionName + ")")
	ab.Add("D=0")
	for i := 0; i < nVars; i++ {
		ab.Add(pushDRegToStack()...)
	}

	cw.writeToFile(ab.Instructions()...)
}

func (cw *CodeWriter) WriteCall(functionName string, nVars int) {
	ab := newAsmBuilder()
	retLabel := functionName + "$ret." + strconv.Itoa(cw.functionCallIndex)
	cw.functionCallIndex++

	ab.Add("@" + retLabel)
	ab.Add("D=A")
	ab.Add(pushDRegToStack()...)
	ab.Add(pushAddressToStack("LCL")...)
	ab.Add(pushAddressToStack("ARG")...)
	ab.Add(pushAddressToStack("THIS")...)
	ab.Add(pushAddressToStack("THAT")...)

	ab.Add("@5")
	ab.Add("D=A")
	ab.Add("@" + strconv.Itoa(nVars))
	ab.Add("D=D+A")
	ab.Add("@SP")
	ab.Add("D=M-D")
	ab.Add("@ARG")
	ab.Add("M=D")

	ab.Add("@SP")
	ab.Add("D=M")
	ab.Add("@LCL")
	ab.Add("M=D")

	ab.Add("@" + functionName)
	ab.Add("0;JMP")

	ab.Add("(" + retLabel + ")")

	cw.writeToFile(ab.Instructions()...)
}

func (cw *CodeWriter) WriteReturn() {
	ab := newAsmBuilder()

	ab.Add("@LCL")
	ab.Add("D=M")
	ab.Add("@R13")
	ab.Add("M=D")

	ab.Add(setSegmentAddressToFrameOffset("R14", "R13", 5)...)

	ab.Add(popStack()...)
	ab.Add("D=M")
	ab.Add("@ARG")
	ab.Add("A=M")
	ab.Add("M=D")

	ab.Add("@ARG")
	ab.Add("D=M+1")
	ab.Add("@SP")
	ab.Add("M=D")

	ab.Add(setSegmentAddressToFrameOffset("THAT", "R13", 1)...)
	ab.Add(setSegmentAddressToFrameOffset("THIS", "R13", 2)...)
	ab.Add(setSegmentAddressToFrameOffset("ARG", "R13", 3)...)
	ab.Add(setSegmentAddressToFrameOffset("LCL", "R13", 4)...)

	ab.Add("@R14")
	ab.Add("A=M")
	ab.Add("0;JMP")

	cw.writeToFile(ab.Instructions()...)
}

func setSegmentAddressToFrameOffset(segment, frame string, offset int) []string {
	return []string{
		"@" + frame,
		"D=M",
		"@" + strconv.Itoa(offset),
		"A=D-A",
		"D=M",
		"@" + segment,
		"M=D",
	}
}

func popStack() []string {
	return []string{
		"@SP",
		"M=M-1",
		"A=M",
	}
}

func pushAddressToStack(address string) []string {
	return append([]string{"@" + address, "D=M"}, pushDRegToStack()...)
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
