package vmwriter

import "os"

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

func (w *VMWriter) WritePush(segment MemorySegment, index int) {}

func (w *VMWriter) WritePop(segment MemorySegment, index int) {}

func (w *VMWriter) WriteArithmetic(command ArithmeticCommand) {}

func (w *VMWriter) WriteLabel(label string) {}

func (w *VMWriter) WriteGoto(label string) {}

func (w *VMWriter) WriteIf(label string) {}

func (w *VMWriter) WriteCall(name string, nArgs int) {}

func (w *VMWriter) WriteFunction(name string, nArgs int) {}

func (w *VMWriter) WriteReturn() {}
