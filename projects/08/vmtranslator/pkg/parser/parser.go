package parser

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type CmdType int

const (
	CmdArithmetic CmdType = iota
	CmdPush
	CmdPop
	CmdLabel
	CmdGoto
	CmdIf
	CmdFunction
	CmdReturn
	CmdCall
)

type Parser struct {
	instructions         []string
	currInstructionIndex int
}

func New(filePath string) *Parser {
	return &Parser{
		instructions:         readInstructions(filePath),
		currInstructionIndex: -1,
	}
}

func readInstructions(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var instructions []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line, _, _ := strings.Cut(scanner.Text(), "//")
		if line = strings.TrimSpace(line); line != "" {
			instructions = append(instructions, line)
		}
	}
	return instructions
}

func (p *Parser) getInstructionParts() []string {
	return strings.Split(p.instructions[p.currInstructionIndex], " ")
}

func (p *Parser) HasMoreLines() bool {
	return p.currInstructionIndex < len(p.instructions)-1
}

func (p *Parser) Advance() {
	if p.HasMoreLines() {
		p.currInstructionIndex++
	}
}

func (p *Parser) CommandType() CmdType {
	switch p.getInstructionParts()[0] {
	case "add":
		return CmdArithmetic
	case "sub":
		return CmdArithmetic
	case "neg":
		return CmdArithmetic
	case "eq":
		return CmdArithmetic
	case "gt":
		return CmdArithmetic
	case "lt":
		return CmdArithmetic
	case "and":
		return CmdArithmetic
	case "or":
		return CmdArithmetic
	case "not":
		return CmdArithmetic
	case "push":
		return CmdPush
	case "pop":
		return CmdPop
	case "label":
		return CmdLabel
	case "goto":
		return CmdGoto
	case "if-goto":
		return CmdIf
	}

	return CmdLabel
}

func (p *Parser) Arg1() string {
	if p.CommandType() == CmdArithmetic {
		return p.getInstructionParts()[0]
	}

	return p.getInstructionParts()[1]
}

func (p *Parser) Arg2() string {
	return p.getInstructionParts()[2]
}
