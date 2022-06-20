package parser

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Parser struct {
	filePath      string
	lines         []string
	currlineIndex int
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
		currlineIndex: -1,
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

func (p *Parser) getCurrentLine() string {
	return strings.TrimSpace(p.lines[p.currlineIndex])
}

func (p *Parser) HasMoreLines() bool {
	return p.currlineIndex < len(p.lines)-1
}

func (p *Parser) Advance() {
	if !p.HasMoreLines() {
		return
	}

	p.currlineIndex++
	line := p.getCurrentLine()
	if strings.HasPrefix(line, "//") || line == "" {
		p.Advance()
	}
}

func (p *Parser) InstructionType() InstructionType {
	line := p.getCurrentLine()

	if strings.HasPrefix(line, "@") {
		return AInstruction
	}

	if strings.HasPrefix(line, "(") && strings.HasSuffix(line, ")") {
		return LInstruction
	}

	return CInstruction
}

func (p *Parser) Symbol() string {
	line := p.getCurrentLine()

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

	line := p.getCurrentLine()
	dest, _, _ := strings.Cut(line, "=")
	return dest
}

func (p *Parser) Comp() string {
	if p.InstructionType() != CInstruction {
		return ""
	}

	line := p.getCurrentLine()
	_, compAndJump, _ := strings.Cut(line, "=")
	comp, _, _ := strings.Cut(compAndJump, ";")
	return comp
}

func (p *Parser) Jump() string {
	if p.InstructionType() != CInstruction {
		return ""
	}

	line := p.getCurrentLine()
	_, jump, _ := strings.Cut(line, ";")
	return jump
}
