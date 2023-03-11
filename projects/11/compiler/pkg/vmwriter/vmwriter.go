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
	w.outputFile.WriteString("push " + getSegmentAlias(segment) + " " + strconv.Itoa(index) + "\n")
}

func (w *VMWriter) WritePop(segment MemorySegment, index int) {
	w.outputFile.WriteString("pop " + getSegmentAlias(segment) + " " + strconv.Itoa(index) + "\n")
}

func getSegmentAlias(segment MemorySegment) string {
	switch segment {
	case Constant:
		return "constant"
	case Argument:
		return "argument"
	case Local:
		return "local"
	case Static:
		return "static"
	case This:
		return "this"
	case That:
		return "that"
	case Pointer:
		return "pointer"
	case Temp:
		return "temp"
	default:
		panic("Undefined alias for segment")
	}
}

func (w *VMWriter) WriteArithmetic(command ArithmeticCommand) {
	var instruction string
	switch command {
	case Add:
		instruction = "add"
	case Sub:
		instruction = "sub"
	case Neg:
		instruction = "neg"
	case Not:
		instruction = "not"
	case Lt:
		instruction = "lt"
	case Gt:
		instruction = "gt"
	case Eq:
		instruction = "eq"
	case And:
		instruction = "and"
	default:
		panic("Undefined alias for arithemtic command")
	}
	w.outputFile.WriteString(instruction + "\n")
}

func (w *VMWriter) WriteLabel(label string) {
	w.outputFile.WriteString("label " + label + "\n")
}

func (w *VMWriter) WriteGoto(label string) {
	w.outputFile.WriteString("goto " + label + "\n")
}

func (w *VMWriter) WriteIf(label string) {
	w.outputFile.WriteString("if-goto " + label + "\n")
}

func (w *VMWriter) WriteCall(name string, nArgs int) {
	w.outputFile.WriteString("call " + name + " " + strconv.Itoa(nArgs) + "\n")
}

func (w *VMWriter) WriteFunction(name string, nArgs int) {
	w.outputFile.WriteString("function " + name + " " + strconv.Itoa(nArgs) + "\n")
}

func (w *VMWriter) WriteReturn() {
	w.outputFile.WriteString("return\n")
}
