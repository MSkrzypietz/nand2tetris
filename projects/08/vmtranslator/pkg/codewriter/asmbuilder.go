package codewriter

type asmBuilder struct {
	instructions []string
}

func newAsmBuilder() *asmBuilder {
	return &asmBuilder{
		instructions: []string{},
	}
}

func (ab *asmBuilder) Add(instructions ...string) {
	ab.instructions = append(ab.instructions, instructions...)
}

func (ab *asmBuilder) Instructions() []string {
	return ab.instructions
}
