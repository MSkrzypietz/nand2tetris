package parser

import (
	"bufio"
	"log"
	"os"
	"strings"
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
