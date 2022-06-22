package parser

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Parser struct {
	filePath             string
	lines                []string
	currLineIndex        int
	currInstructionIndex int
}

type InstructionType int

const (
	AInstruction InstructionType = iota
	CInstruction
	LInstruction
)

func New(filePath string) *Parser {
	return &Parser{
		filePath:      filePath,
		lines:         readLines(filePath),
		currLineIndex: -1,
	}
}

func readLines(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func (p *Parser) GetCurrentLine() string {
	line, _, _ := strings.Cut(p.lines[p.currLineIndex], "//")
	return strings.TrimSpace(line)
}

func (p *Parser) HasMoreLines() bool {
	return p.currLineIndex < len(p.lines)-1
}

func (p *Parser) Advance() {
	if !p.HasMoreLines() {
		return
	}

	p.currLineIndex++
	if p.GetCurrentLine() == "" {
		p.Advance()
	}
}

func (p *Parser) InstructionType() InstructionType {
	line := p.GetCurrentLine()

	if strings.HasPrefix(line, "@") {
		return AInstruction
	}

	if strings.HasPrefix(line, "(") && strings.HasSuffix(line, ")") {
		return LInstruction
	}

	return CInstruction
}

func (p *Parser) Symbol() string {
	line := p.GetCurrentLine()

	if p.InstructionType() == AInstruction {
		return line[1:]
	}

	if p.InstructionType() == LInstruction {
		return line[1:][:len(line)-2]
	}

	return ""
}

func (p *Parser) Dest() string {
	if p.InstructionType() != CInstruction {
		return ""
	}

	line := p.GetCurrentLine()
	dest, _, found := strings.Cut(line, "=")
	if !found {
		return ""
	}
	return dest
}

func (p *Parser) Comp() string {
	if p.InstructionType() != CInstruction {
		return ""
	}

	line := p.GetCurrentLine()
	_, compAndJump, found := strings.Cut(line, "=")
	if found {
		comp, _, _ := strings.Cut(compAndJump, ";")
		return comp
	}
	comp, _, _ := strings.Cut(line, ";")
	return comp
}

func (p *Parser) Jump() string {
	if p.InstructionType() != CInstruction {
		return ""
	}

	line := p.GetCurrentLine()
	_, jump, _ := strings.Cut(line, ";")
	return jump
}
