package main

import (
	"fmt"
	"os"

	"vmtranslator/pkg/parser"
)

func main() {
	filePath := os.Args[1]
	p := parser.New(filePath)

	fmt.Printf("%v", p)
}
