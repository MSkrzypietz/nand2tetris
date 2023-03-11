package vmwriter

import (
	"os"
	"strconv"
)

type MemorySegment int

const (
	Constant MemorySegment = iota
	Argument
	Local
	Static
	This
	That
	Pointer
	Temp
)

type ArithmeticCommand int

const (
	Add ArithmeticCommand = iota
	Sub
	Neg
	Eq
	Gt
	Lt
	And
	Or
	Not
)

type VMWriter struct {
	outputFile *os.File
}

func New(outputFile *os.File) *VMWriter {
	return &VMWriter{
		outputFile: outputFile,
	}
}

func (w *VMWriter) WritePush(segment MemorySegment, index int) {
	switch segment {
	case Constant:
		w.outputFile.WriteString("push constant " + strconv.Itoa(index) + "\n")
	}
}

func (w *VMWriter) WritePop(segment MemorySegment, index int) {
	switch segment {
	case Temp:
		w.outputFile.WriteString("pop temp " + strconv.Itoa(index) + "\n")
	}
}

func (w *VMWriter) WriteArithmetic(command ArithmeticCommand) {
	var instruction string
	switch command {
	case Add:
		instruction = "add"
	}
	w.outputFile.WriteString(instruction + "\n")
}

func (w *VMWriter) WriteLabel(label string) {}

func (w *VMWriter) WriteGoto(label string) {}

func (w *VMWriter) WriteIf(label string) {}

func (w *VMWriter) WriteCall(name string, nArgs int) {
	w.outputFile.WriteString("call " + name + " " + strconv.Itoa(nArgs) + "\n")
}

func (w *VMWriter) WriteFunction(name string, nArgs int) {
	w.outputFile.WriteString("function " + name + " " + strconv.Itoa(nArgs) + "\n")
}

func (w *VMWriter) WriteReturn() {
	w.outputFile.WriteString("return\n")
}
